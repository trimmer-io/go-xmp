// Copyright (c) 2017-2018 Alexander Eichhorn
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

func (x *PathValueList) Add(p Path, v string) {
	if v == "" {
		return
	}
	*x = append(*x, PathValue{Path: p, Value: v})
}

func (x *PathValueList) AddFlags(p Path, v string, f SyncFlags) {
	if v == "" {
		return
	}
	*x = append(*x, PathValue{Path: p, Value: v, Flags: f})
}

func (x PathValueList) Find(p Path) *PathValue {
	for _, v := range x {
		if v.Path == p {
			return &v
		}
	}
	return nil
}

// assumes a sorted list
func (x PathValueList) Unique() PathValueList {
	l := make(PathValueList, 0, len(x))
	var last Path
	for _, v := range x {
		if last != v.Path {
			l = append(l, v)
			last = v.Path
		}
	}
	return l
}

func (x PathValueList) Diff(y PathValueList) PathValueList {
	if len(x) == 0 {
		return y
	}
	if len(y) == 0 {
		return x
	}
	diff := make(PathValueList, 0, Max(len(x), len(y)))
	for _, xv := range x {
		yv := y.Find(xv.Path)
		if yv == nil {
			diff.AddFlags(xv.Path, xv.Value, DELETE)
			continue
		}
		if yv.Value != xv.Value {
			diff.AddFlags(xv.Path, xv.Value, REPLACE)
		}
	}
	for _, yv := range y {
		xv := x.Find(yv.Path)
		if xv == nil {
			diff.AddFlags(yv.Path, yv.Value, CREATE)
		}
	}
	return diff
}

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

