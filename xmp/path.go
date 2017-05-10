// Copyright (c) 2017 Alexander Eichhorn
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package xmp

import (
	"encoding"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Path string

type PathValue struct {
	Path      Path      `json:"path"`
	Namespace string    `json:"namespace,omitempty"`
	Value     string    `json:"value"`
	Flags     SyncFlags `json:"flags,omitempty"`
}

type PathValueList []PathValue

type byPath PathValueList

func (l byPath) Len() int           { return len(l) }
func (l byPath) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byPath) Less(i, j int) bool { return l[i].Path.String() < l[j].Path.String() }

type byValue PathValueList

func (l byValue) Len() int           { return len(l) }
func (l byValue) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byValue) Less(i, j int) bool { return l[i].Value < l[j].Value }

func (x Path) String() string {
	return string(x)
}

// ns:tagname/subtagname
func (x Path) IsXmpPath() bool {
	return strings.Index(string(x), ":") > -1
}

func (x Path) NamespacePrefix() string {
	if i := strings.Index(string(x), ":"); i > -1 {
		return string(x[:i])
	}
	return string(x)
}

func (x Path) Namespace() (*Namespace, error) {
	if i := strings.Index(string(x), ":"); i > -1 {
		if ns, err := GetNamespace(string(x[:i])); err != nil {
			return nil, err
		} else {
			return ns, nil
		}
	}
	return nil, fmt.Errorf("xmp: invalid path '%s'", x.String())
}

func (x Path) MatchNamespace(ns *Namespace) bool {
	if ns == nil {
		return false
	}
	if i := strings.Index(string(x), ":"); i > -1 {
		return ns.GetName() == string(x[:i])
	}
	return false
}

func (x Path) Len() int {
	return strings.Count(string(x), "/") + 1
}

func (x Path) Fields() []string {
	s := string(x)
	if i := strings.Index(string(x), ":"); i > -1 {
		return strings.Split(s[i+1:], "/")
	}
	return nil
}

func (x *Path) UnmarshalText(data []byte) error {
	p := Path(string(data))
	if !p.IsXmpPath() {
		return fmt.Errorf("xmp: invalid path '%s'", p.String())
	}
	*x = p
	return nil
}

func (d *Document) GetPath(path Path) (string, error) {
	if !path.IsXmpPath() {
		return "", fmt.Errorf("xmp: invalid path '%s'", path.String())
	}
	ns, err := path.Namespace()
	if err != nil {
		return "", err
	}
	m := d.FindModel(ns)
	if m == nil {
		return "", fmt.Errorf("xmp: model not found for path '%s'", path.String())
	}
	return GetPath(m, path)
}

func (d *Document) SetPath(desc PathValue) error {
	flags := desc.Flags
	path := desc.Path
	value := desc.Value

	if flags == 0 {
		flags = S_DEFAULT
	}
	if !path.IsXmpPath() {
		return fmt.Errorf("xmp: invalid path '%s'", path.String())
	}

	ns, err := path.Namespace()
	if err != nil {
		return err
	}

	m := d.FindModel(ns)
	if m == nil {
		if flags&S_CREATE == 0 {
			return fmt.Errorf("xmp: create flag required to make model for '%s'", path.String())
		}
		m = ns.NewModel()
		d.AddModel(m)
	}

	// empty source will only be used with delete flag
	if value == "" && flags&S_DELETE == 0 {
		return fmt.Errorf("xmp: delete flag required for empty '%s'", path.String())
	}

	// get current version of path value
	dest, err := GetPath(m, path)
	if err != nil {
		return err
	}

	// skip equal
	if dest == value {
		return fmt.Errorf("xmp: no update required for '%s'", path.String())
	}

	// empty destination values require create flag
	if dest == "" && flags&S_CREATE == 0 {
		return fmt.Errorf("xmp: create flag required to make new attribute at '%s'", path.String())
	}

	// existing destination values require replace/delete/append/unique flag
	if dest != "" && flags&(S_REPLACE|S_DELETE|S_APPEND|S_UNIQUE) == 0 {
		return fmt.Errorf("xmp: update flag required to change existing attribute at '%s'", path.String())
	}

	err = SetPath(m, path, value, flags)
	if err == nil {
		d.SetDirty()
	}
	return err
}

