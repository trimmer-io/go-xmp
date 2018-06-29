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
	"strings"
	"sync"
)

// typeInfo holds details for the xml representation of a type.
type typeInfo struct {
	fields []fieldInfo
}

// fieldInfo holds details for the xmp representation of a single field.
type fieldInfo struct {
	idx        []int
	name       string
	minVersion Version
	maxVersion Version
	flags      fieldFlags
}

func (f fieldInfo) String() string {
	s := []string{fmt.Sprintf("field %s (%v)", f.name, f.idx)}
	if !f.minVersion.IsZero() {
		s = append(s, "vmin", f.minVersion.String())
	}
	if !f.maxVersion.IsZero() {
		s = append(s, "vmax", f.maxVersion.String())
	}
	if f.flags&fAttr > 0 {
		s = append(s, "Attr")
	}
	if f.flags&fEmpty > 0 {
		s = append(s, "Empty")
	}
	if f.flags&fOmit > 0 {
		s = append(s, "Omit")
	}
	if f.flags&fAny > 0 {
		s = append(s, "Any")
	}
	if f.flags&fFlat > 0 {
		s = append(s, "Flat")
	}
	if f.flags&fArray > 0 {
		s = append(s, "Array")
	}
	if f.flags&fBinaryMarshal > 0 {
		s = append(s, "BinaryMarshal")
	}
	if f.flags&fBinaryUnmarshal > 0 {
		s = append(s, "BinaryUnmarshal")
	}
	if f.flags&fTextMarshal > 0 {
		s = append(s, "TextMarshal")
	}
	if f.flags&fTextUnmarshal > 0 {
		s = append(s, "TextUnmarshal")
	}
	if f.flags&fMarshal > 0 {
		s = append(s, "Marshal")
	}
	if f.flags&fUnmarshal > 0 {
		s = append(s, "Unmarshal")
	}
	return strings.Join(s, " ")
}

type fieldFlags int

const (
	fElement fieldFlags = 1 << iota
	fAttr
	fEmpty
	fOmit
	fAny
	fFlat
	fArray
	fBinaryMarshal
	fBinaryUnmarshal
	fTextMarshal
	fTextUnmarshal
	fMarshal
	fUnmarshal
	fMarshalAttr
	fUnmarshalAttr
	fMode = fElement | fAttr | fEmpty | fOmit | fAny | fFlat | fArray | fBinaryMarshal | fBinaryUnmarshal | fTextMarshal | fTextUnmarshal | fMarshal | fUnmarshal | fMarshalAttr | fUnmarshalAttr
)

type tinfoMap map[reflect.Type]*typeInfo

var tinfoNsMap = make(map[string]tinfoMap)
var tinfoLock sync.RWMutex

var (
	binaryUnmarshalerType = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
	binaryMarshalerType   = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	textUnmarshalerType   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	textMarshalerType     = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	marshalerType         = reflect.TypeOf((*Marshaler)(nil)).Elem()
	unmarshalerType       = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	attrMarshalerType     = reflect.TypeOf((*MarshalerAttr)(nil)).Elem()
	attrUnmarshalerType   = reflect.TypeOf((*UnmarshalerAttr)(nil)).Elem()
	arrayType             = reflect.TypeOf((*Array)(nil)).Elem()
	zeroType              = reflect.TypeOf((*Zero)(nil)).Elem()
	stringerType          = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

// getTypeInfo returns the typeInfo structure with details necessary
// for marshaling and unmarshaling typ.
func getTypeInfo(typ reflect.Type, ns string) (*typeInfo, error) {
	if ns == "" {
		ns = "xmp"
	}
	tinfoLock.RLock()
	m, ok := tinfoNsMap[ns]
	if !ok {
		m = make(tinfoMap)
		tinfoLock.RUnlock()
		tinfoLock.Lock()
		tinfoNsMap[ns] = m
		tinfoLock.Unlock()
		tinfoLock.RLock()
	}
	tinfo, ok := m[typ]
	tinfoLock.RUnlock()
	if ok {
		return tinfo, nil
	}
	tinfo = &typeInfo{}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xmp: type %s is not a struct", typ.String())
	}
	n := typ.NumField()
	for i := 0; i < n; i++ {
		f := typ.Field(i)
		if (f.PkgPath != "" && !f.Anonymous) || f.Tag.Get(ns) == "-" {
			continue // Private field
		}

		// For embedded structs, embed its fields.
		if f.Anonymous {
			t := f.Type
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() == reflect.Struct {
				inner, err := getTypeInfo(t, ns)
				if err != nil {
					return nil, err
				}
				for _, finfo := range inner.fields {
					finfo.idx = append([]int{i}, finfo.idx...)
					if err := addFieldInfo(typ, tinfo, &finfo, ns); err != nil {
						return nil, err
					}
				}
				continue
			}
		}

		finfo, err := structFieldInfo(typ, &f, ns)
		if err != nil {
			return nil, err
		}

		// Add the field if it doesn't conflict with other fields.
		if err := addFieldInfo(typ, tinfo, finfo, ns); err != nil {
			return nil, err
		}
	}
	tinfoLock.Lock()
	m[typ] = tinfo
	tinfoLock.Unlock()
	return tinfo, nil
}