func NewPath(prefix string, segments ...string) Path {
	return Path(prefix + ":" + strings.Join(segments, "/"))
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

func (x Path) Namespace(d *Document) (*Namespace, error) {
	if i := strings.Index(string(x), ":"); i > -1 {
		if ns := d.findNsByPrefix(string(x[:i])); ns != nil {
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
	s := string(x)
	if i := strings.Index(s, ":"); i > -1 {
		s = strings.TrimPrefix(s[i+1:], "/")
		if len(s) > 0 {
			return strings.Count(s, "/") + 1
		}
	}
	return 0
}

func (x Path) Fields() []string {
	s := string(x)
	if i := strings.Index(s, ":"); i > -1 {
		y := strings.TrimPrefix(s[i+1:], "/")
		if len(y) > 0 {
			return strings.Split(y, "/")
		}
	}
	return nil
}

func (x Path) AppendIndex(i int) Path {
	f := "%s[%d]"
	if strings.HasSuffix(x.String(), "]") {
		f = "%s/[%d]"
	}
	return Path(fmt.Sprintf(f, x.String(), i))
}

func (x Path) AppendIndexString(s string) Path {
	f := "%s[%s]"
	if strings.HasSuffix(x.String(), "]") {
		f = "%s/[%s]"
	}
	return Path(fmt.Sprintf(f, x.String(), s))
}

func (x Path) Push(segment ...string) Path {
	f := x.Fields()
	f = append(f, segment...)
	return NewPath(x.NamespacePrefix(), f...)
}

func (x Path) Pop() (string, Path) {
	f := x.Fields()
	switch l := len(f); l {
	case 0:
		return "", x
	case 1:
		return f[0], NewPath(x.NamespacePrefix())
	default:
		return f[l-1], NewPath(x.NamespacePrefix(), f[:l-1]...)
	}
}

func (x Path) PopFront() (string, Path) {
	switch f := x.Fields(); len(f) {
	case 0:
		return "", x
	case 1:
		seg := f[0]
		ns := x.NamespacePrefix()
		if hasPrefix(seg) {
			ns = getPrefix(seg)
			seg = stripPrefix(seg)
		}
		return seg, NewPath(ns)
	default:
		seg := f[0]
		ns := x.NamespacePrefix()
		if hasPrefix(seg) {
			ns = getPrefix(seg)
			seg = stripPrefix(seg)
		}
		return seg, NewPath(ns, f[1:]...)
	}
}

func (x Path) PeekNamespacePrefix() string {
	switch f := x.Fields(); len(f) {
	case 0:
		return x.NamespacePrefix()
	default:
		if hasPrefix(f[0]) {
			return getPrefix(f[0])
		}
		return x.NamespacePrefix()
	}
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
	ns, err := path.Namespace(d)
	if err != nil {
		return "", err
	}
	n := d.FindNode(ns)
	if n == nil {
		return "", fmt.Errorf("xmp: path '%s' not found", path.String())
	}
	if n.Model != nil {
		if s, err := GetModelPath(n.Model, path); err == nil {
			return s, nil
		} else if err != errNotFound {
			return "", fmt.Errorf("xmp: path '%s' error: %v", path.String(), err)
		}
	}
	if v, err := n.GetPath(path); err != nil {
		if err == errNotFound {
			return "", fmt.Errorf("xmp: path '%s' not found", path.String())
		} else {
			return "", fmt.Errorf("xmp: path '%s' error: %v", path.String(), err)
		}
	} else {
		return v, nil
	}
}

func GetModelPath(v Model, path Path) (string, error) {
	val := derefIndirect(v)
	l := path.Len()
	for n, walker := path.PopFront(); n != ""; n, walker = walker.PopFront() {
		// fmt.Printf("%d name=%s walker=%s\n", l-walker.Len(), n, walker.String())
		name, idx, lang := parsePathSegment(n)
		if idx < -1 {
			return "", fmt.Errorf("path field %d (%s): invalid index", l-walker.Len(), n)
		}

		finfo, err := findField(val, name, "xmp")
		if err != nil {
			return "", errNotFound
			// return "", fmt.Errorf("path field %d (%s) not found: %v", i, name, err)
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
		if finfo.flags&fOmit > 0 || (finfo.flags&fEmpty == 0 && isEmptyValue(fv)) {
			return "", nil
		}

		// Drill into interfaces and pointers.
		for fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}

		// continue loop when field is a struct and we're not at the end yet
		if fv.Kind() == reflect.Struct && walker.Len() > 0 {
			val = fv
			continue
		}

		// handle XMP array types
		av := fv
		isArray := false
		if fv.CanInterface() && (finfo != nil && finfo.flags&fArray > 0 || typ.Implements(arrayType)) {
			isArray = true
		} else if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo != nil && finfo.flags&fArray > 0 || pv.Type().Implements(arrayType)) {
				av = pv
				isArray = true
			}
		}

		if isArray {
			switch arr := av.Interface().(type) {
			case ExtensionArray:
				if walker.Len() == 0 || fv.Len() <= idx {
					return "", errNotFound
				}
				ext := arr[idx]
				name = walker.PeekNamespacePrefix()
				node := ext.FindNodeByName(name)
				if node == nil {
					return "", errNotFound
				}
				if node.Model != nil {
					val = reflect.Indirect(reflect.ValueOf(node.Model))
					continue
				}
				return node.GetPath(walker)

			case NamedExtensionArray:
				node := arr.FindNodeByName(name)
				if node == nil {
					return "", errNotFound
				}
				return node.GetPath(walker)

			case AltString:
				if lang != "" {
					return arr.Get(lang), nil
				}
				return arr.Default(), nil
			default:
				// sanitize array index
				if idx < 0 {
					idx = 0
				}
				if fv.Len() < idx {
					return "", errNotFound
				}
				if walker.Len() > 0 {
					if av.Len() <= idx {
						return "", errNotFound
					}
					val = derefValue(av.Index(idx))
					continue
				} else {
					if av.Len() <= idx {
						return "", errNotFound
					}
					fv = derefValue(av.Index(idx))
					typ = fv.Type()
					finfo = nil
				}
			}
		}

		// Check for text marshaler and marshal as string
		av = fv
		isText := false
		if fv.CanInterface() && (finfo != nil && finfo.flags&fTextMarshal > 0 || typ.Implements(textMarshalerType)) {
			isText = true
		} else if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo != nil && finfo.flags&fTextMarshal > 0 || pv.Type().Implements(textMarshalerType)) {
				av = pv
				isText = true
			}
		}

		if isText {
			b, err := av.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil || b == nil {
				return "", err
			}
			return string(b), nil
		}

		// handle maps
		if fv.Kind() == reflect.Map {
			if finfo.flags&fFlat == 0 {
				n, walker = walker.PopFront()
				name, _, _ = parsePathSegment(n)
			}
			val := fv.MapIndex(reflect.ValueOf(name))
			if !val.IsValid() {
				return "", errNotFound
			}
			// process as simple value below
			fv = val
			typ = val.Type()
		}

		// simple values are just fine, but any other type (slice, array, struct)
		// without textmarshaler will fail here
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

func (d *Document) SetPath(desc PathValue) error {
	flags := desc.Flags
	path := desc.Path
	value := desc.Value

	if flags == 0 {
		flags = DEFAULT
	}
	if !path.IsXmpPath() {
		return fmt.Errorf("xmp: invalid path '%s'", path.String())
	}

	ns, err := path.Namespace(d)
	if ns == nil || err != nil {
		if desc.Namespace != "" {
			ns = &Namespace{path.NamespacePrefix(), desc.Namespace, nil}
			Register(ns)
		} else {
			return err
		}
	}

	m := d.FindModel(ns)
	n := d.FindNode(ns)
	if m == nil && n == nil && flags&CREATE == 0 {
		if flags&NOFAIL > 0 {
			return nil
		}
		return fmt.Errorf("xmp: create flag required to make model for '%s'", path.String())
	}
	if m == nil {
		m = ns.NewModel()
		if m != nil {
			n, _ = d.AddModel(m)
		}
	}
	if m == nil && n == nil {
		n = d.nodes.AddNode(NewNode(ns.RootName()))
	}

	// empty source will only be used with delete flag
	if value == "" && flags&DELETE == 0 {
		if flags&NOFAIL > 0 {
			return nil
		}
		return fmt.Errorf("xmp: delete flag required for empty '%s'", path.String())
	}

	// delete model when only namespace is set
	if path.Len() == 0 && flags&DELETE > 0 {
		d.RemoveNamespace(ns)
	}

	// get current version of path value
	var dest string
	if m != nil {
		dest, err = GetModelPath(m, path)
	}
	if m == nil || err == errNotFound {
		dest, err = n.GetPath(path)
	}
	if err != nil {
		return err
	}

	// skip when equal
	if dest == value {
		return nil
	}

	// empty destination values require create flag
	if dest == "" && flags&(CREATE|APPEND|UNIQUE) == 0 {
		if flags&NOFAIL > 0 {
			return nil
		}
		return fmt.Errorf("xmp: create flag required to make new attribute at '%s'", path.String())
	}

	// existing destination values require replace/delete/append/unique flag
	if dest != "" && flags&(REPLACE|DELETE|APPEND|UNIQUE) == 0 {
		if flags&NOFAIL > 0 {
			return nil
		}
		return fmt.Errorf("xmp: update flag required to change existing attribute at '%s'", path.String())
	}

	if m != nil {
		if err = SetModelPath(m, path, value, flags); err != nil && err == errNotFound {
			err = n.SetPath(path, value, flags)
		}
	} else {
		err = n.SetPath(path, value, flags)
	}
	if err == nil {
		d.SetDirty()
	}

	if flags&NOFAIL > 0 {
		return nil
	}
	return err
}

func SetModelPath(v Model, path Path, value string, flags SyncFlags) error {
	if flags == 0 {
		flags = DEFAULT
	}
	if !path.IsXmpPath() {
		return fmt.Errorf("xmp: invalid path '%s'", path.String())
	}

	val := derefIndirect(v)

	l := path.Len()
	for n, walker := path.PopFront(); n != ""; n, walker = walker.PopFront() {
		name, idx, lang := parsePathSegment(n)
		if idx < -1 {
			return fmt.Errorf("path field %d (%s): invalid index", l-walker.Len(), n)
		}

		// using the short-form of name here to find attribute names across namespaces
		// (e.g. some models use different namespace structs internally)
		finfo, err := findField(val, name, "xmp")
		if err != nil {
			// special error is catched by caller to retry setting as raw node
			return errNotFound
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

		// continue loop when field is a struct and we're not at the end
		if fv.Kind() == reflect.Struct && walker.Len() > 0 {
			val = fv
			continue
		}

		// handle maps
		if fv.Kind() == reflect.Map {
			// use proper name depending on flattening
			if finfo.flags&fFlat == 0 {
				n, walker = walker.PopFront()
				name, _, _ = parsePathSegment(n)
			}

			t := fv.Type()
			if fv.IsNil() {
				fv.Set(reflect.MakeMap(t))
			}
			switch t.Key().Kind() {
			case reflect.String:
			default:
				return fmt.Errorf("map key type must be string")
			}

			switch {
			case flags&DELETE > 0 && value == "":
				// remove map value
				fv.SetMapIndex(reflect.ValueOf(name), reflect.Zero(t.Elem()))
				return nil
			case flags&(REPLACE|CREATE) > 0 && value != "":
				// set map value
				switch t.Elem().Kind() {
				case reflect.String:
					fv.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value).Convert(t.Elem()))
				case reflect.Struct:
					val = reflect.New(t.Elem()).Elem()
					// FIXME: this recursion may not work as expected because the map
					// value is not updated when set later on
					fv.SetMapIndex(reflect.ValueOf(name), val)
					continue
				default:
					// FIXME: this does not allow for struct pointers as map values
					mval := reflect.New(t.Elem()).Elem()
					if mval.Type().Kind() == reflect.Ptr && mval.IsNil() && mval.CanSet() {
						mval.Set(reflect.New(mval.Type().Elem()))
					}
					if mval.CanInterface() && mval.Type().Implements(textUnmarshalerType) {
						if err := mval.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
							return err
						}
					} else if mval.CanAddr() {
						pv := mval.Addr()
						if pv.CanInterface() && pv.Type().Implements(textUnmarshalerType) {
							if err := pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
								return err
							}
						}
					} else {
						if err := setValue(mval, value); err != nil {
							return err
						}
					}
					fv.SetMapIndex(reflect.ValueOf(name), mval)
					return nil
				}
			}
			continue
		}

		// handle Slice/Array type fields as pointers to make them settable
		av := fv
		if fv.CanAddr() {
			av = fv.Addr()
		}
		if av.CanInterface() && (finfo != nil && finfo.flags&fArray > 0 || av.Type().Implements(arrayType)) {
			switch arr := av.Interface().(type) {
			case *ExtensionArray:
				// unwrap extension namespace and keep path remainder
				name = walker.PeekNamespacePrefix()

				// sanitize idx (it is -1 when no array index was specified)
				if idx < 0 {
					idx = 0
				}

				// add empty nodes up to idx if necessary
				if l := len(*arr); l <= idx {
					if flags&(CREATE|APPEND) == 0 && value != "" {
						return fmt.Errorf("CREATE flag required to add extension %s on path %s", name, path)
					}
					// grow slice and fill with initialized nodes
					for ; l <= idx; l++ {
						*arr = append(*arr, (*Extension)(NewNode(EmptyName)))
					}
				}

				// use the node at index
				ext := (*arr)[idx]
				node := (*Node)(ext)

				// there is two types of extensions
				if node.Name() != "" || node.FullName() == "rdf:Description" {
					// Type 1: model and nodes on top-level (name not empty)
					node = ext.FindNodeByName(name)
					if node == nil {
						if ns, err := GetNamespace(name); err == nil {
							node.Model = ns.NewModel()
						}
					}
					if node.Model != nil {
						// FIXME: when the model does not contain the path we might
						// want to store it as child node, however, the for{} loop
						// allows no back-tracking on error
						val = reflect.Indirect(reflect.ValueOf(node.Model))
						continue
					}
					// store as child node
					return node.SetPath(walker, value, flags)
				} else {
					// Type 2: models on child level (empty name); used for xmpMM:Pantry
					child := node.Nodes.FindNodeByName(name)
					if child == nil {
						child = node.AddNode(NewNode(NewName(name)))
					}
					if child.Model == nil {
						if ns, err := GetNamespace(name); err == nil {
							child.Model = ns.NewModel()
						}
					}
					if child.Model != nil {
						// FIXME: when the model does not contain the path we might
						// want to store it as extension model child node. however,
						// the for{} loop allows no back-tracking on error
						val = reflect.Indirect(reflect.ValueOf(child.Model))
						continue
					}
					// store as child node
					return child.SetPath(NewPath(name, walker.Fields()...), value, flags)
				}

			case *NamedExtensionArray:
				// Named extensions do not contain models
				// unwrap extension namespace and keep path remainder
				name, walker = walker.PopFront()
				node := arr.FindNodeByName(name)
				if node == nil {
					// create new extension node without model
					if flags&(CREATE|APPEND) == 0 && value != "" {
						return fmt.Errorf("CREATE flag required to add extension %s on path %s", name, path)
					}
					node = NewNode(NewName(name))
					ext := (*Extension)(node)
					fv.Set(reflect.Append(fv, reflect.ValueOf(ext)))
				}
				// store as child node
				return node.SetPath(walker, value, flags)

			case *AltString:
				switch {
				case flags&UNIQUE > 0 && value != "":
					// append source when not exist
					if !arr.AddUnique(lang, value) {
						return fmt.Errorf("equal value exists for %s on path %s", name, path)
					}
				case flags&APPEND > 0 && value != "":
					// append source value
					arr.Add(lang, value)
				case flags&(REPLACE|CREATE) > 0 && value != "":
					// replace entire AltString with a new version
					if lang != "" {
						arr.Set(lang, value)
					} else {
						*arr = NewAltString(value)
					}
				case flags&DELETE > 0 && value == "":
					// delete the entire AltString or just a specific language
					if lang != "" {
						arr.RemoveLang(lang)
					} else {
						fv.Set(reflect.Zero(fv.Type()))
					}
				}
				return nil

			default:
				// deref the pointer
				av = reflect.Indirect(av)
				// are we at the end of the path?
				if walker.Len() > 0 {
					// sanitize idx (it is -1 when no array index was specified)
					if idx < 0 {
						idx = 0
					}
					// handle arrays along the path
					if av.Len() <= idx {
						if flags&DELETE > 0 && value == "" {
							// would be deleting smth inside a non-existent element
							return nil
						}
						if flags&(CREATE|APPEND) == 0 {
							return fmt.Errorf("CREATE flag required to grow slice %s to index %d", name, idx)
						}
						if av.Kind() == reflect.Array {
							return fmt.Errorf("array %s index %d out of bounds", name, idx)
						}

						// add empty items up to idx
						// FIXME: when an error occurs later on, we cannot remove slice entry
						growSlice(fv, idx)
						val = reflect.New(fv.Type().Elem())
						fv.Index(idx).Set(val.Elem())
						val = derefValue(fv.Index(idx))

						// unmarshal text
						if val.CanInterface() && val.Type().Implements(textUnmarshalerType) {
							return val.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
						}

						// recurse into struct
						if reflect.Indirect(val).Kind() == reflect.Struct {
							continue
						}

						// when reaching this, path is too long
						return errNotFound
					}

					// recurse to next path segment using the value at slice index
					val = derefValue(fv.Index(idx))
					continue
				}

				// arrays are fixed size
				if av.Kind() == reflect.Array {
					if flags&REPLACE == 0 && value != "" {
						return fmt.Errorf("REPLACE flag required for setting array value in %s", name)
					}
					if flags&DELETE == 0 && value == "" {
						return fmt.Errorf("DELETE flag required to clear array value in %s", name)
					}

					// sanitize idx (it is -1 when no array index was specified)
					if idx < 0 {
						idx = 0
					}
					if av.Len() <= idx {
						return fmt.Errorf("array %s index %d out of bounds", name, idx)
					}

					v := reflect.New(fv.Type().Elem())
					if v.CanInterface() && v.Type().Implements(textUnmarshalerType) {
						if err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
							return err
						}
					} else if err := setValue(v, value); err != nil {
						return err
					}
					av.Index(idx).Set(v.Elem())
					return nil
				}

				switch {
				// case flags&UNIQUE > 0 && value != "" && idx > -1:
				// ignore the unique flag when setting with absolute slice index
				case flags&UNIQUE > 0 && value != "" && idx == -1:
					// append if unique
					for i, l := 0, fv.Len(); i < l; i++ {
						v := derefValue(fv.Index(i))

						isText := false
						vv := v
						if v.CanInterface() && v.Type().Implements(textMarshalerType) {
							isText = true
						} else if v.CanAddr() {
							pv := v.Addr()
							if pv.CanInterface() && pv.Type().Implements(textMarshalerType) {
								vv = pv
								isText = true
							}
						}

						if isText {
							b, err := vv.Interface().(encoding.TextMarshaler).MarshalText()
							if err != nil || b == nil {
								return err
							}
							if value == string(b) {
								return fmt.Errorf("equal value exists for %s on path %s", name, path)
							}
						} else {
							if s, b, err := marshalSimple(v.Type(), v); err != nil {
								return err
							} else {
								if b != nil {
									s = string(b)
								}
								if value == s {
									return fmt.Errorf("equal value exists for %s on path %s", name, path)
								}
							}
						}
					}

					v := reflect.New(fv.Type().Elem())
					if v.CanInterface() && v.Type().Implements(textUnmarshalerType) {
						if err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
							return err
						}
					} else if err := setValue(v, value); err != nil {
						return err
					}

					fv.Set(reflect.Append(fv, v.Elem()))

				case flags&APPEND > 0 && value != "" && idx == -1:
					// always append
					v := reflect.New(fv.Type().Elem())
					if v.CanInterface() && v.Type().Implements(textUnmarshalerType) {
						if err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
							return err
						}
					} else if err := setValue(v, value); err != nil {
						return err
					}

					fv.Set(reflect.Append(fv, v.Elem()))

				case flags&(REPLACE|CREATE) > 0 && value != "" && idx > -1:
					// replace or create slice element at index
					v := reflect.New(fv.Type().Elem())
					if v.CanInterface() && v.Type().Implements(textUnmarshalerType) {
						if err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
							return err
						}
					} else if err := setValue(v, value); err != nil {
						return err
					}
					// optionally growing the slice
					growSlice(fv, idx)
					// overwrite current entry,
					fv.Index(idx).Set(v.Elem())

				case flags&(REPLACE|CREATE) > 0 && value != "" && idx == -1:
					// replace entire slice with a single element
					v := reflect.New(fv.Type().Elem())
					if v.CanInterface() && v.Type().Implements(textUnmarshalerType) {
						if err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value)); err != nil {
							return err
						}
					} else if err := setValue(v, value); err != nil {
						return err
					}
					fv.Set(reflect.Zero(fv.Type()))
					fv.Set(reflect.Append(fv, v.Elem()))

				case flags&DELETE > 0 && value == "" && idx > -1:
					// delete slice element at index
					l := fv.Len()
					if l <= idx {
						return fmt.Errorf("slice %s index %d out of bounds", name, idx)
					}
					switch {
					case idx == 0:
						fv.Set(fv.Slice(1, l))
					case idx == l-1:
						fv.Set(fv.Slice(0, idx))
					default:
						fv.Set(reflect.AppendSlice(fv.Slice3(0, idx, l-1), fv.Slice(idx+1, l)))
					}

				case flags&DELETE > 0 && value == "" && idx == -1:
					// delete the entire slice
					fv.Set(reflect.Zero(fv.Type()))
				default:
					return fmt.Errorf("unsupported flag combination %v", flags)
				}
				return nil
			}
		}

		// for simple values, check current value before replacing
		if flags&REPLACE == 0 && value != "" && !isEmptyValue(fv) {
			return fmt.Errorf("REPLACE flag required for overwriting value at %s on path %s", name, path)
		}

		// Text
		if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo != nil && finfo.flags&fTextUnmarshal > 0 || pv.Type().Implements(textUnmarshalerType)) {
				return pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
			}
		}

		// otherwise set simple value directly
		return setValue(fv, value)
	}
	return nil
}