func (d *Document) ListPaths() (PathValueList, error) {
	l := make(PathValueList, 0)
	for _, v := range d.Nodes {
		if v.Model != nil {
			if pvl, err := ListPaths(v.Model); err != nil {
				return nil, err
			} else {
				l = append(l, pvl...)
			}
		}
	}
	sort.Sort(byPath(l))
	return l, nil
}

func SetPath(v Model, path Path, value string, flags SyncFlags) error {
	if flags == 0 {
		flags = S_DEFAULT
	}
	if !path.IsXmpPath() {
		return fmt.Errorf("xmp: invalid path '%s'", path.String())
	}

	val := derefIndirect(v)
	ns, _ := path.Namespace()
	// fmt.Printf("SET path=%s (%d)\n", path.String(), path.Length())

	// follow path
	l := path.Len()
	for i, name := range path.Fields() {

		var lang string
		var idx int
		// split lang or array index from name, be safe with slice indexes
		if k := strings.Index(name, "["); k > 0 && len(name) > k+2 {
			s := strings.TrimSuffix(name[k+1:], "]")
			if j, err := strconv.Atoi(s); err == nil {
				idx = j
			} else {
				lang = s
			}
			name = name[:k]
		}
		// fmt.Printf("> name=%s lang=%s idx=%d\n", name, lang, idx)

		finfo, err := findField(val, ns.Expand(name), "xmp")
		if err != nil {
			return fmt.Errorf("path field %d (%s) not found: %v", i, name, err)
		}

		// ignore `any` fields
		if finfo.flags&fAny > 0 {
			return nil
		}

		// allocate memory for pointer values in structs
		fv := finfo.value(val)
		if !fv.IsValid() {
			return nil
		}
		if fv.Type().Kind() == reflect.Ptr && fv.IsNil() && fv.CanSet() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}
		fv = derefValue(fv)
		typ := fv.Type()
		// fmt.Printf("> found %s (%s)\n", fv.Type().Name(), fv.Kind().String())

		// continue loop when field is a struct and we're not at the end
		if fv.Kind() == reflect.Struct && i < l-1 {
			// fmt.Printf("> is a struct, recursing...\n")
			val = fv
			continue
		}

		// handle Array type fields, AltString type is always at the end of a path
		if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo != nil && finfo.flags&fArray > 0 || pv.Type().Implements(arrayType)) {
				if _, ok := pv.Interface().(ExtensionArray); ok {
					// skip extension arrays for now
					return nil
				} else if altArr, ok := pv.Interface().(*AltString); ok {
					// fmt.Printf("Setting AltString ptr value\n")
					switch {
					case flags&S_UNIQUE > 0 && value != "":
						// append source when not exist
						altArr.AddUnique(lang, value)
						return nil
					case flags&S_APPEND > 0 && value != "":
						// append source value
						altArr.Add(lang, value)
						return nil
					case flags&S_REPLACE > 0 && value != "":
						// replace entire AltString with a new version
						if lang != "" {
							altArr.Set(lang, value)
						} else {
							*altArr = NewAltString(value)
						}
						return nil
					case flags&S_DELETE > 0 && value == "":
						// delete the entire AltString or specific language
						if lang != "" {
							altArr.RemoveLang(lang)
						} else {
							pv.Set(reflect.Zero(typ))
						}
						return nil
					}
					return nil
				} else {
					// are we at the end of the path?
					if i < l-1 {
						if fv.Len() <= idx {
							return fmt.Errorf("cannot set index %d of slice %s", idx, name)
						}
						// fmt.Printf("Recursing into array ptr %v on index %d\n", typ.Name(), idx)
						val = derefValue(fv.Index(idx))
						continue
					}

					// Note: we're not setting values at slice/array index; if the path
					// ends in an array operator it is ignored
					switch {
					case flags&S_UNIQUE > 0 && value != "":
						// append source when not exist
						if strArr, ok := pv.Interface().(*StringArray); ok {
							strArr.AddUnique(value)
						} else if strList, ok := pv.Interface().(*StringList); ok {
							strList.AddUnique(value)
						}
						//  else if idList, ok := pv.Interface().(*IdentifierArray); ok {
						// 	idList.AddUnique(value)
						// }
						return nil
					case flags&S_APPEND > 0 && value != "":
						// append source value
						if strArr, ok := pv.Interface().(*StringArray); ok {
							strArr.Add(value)
						} else if strList, ok := pv.Interface().(*StringList); ok {
							strList.Add(value)
						}
						//  else if idList, ok := pv.Interface().(*IdentifierArray); ok {
						// 	idList.Add(value)
						// }
						return nil
					case flags&S_REPLACE > 0 && value != "":
						// replace entire list with a new version
						if strArr, ok := pv.Interface().(*StringArray); ok {
							*strArr = NewStringArray(value)
						} else if strList, ok := pv.Interface().(*StringList); ok {
							*strList = NewStringList(value)
						}
						//  else if idList, ok := pv.Interface().(*IdentifierArray); ok {
						// 	*idList = NewIdentifierArray(value)
						// }
						return nil
					case flags&S_DELETE > 0 && value == "":
						// delete the entire array
						pv.Set(reflect.Zero(typ))
						return nil
					}
					return nil
				}
			}
		}

		// Text
		if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo != nil && finfo.flags&fTextUnmarshal > 0 || pv.Type().Implements(textUnmarshalerType)) {
				// log.Debugf(">> unmarshal model %s calling UnmarshalText addr with '%s'\n", src.FullName(), src.Value)
				return pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
			}
		}

		// otherwise set simple value directly
		if err := setValue(fv, value); err != nil {
			return err
		}
	}
	return nil
}