// structFieldInfo builds and returns a fieldInfo for f.
func structFieldInfo(typ reflect.Type, f *reflect.StructField, ns string) (*fieldInfo, error) {
	finfo := &fieldInfo{idx: f.Index}
	// Split the tag from the xml namespace if necessary.
	tag := f.Tag.Get(ns)

	// Parse flags.
	tokens := strings.Split(tag, ",")
	if len(tokens) == 1 {
		finfo.flags = fElement
	} else {
		tag = tokens[0]
		for _, flag := range tokens[1:] {
			switch flag {
			case "attr":
				finfo.flags |= fAttr
			case "empty":
				finfo.flags |= fEmpty
			case "omit":
				finfo.flags |= fOmit
			case "any":
				finfo.flags |= fAny
			case "flat":
				finfo.flags |= fFlat
			}

			// dissect version(s)
			//   v1.0     - only write in version v1.0
			//   v1.0+    - starting at and after v1.0
			//   v1.0-    - only write before and including v1.0
			//   v1.0<1.2 - write from v1.0 until v1.2
			if strings.HasPrefix(flag, "v") {
				flag = flag[1:]
				var op rune
				tokens := strings.FieldsFunc(flag, func(r rune) bool {
					switch r {
					case '+', '-', '<':
						op = r
						return true
					default:
						return false
					}
				})
				var err error
				switch op {
				case '+':
					finfo.minVersion, err = ParseVersion(tokens[0])
				case '-':
					finfo.maxVersion, err = ParseVersion(tokens[0])
				case '<':
					finfo.minVersion, err = ParseVersion(tokens[0])
					if err == nil {
						finfo.maxVersion, err = ParseVersion(tokens[1])
					}
				default:
					finfo.minVersion, err = ParseVersion(flag)
					if err == nil {
						finfo.maxVersion, err = ParseVersion(flag)
					}
				}

				if err != nil {
					return nil, fmt.Errorf("invalid %s version on field %s of type %s (%q): %v", ns, f.Name, typ, f.Tag.Get(ns), err)
				}
			}
		}

		// When any flag except `attr` is used it defaults to `element`
		if finfo.flags&fAttr == 0 {
			finfo.flags |= fElement
		}
	}

	if tag != "" {
		finfo.name = tag
	} else {
		// Use field name as default.
		finfo.name = f.Name
	}

	// add static type info about interfaces the type implements
	if f.Type.Implements(arrayType) {
		finfo.flags |= fArray
	}
	if f.Type.Implements(binaryUnmarshalerType) {
		finfo.flags |= fBinaryUnmarshal
	}
	if f.Type.Implements(binaryMarshalerType) {
		finfo.flags |= fBinaryMarshal
	}
	if f.Type.Implements(textUnmarshalerType) {
		finfo.flags |= fTextUnmarshal
	}
	if f.Type.Implements(textMarshalerType) {
		finfo.flags |= fTextMarshal
	}
	if f.Type.Implements(unmarshalerType) {
		finfo.flags |= fUnmarshal
	}
	if f.Type.Implements(marshalerType) {
		finfo.flags |= fMarshal
	}
	if f.Type.Implements(attrUnmarshalerType) {
		finfo.flags |= fUnmarshalAttr
	}
	if f.Type.Implements(attrMarshalerType) {
		finfo.flags |= fMarshalAttr
	}

	return finfo, nil
}

func addFieldInfo(typ reflect.Type, tinfo *typeInfo, newf *fieldInfo, ns string) error {
	var conflicts []int
	// Find all conflicts.
	for i := range tinfo.fields {
		oldf := &tinfo.fields[i]

		// Same name is a conflict unless versions don't overlap.
		if newf.name == oldf.name {
			if !newf.minVersion.Between(oldf.minVersion, oldf.maxVersion) {
				continue
			}
			if !newf.maxVersion.Between(oldf.minVersion, oldf.maxVersion) {
				continue
			}
			conflicts = append(conflicts, i)
		}
	}

	// Return the first error.
	for _, i := range conflicts {
		oldf := &tinfo.fields[i]
		f1 := typ.FieldByIndex(oldf.idx)
		f2 := typ.FieldByIndex(newf.idx)
		return fmt.Errorf("xmp: %s field %q with tag %q conflicts with field %q with tag %q", typ, f1.Name, f1.Tag.Get(ns), f2.Name, f2.Tag.Get(ns))
	}

	// Without conflicts, add the new field and return.
	tinfo.fields = append(tinfo.fields, *newf)
	return nil
}

// value returns v's field value corresponding to finfo.
// It's equivalent to v.FieldByIndex(finfo.idx), but initializes
// and dereferences pointers as necessary.
func (finfo *fieldInfo) value(v reflect.Value) reflect.Value {
	for i, x := range finfo.idx {
		if i > 0 {
			t := v.Type()
			if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}

	return v
}

// Load value from interface, but only if the result will be
// usefully addressable.
func derefIndirect(v interface{}) reflect.Value {
	return derefValue(reflect.ValueOf(v))
}

func derefValue(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Interface && !val.IsNil() {
		e := val.Elem()
		if e.Kind() == reflect.Ptr && !e.IsNil() {
			val = e
		}
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	return val
}