func (d *Document) ListPaths() (PathValueList, error) {
	// sync individual models to establish correct XMP entries
	if err := d.syncToXMP(); err != nil {
		return nil, err
	}
	l := make(PathValueList, 0)
	for _, v := range d.nodes {
		if v.Model != nil {
			if pvl, err := ListModelPaths(v.Model); err != nil {
				return nil, err
			} else {
				l = append(l, pvl...)
			}
		}
		r, err := v.ListPaths(NewPath(v.Name()))
		if err != nil {
			return nil, err
		}
		l = append(l, r...)
	}
	sort.Sort(byPath(l))
	return l.Unique(), nil
}

func ListModelPaths(v Model) (PathValueList, error) {
	return listPaths(reflect.ValueOf(v), NewPath(v.Namespaces()[0].GetName()))
}

func listPaths(val reflect.Value, path Path) (PathValueList, error) {
	if !val.IsValid() {
		return nil, nil
	}

	if isEmptyValue(val) {
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

	// walk all fields
	for _, finfo := range tinfo.fields {
		fv := finfo.value(val)

		if !fv.IsValid() {
			continue
		}

		if (fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr) && fv.IsNil() {
			continue
		}

		if finfo.flags&fOmit > 0 || (finfo.flags&fEmpty == 0 && isEmptyValue(fv)) {
			continue
		}

		// allow changing the namespace on first-level paths to support
		// multi-namespace models like exif, iptc and quicktime
		if path.Len() == 0 {
			ns := path.NamespacePrefix()
			realNs := getPrefix(finfo.name)
			if ns != realNs {
				path = NewPath(realNs)
			}
		}

		fname := stripPrefix(finfo.name)
		if hasPrefix(finfo.name) && getPrefix(finfo.name) != path.NamespacePrefix() {
			fname = finfo.name
		}

		// Drill into interfaces and pointers.
		for fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}
		typ = fv.Type()

		// handle XMP array types
		av := fv
		isArray := false
		if fv.CanInterface() && (finfo.flags&fArray > 0 || typ.Implements(arrayType)) {
			isArray = true
		} else if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo.flags&fArray > 0 || pv.Type().Implements(arrayType)) {
				av = pv
				isArray = true
			}
		}

		if isArray {
			switch arr := av.Interface().(type) {
			case ExtensionArray:
				for i, v := range arr {
					subpath := path.Push(fname).AppendIndex(i)
					// name := fmt.Sprintf("%s[%d]", fname, i)
					if v.Model != nil {
						if l, err := listPaths(reflect.ValueOf(v.Model), subpath); err != nil {
							return nil, err
						} else {
							pvl = append(pvl, l...)
						}
					}
					if v.XMLName.Local == "" || (*Node)(v).FullName() == "rdf:Description" {
						for _, child := range v.Nodes {
							if child.Model != nil {
								if l, err := listPaths(reflect.ValueOf(child.Model), subpath); err != nil {
									return nil, err
								} else {
									pvl = append(pvl, l...)
								}
							} else {
								if l, err := child.ListPaths(subpath); err != nil {
									return nil, err
								} else {
									pvl = append(pvl, l...)
								}
							}
						}
					} else {
						node := (*Node)(v)
						if l, err := node.ListPaths(subpath); err != nil {
							return nil, err
						} else {
							pvl = append(pvl, l...)
						}
					}
				}
			case NamedExtensionArray:
				for _, v := range arr {
					node := (*Node)(v)
					if l, err := node.ListPaths(path.Push(fname, node.Name())); err != nil {
						return nil, err
					} else {
						pvl = append(pvl, l...)
					}
				}

			case AltString:
				// AltString types are always at the end of a path
				for _, v := range arr {
					pvl.Add(path.Push(fname).AppendIndexString(v.GetLang()), v.Value)
				}
			default:
				for i, l := 0, av.Len(); i < l; i++ {
					v := derefValue(av.Index(i))

					// check for text marshaler
					isText := false
					vv := v
					if v.CanInterface() && v.Type().Implements(textMarshalerType) {
						isText = true
					} else if v.CanAddr() {
						pv := v.Addr()
						if pv.CanInterface() && pv.Type().Implements(textMarshalerType) {
							vv = pv
							isText = true
						}
					}

					if isText {
						b, err := vv.Interface().(encoding.TextMarshaler).MarshalText()
						if err != nil || b == nil {
							return nil, err
						}
						pvl.Add(path.Push(fname).AppendIndex(i), string(b))
						continue
					}

					switch v.Kind() {
					case reflect.Struct, reflect.Slice, reflect.Array:
						l, err := listPaths(v, path.Push(fname).AppendIndex(i))
						if err != nil {
							return nil, err
						}
						pvl = append(pvl, l...)
					default:
						if s, b, err := marshalSimple(v.Type(), v); err != nil {
							return nil, err
						} else {
							if b != nil {
								s = string(b)
							}
							pvl.Add(path.Push(fname).AppendIndex(i), s)
						}
					}
				}
			}
			continue
		}

		// Check for text marshaler and marshal as string
		av = fv
		isText := false
		if fv.CanInterface() && (finfo.flags&fTextMarshal > 0 || typ.Implements(textMarshalerType)) {
			isText = true
		} else if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && (finfo.flags&fTextMarshal > 0 || pv.Type().Implements(textMarshalerType)) {
				av = pv
				isText = true
			}
		}

		if isText {
			b, err := av.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil || b == nil {
				return nil, err
			}
			pvl.Add(path.Push(fname), string(b))
			continue
		}

		// continue iteration when field is a struct without text marshaler
		if fv.Kind() == reflect.Struct {
			l, err := listPaths(fv, path.Push(fname))
			if err != nil {
				return nil, err
			}
			pvl = append(pvl, l...)
			continue
		}

		// handle maps
		if fv.Kind() == reflect.Map {
			for _, key := range fv.MapKeys() {
				// need a string representation of key and value here
				val := fv.MapIndex(key)
				ks, kb, kerr := marshalSimple(key.Type(), key)
				if kerr != nil {
					return nil, kerr
				}
				vs, vb, verr := marshalSimple(val.Type(), val)
				if verr != nil {
					return nil, verr
				}
				if kb != nil {
					ks = string(kb)
				}
				if vb != nil {
					vs = string(vb)
				}
				if finfo.flags&fFlat == 0 {
					pvl.Add(path.Push(fname, ks), vs)
				} else {
					pvl.Add(path.Push(ks), vs)
				}
			}
			continue
		}

		// otherwise marshal as value
		if s, b, err := marshalSimple(typ, fv); err != nil {
			return nil, err
		} else {
			if b != nil {
				s = string(b)
			}
			pvl.Add(path.Push(fname), s)
		}
	}

	sort.Sort(byPath(pvl))
	return pvl, nil
}

func parsePathSegment(name string) (string, int, string) {
	var lang string
	var idx int = -1
	// split lang or array index from name, be safe with slice indexes
	if k := strings.Index(name, "["); k > 0 && len(name) > k+1 {
		s := strings.TrimSuffix(name[k+1:], "]")
		if len(s) == 0 {
			idx = 0
		} else if j, err := strconv.Atoi(s); err == nil {
			idx = j
		} else {
			lang = s
		}
		name = name[:k]
	}
	return name, idx, lang
}

func growSlice(v reflect.Value, n int) {
	if v.Kind() != reflect.Slice {
		return
	}
	if l := v.Len(); l <= n {
		if v.Cap() <= n {
			ncap := Max(n*2, 4)
			nv := reflect.MakeSlice(v.Type(), l, ncap)
			reflect.Copy(nv, v)
			v.Set(nv)
		}
		v.SetLen(n + 1)
	}
}