func GetPath(v Model, path Path) (string, error) {

	val := derefIndirect(v)
	ns, _ := path.Namespace()
	// fmt.Printf("GET path=%s (%d)\n", path.String(), path.Length())

	// follow path
	l := path.Len()
	for i, name := range path.Fields() {
		var lang string
		var idx int
		// split lang or array index from name, be safe with slice indexes
		if k := strings.Index(name, "["); k > 0 && len(name) > k+2 {
			s := strings.TrimSuffix(name[k+1:], "]")
			if j, err := strconv.Atoi(s); err == nil {
				idx = j
			} else {
				lang = s
			}
			name = name[:k]
		}
		// fmt.Printf("> name=%s lang=%s idx=%d\n", name, lang, idx)

		finfo, err := findField(val, ns.Expand(name), "xmp")
		if err != nil {
			return "", fmt.Errorf("path field %d (%s) not found: %v", i, name, err)
		}

		// ignore `any` fields
		if finfo.flags&fAny > 0 {
			return "", nil
		}

		fv := finfo.value(val)
		typ := fv.Type()

		// ignore empty fields
		if !fv.IsValid() {
			return "", nil
		}
		if (fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr) && fv.IsNil() {
			return "", nil
		}
		if finfo.flags&fEmpty == 0 && isEmptyValue(fv) {
			return "", nil
		}

		// Drill into interfaces and pointers.
		// This can turn into an infinite loop given a cyclic chain,
		// but it matches the Go 1 behavior.
		for fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}
		// fmt.Printf("> found %s (%s)\n", fv.Type().Name(), fv.Kind().String())

		// continue loop when field is a struct and we're not at the end
		if fv.Kind() == reflect.Struct && i < l-1 {
			val = fv
			// fmt.Printf("> is a struct, recursing...\n")
			continue
		}

		// handle Array type fields (AltString types are always at the end of a path)
		if fv.CanInterface() && (finfo != nil && finfo.flags&fArray > 0 || typ.Implements(arrayType)) {
			if _, ok := fv.Interface().(ExtensionArray); ok {
				// skip extension arrays for now
				return "", nil
			} else if altArr, ok := fv.Interface().(AltString); ok {
				// fmt.Printf("> is an Alt Array\n")
				if lang != "" {
					return altArr.Get(lang), nil
				}
				return altArr.Default(), nil
			} else {
				// fmt.Printf("> is normal Array\n")
				if fv.Len() <= idx {
					// fmt.Printf("> is a short Array\n")
					return "", nil
				}
				if i < l {
					// fmt.Printf("> recursing to index %d...\n", idx)
					val = derefValue(fv.Index(idx))
					continue
				} else {
					fv = derefValue(fv.Index(idx))
					typ = fv.Type()
					// fmt.Printf("> using index %d %s\n", idx, typ.Name())
				}
			}
		} else if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo != nil && finfo.flags&fArray > 0 || pv.Type().Implements(arrayType)) {
				if _, ok := pv.Interface().(*ExtensionArray); ok {
					// skip extension arrays for now
					return "", nil
				} else if altArr, ok := pv.Interface().(*AltString); ok {
					// fmt.Printf("> is an Alt Array ptr\n")
					if lang != "" {
						return altArr.Get(lang), nil
					}
					return altArr.Default(), nil
				} else {
					// fmt.Printf("> is normal Array ptr\n")
					if fv.Len() <= idx {
						// fmt.Printf("> is a short Array\n")
						return "", nil
					}
					if i < l-1 {
						// fmt.Printf("> recursing to index %d...\n", idx)
						val = derefValue(fv.Index(idx))
						continue
					} else {
						fv = derefValue(fv.Index(idx))
						typ = fv.Type()
						// fmt.Printf("> using index %d %s\n", idx, typ.Name())
					}
				}
			}
		}

		// Check for text marshaler and marshal as string
		if fv.CanInterface() && (finfo != nil && finfo.flags&fTextMarshal > 0 || typ.Implements(textMarshalerType)) {
			b, err := fv.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil || b == nil {
				return "", err
			}
			return string(b), nil
		}

		if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo != nil && finfo.flags&fTextMarshal > 0 || pv.Type().Implements(textMarshalerType)) {
				b, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil || b == nil {
					return "", err
				}
				return string(b), nil
			}
		}

		// simple values are just fine, but any other type (slice, array, struct)
		// without textmarshaler will fail
		if s, b, err := marshalSimple(typ, fv); err != nil {
			return "", err
		} else {
			if b != nil {
				s = string(b)
			}
			return s, nil
		}
	}
	return "", nil
}

