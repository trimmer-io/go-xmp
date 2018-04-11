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
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
)

type Tag struct {
	Key   string `xmp:"-" json:"key,omitempty"`
	Value string `xmp:"-" json:"value,omitempty"`
	Lang  string `xmp:"-" json:"lang,omitempty"`
}

type TagList []Tag

func (x TagList) MarshalXMP(e *Encoder, node *Node, m Model) error {
	if len(x) == 0 {
		return nil
	}
	for _, v := range x {
		name := xml.Name{Local: v.Key}
		if v.Lang != "" {
			name.Local += "-" + v.Lang
		}
		node.AddAttr(Attr{
			Name:  name,
			Value: v.Value,
		})
	}
	return nil
}

func (x *TagList) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	for _, v := range node.Attr {
		tag := Tag{
			Key:   v.Name.Local,
			Value: v.Value,
		}
		if i := strings.Index(tag.Key, "-"); i > -1 {
			// FIXME: parsing would be better here
			tag.Key, tag.Lang = tag.Key[:i], tag.Key[i+1:]
		}
		*x = append(*x, tag)
	}
	return nil
}

func GetNativeField(v Model, name string) (string, error) {
	nsName, err := getNamespaceName(v)
	if err != nil {
		return "", err
	}

	val := derefIndirect(v)
	finfo, err := findField(val, name, nsName)
	if err != nil {
		return "", err
	}

	fv := finfo.value(val)
	typ := fv.Type()

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

	// Check for text marshaler and marshal as node value
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

	if fv.CanInterface() && (finfo != nil && finfo.flags&fTextMarshal > 0 || typ.Implements(textMarshalerType)) {
		b, err := fv.Interface().(encoding.TextMarshaler).MarshalText()
		if err != nil || b == nil {
			return "", err
		}
		return string(b), nil
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

func SetNativeField(v Model, name, value string) error {
	nsName, err := getNamespaceName(v)
	if err != nil {
		return err
	}

	val := derefIndirect(v)
	finfo, err := findField(val, name, nsName)
	if err != nil {
		return err
	}

	f := finfo.value(val)

	// allocate memory for pointer values in structs
	if f.Type().Kind() == reflect.Ptr && f.IsNil() && f.CanSet() {
		f.Set(reflect.New(f.Type().Elem()))
	}

	// load and potentially create value
	f = derefValue(f)

	// try unmarshalers
	if f.CanAddr() {
		pv := f.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fBinaryUnmarshal > 0 || pv.Type().Implements(binaryUnmarshalerType)) {
			return pv.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(value))
		}
	}

	if f.CanInterface() && (finfo != nil && finfo.flags&fBinaryUnmarshal > 0 || f.Type().Implements(binaryUnmarshalerType)) {
		return f.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(value))
	}

	if f.CanAddr() {
		pv := f.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fTextUnmarshal > 0 || pv.Type().Implements(textUnmarshalerType)) {
			return pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
		}
	}

	if f.CanInterface() && (finfo != nil && finfo.flags&fTextUnmarshal > 0 || f.Type().Implements(textUnmarshalerType)) {
		return f.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
	}

	// otherwise set simple field value directly or fail
	return setValue(f, value)
}

func SetLocaleField(v Model, lang string, name, value string) error {
	nsName, err := getNamespaceName(v)
	if err != nil {
		return err
	}

	val := derefIndirect(v)
	finfo, err := findField(val, name, nsName)
	if err != nil {
		return err
	}

	f := finfo.value(val)

	// allocate memory for pointer values in structs
	if f.Type().Kind() == reflect.Ptr && f.IsNil() && f.CanSet() {
		f.Set(reflect.New(f.Type().Elem()))
	}

	// load and potentially create value
	f = derefValue(f)

	if f.Kind() != reflect.Slice || f.Type().Elem() != reflect.TypeOf(AltItem{}) {
		return fmt.Errorf("field '%s' must be of type xmp.AltString, found type '%s' kind '%s'", name, f.Type().String(), f.Kind())
	}

	// we need a pointer to AltString slices for appending
	a, ok := f.Addr().Interface().(*AltString)
	if !ok {
		return fmt.Errorf("field '%s' must be of type xmp.AltString", name)
	}

	// use AltString interface to add value
	a.Set(lang, value)
	return nil
}