func ListPaths(v Model) (PathValueList, error) {
	return listPaths(reflect.ValueOf(v), v.Namespaces()[0].GetName(), "")
}

func listPaths(val reflect.Value, ns, prefix string) (PathValueList, error) {

	if !val.IsValid() {
		return nil, nil
	}

	if isEmptyValue(val) {
		// log.Debugf("xmp: skipping empty value for %s\n", name)
		return nil, nil
	}

	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}

	typ := val.Type()
	tinfo, err := getTypeInfo(typ, "xmp")
	if err != nil {
		return nil, err
	}

	pvl := make(PathValueList, 0)
	if prefix != "" {
		prefix = prefix + "/"
	}

	// go through all fields
	for _, finfo := range tinfo.fields {

		fv := finfo.value(val)

		if !fv.IsValid() {
			continue
		}

		if (fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr) && fv.IsNil() {
			continue
		}

		if finfo.flags&fEmpty == 0 && isEmptyValue(fv) {
			continue
		}

		fname := dropNsName(finfo.name)
		// fmt.Printf("> %s:%s%s", ns, prefix, fname)

		// Drill into interfaces and pointers.
		// This can turn into an infinite loop given a cyclic chain,
		// but it matches the Go 1 behavior.
		for fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}
		typ = fv.Type()

		// handle Array type fields (AltString types are always at the end of a path)
		if fv.CanInterface() && (finfo.flags&fArray > 0 || typ.Implements(arrayType)) {
			if _, ok := fv.Interface().(ExtensionArray); ok {
				// skip extension arrays for now
				continue
			} else if altArr, ok := fv.Interface().(AltString); ok {
				// fmt.Printf("--> is AltString, unpacking...\n")
				for _, v := range altArr {
					lang := v.Lang
					if lang == "" && v.IsDefault {
						lang = "x-default"
					}
					// fmt.Printf("%s\n", fmt.Sprintf("%s:%s%s[%s]", ns, prefix, fname, lang))
					pvl = append(pvl, PathValue{
						Path:  Path(fmt.Sprintf("%s:%s%s[%s]", ns, prefix, fname, lang)),
						Value: v.Value,
					})
				}
			} else {
				// fmt.Printf("--> is Array, unpacking...\n")
				for i, l := 0, fv.Len(); i < l; i++ {
					name := fmt.Sprintf("%s%s[%d]", prefix, fname, i)
					v := derefValue(fv.Index(i))
					switch v.Kind() {
					case reflect.Struct, reflect.Slice, reflect.Array:
						// fmt.Printf("> recursing to %s\n", name)
						l, err := listPaths(v, ns, name)
						if err != nil {
							return nil, err
						}
						pvl = append(pvl, l...)
					default:
						// fmt.Printf("> marshaling simple2 ptr\n")
						if s, b, err := marshalSimple(v.Type(), v); err != nil {
							return nil, err
						} else {
							if b != nil {
								s = string(b)
							}
							pvl = append(pvl, PathValue{
								Path:  Path(ns + ":" + name),
								Value: s,
							})
						}
					}
				}
			}
			continue
		} else if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo.flags&fArray > 0 || pv.Type().Implements(arrayType)) {
				if _, ok := pv.Interface().(*ExtensionArray); ok {
					// skip extension arrays for now
					continue
				} else if altArr, ok := pv.Interface().(*AltString); ok {
					// fmt.Printf("--> is AltString ptr, unpacking...\n")
					for _, v := range *altArr {
						lang := v.Lang
						if lang == "" && v.IsDefault {
							lang = "x-default"
						}
						pvl = append(pvl, PathValue{
							Path:  Path(fmt.Sprintf("%s:%s%s[%s]", ns, prefix, fname, lang)),
							Value: v.Value,
						})
					}
				} else {
					// fmt.Printf("--> is Array ptr, unpacking...\n")
					for i, l := 0, fv.Len(); i < l; i++ {
						name := fmt.Sprintf("%s%s[%d]", prefix, fname, i)
						v := derefValue(fv.Index(i))
						switch v.Kind() {
						case reflect.Struct, reflect.Slice, reflect.Array:
							// fmt.Printf("> recursing to %s\n", name)
							l, err := listPaths(v, ns, name)
							if err != nil {
								return nil, err
							}
							pvl = append(pvl, l...)
						default:
							// fmt.Printf("> marshaling simple2\n")
							if s, b, err := marshalSimple(v.Type(), v); err != nil {
								return nil, err
							} else {
								if b != nil {
									s = string(b)
								}
								pvl = append(pvl, PathValue{
									Path:  Path(ns + ":" + name),
									Value: s,
								})
							}
						}
					}
				}
				continue
			}
		}

		// Check for text marshaler and marshal as string
		if fv.CanInterface() && (finfo.flags&fTextMarshal > 0 || typ.Implements(textMarshalerType)) {
			// fmt.Printf("> marshaling text\n")
			b, err := fv.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil || b == nil {
				return nil, err
			}
			pvl = append(pvl, PathValue{
				Path:  Path(ns + ":" + prefix + fname),
				Value: string(b),
			})
			continue
		}

		if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo.flags&fTextMarshal > 0 || pv.Type().Implements(textMarshalerType)) {
				// fmt.Printf("> marshaling text ptr\n")
				b, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil || b == nil {
					return nil, err
				}
				pvl = append(pvl, PathValue{
					Path:  Path(ns + ":" + prefix + fname),
					Value: string(b),
				})
				continue
			}
		}

		// continue iteration when field is a struct without text marshaler
		if fv.Kind() == reflect.Struct {
			// fmt.Printf("--> is a struct, recursing...\n")
			l, err := listPaths(fv, ns, prefix+fname)
			if err != nil {
				return nil, err
			}
			pvl = append(pvl, l...)
			continue
		}

		// otherwise marshal as value
		// fmt.Printf("> marshaling simple\n")
		if s, b, err := marshalSimple(typ, fv); err != nil {
			return nil, err
		} else {
			if b != nil {
				s = string(b)
			}
			pvl = append(pvl, PathValue{
				Path:  Path(ns + ":" + prefix + fname),
				Value: s,
			})
		}
	}

	sort.Sort(byPath(pvl))
	return pvl, nil
}