func GetLocaleField(v Model, lang string, name string) (string, error) {
	nsName, err := getNamespaceName(v)
	if err != nil {
		return "", err
	}

	val := derefIndirect(v)
	finfo, err := findField(val, name, nsName)
	if err != nil {
		return "", err
	}

	f := finfo.value(val)

	if f.Kind() != reflect.Slice && f.Type().Elem() != reflect.TypeOf(AltItem{}) {
		return "", fmt.Errorf("field '%s' must be of type AltString, found %s (%s)", name, f.Type().String(), f.Kind())
	}

	a, ok := f.Interface().(AltString)
	if !ok {
		return "", fmt.Errorf("field '%s' must be of type AltString", name)
	}

	// use AltString interface to get value
	return a.Get(lang), nil
}

func ListNativeFields(v Model) (TagList, error) {
	nsName, err := getNamespaceName(v)
	if err != nil {
		return nil, err
	}

	val := derefIndirect(v)
	typ := val.Type()

	tinfo, err := getTypeInfo(typ, nsName)
	if err != nil {
		return nil, err
	}

	tagList := make(TagList, 0)

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

		// Drill into interfaces and pointers.
		// This can turn into an infinite loop given a cyclic chain,
		// but it matches the Go 1 behavior.
		for fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}

		tag := Tag{
			Key: finfo.name,
		}

		// Check for text marshaler and marshal as node value
		if fv.CanInterface() && typ.Implements(textMarshalerType) {
			b, err := fv.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return nil, err
			}
			if len(b) == 0 {
				continue
			}
			tag.Value = string(b)
			tagList = append(tagList, tag)
			continue
		}

		if fv.CanAddr() {
			pv := fv.Addr()
			if pv.CanInterface() && pv.Type().Implements(textMarshalerType) {
				b, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, err
				}
				if len(b) == 0 {
					continue
				}
				tag.Value = string(b)
				tagList = append(tagList, tag)
				continue
			}
		}

		// handle multi-language arrays
		if fv.Kind() == reflect.Slice && fv.Type().Elem() == reflect.TypeOf(AltItem{}) {
			a, ok := fv.Interface().(AltString)
			if !ok {
				return nil, fmt.Errorf("field '%s' must be of type AltString", finfo.name)
			}
			for _, v := range a {
				tagList = append(tagList, Tag{
					Key:   finfo.name,
					Lang:  v.Lang,
					Value: v.Value,
				})
			}
			continue
		}

		// simple values are just fine, but any other type (slice, array, struct)
		// without textmarshaler will fail here
		if s, b, err := marshalSimple(typ, fv); err != nil {
			return nil, err
		} else {
			if b != nil {
				s = string(b)
			}
			if len(s) == 0 {
				continue
			}
			tag.Value = s
			tagList = append(tagList, tag)
		}
	}

	return tagList, nil
}

func findField(val reflect.Value, name, ns string) (*fieldInfo, error) {
	typ := val.Type()
	tinfo, err := getTypeInfo(typ, ns)
	if err != nil {
		return nil, err
	}

	// pick the correct field based on name, flags and version
	var finfo *fieldInfo
	any := -1
	for i, v := range tinfo.fields {
		// version must always match
		// if !d.version.Between(v.minVersion, v.maxVersion) {
		// 	continue
		// }

		// save `any` field in case
		if v.flags&fAny > 0 {
			any = i
		}

		// field name must match
		if hasPrefix(name) {
			// exact match when namespace is specified
			if v.name != name {
				continue
			}
		} else {
			// suffix match without namespace
			if stripPrefix(v.name) != name {
				continue
			}
		}

		finfo = &v
		break
	}

	if finfo == nil && any >= 0 {
		finfo = &tinfo.fields[any]
	}

	// nothing found
	if finfo == nil {
		return nil, fmt.Errorf("no field with tag '%s' in type '%s'", name, typ.String())
	}

	return finfo, nil
}

func getNamespaceName(v Model) (string, error) {
	if n := v.Namespaces(); len(n) == 0 {
		return "", fmt.Errorf("model '%s' must implement at least one namespace", reflect.TypeOf(v).String())
	} else {
		return n[0].GetName(), nil
	}
}
