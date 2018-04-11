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

// Types as defined in ISO 16684-1:2011(E) 8.2.1 (Core value types)
// - IntArray (ordered)
// - StringList (ordered)
// - StringArray (unordered)
// - AltString (string, xml:lang support)

package xmp

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Array interface {
	Typ() ArrayType
}

// Choice / Alternative Arrays with xml:lang support
//
type ArrayType string

const (
	ArrayTypeOrdered     ArrayType = "Seq"
	ArrayTypeUnordered   ArrayType = "Bag"
	ArrayTypeAlternative ArrayType = "Alt"
)

// Note: when changing these tags, also update json.go for proper marshal/demarshal
type AltItem struct {
	Value     string `xmp:",chardata" json:"value"`
	Lang      string `xmp:"xml:lang"  json:"lang"`
	IsDefault bool   `xmp:"-"         json:"isDefault"`
}

func (x *AltItem) UnmarshalText(data []byte) error {
	x.Value = string(data)
	x.Lang = ""
	x.IsDefault = true
	return nil
}

func (x AltItem) GetLang() string {
	if x.Lang == "" && x.IsDefault {
		return "x-default"
	}
	return x.Lang
}

type AltString []AltItem

func (x AltString) IsZero() bool {
	return len(x) == 0
}

func NewAltString(items ...interface{}) AltString {
	if len(items) == 0 {
		return nil
	}
	a := make(AltString, 0)
	for _, v := range items {
		switch val := v.(type) {
		case AltItem:
			if val.Value != "" {
				a = append(a, val)
			}
		case string:
			if val != "" {
				a = append(a, AltItem{Value: val})
			}
		case fmt.Stringer:
			if s := val.String(); s != "" {
				a = append(a, AltItem{Value: s})
			}
		}
	}
	a.EnsureDefault()
	return a
}

// make sure there is exactly one a default and it's the first item
func (a *AltString) EnsureDefault() {
	l := len(*a)
	if l > 1 {
		idx := -1
		for i := 0; i < l; i++ {
			if (*a)[i].IsDefault {
				if idx > -1 {
					(*a)[i].IsDefault = false
				} else {
					idx = i
				}
			}
		}
		if idx != 0 {
			(*a)[idx], (*a)[0] = (*a)[0], (*a)[idx]
		}
	} else if l == 1 {
		(*a)[0].IsDefault = true
	}
}

func (a AltString) Default() string {
	for _, v := range a {
		if v.IsDefault {
			return v.Value
		}
	}
	return ""
}

func (a AltString) Index(lang string) int {
	for i, v := range a {
		if v.Lang == lang {
			return i
		}
	}
	return -1
}

func (a AltString) Get(lang string) string {
	if lang == "" {
		return a.Default()
	}
	for _, v := range a {
		if v.Lang == lang {
			return v.Value
		}
	}
	return ""
}

func (a *AltString) AddDefault(lang string, value string) {
	if value == "" {
		return
	}
	i := AltItem{
		Value:     value,
		Lang:      lang,
		IsDefault: true,
	}

	// clear any previous default
	for i, l := 0, len(*a); i < l; i++ {
		(*a)[i].IsDefault = false
	}

	// add new default as first element
	*a = append(AltString{i}, (*a)...)
}

func (a *AltString) AddUnique(lang string, value string) bool {
	if value == "" {
		return false
	}
	if idx := a.Index(lang); idx > -1 && (*a)[idx].Value == value {
		return false
	}
	a.Add(lang, value)
	return true
}

func (a *AltString) Add(lang string, value string) {
	if value == "" {
		return
	}
	*a = append(*a, AltItem{
		Value:     value,
		Lang:      lang,
		IsDefault: false,
	})
	a.EnsureDefault()
}

func (a *AltString) Set(lang string, value string) {
	if value == "" {
		return
	}
	if i := a.Index(lang); i > -1 {
		(*a)[i].Value = value
	} else {
		a.Add(lang, value)
	}
}

func (a *AltString) RemoveLang(lang string) {
	idx := -1
	for i, v := range *a {
		if v.Lang == lang {
			idx = i
			break
		}
	}
	if idx > -1 {
		*a = append((*a)[:idx], (*a)[idx+1:]...)
		a.EnsureDefault()
	}
}

func (a AltString) Typ() ArrayType {
	return ArrayTypeAlternative
}

func (x AltString) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *AltString) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	if err := UnmarshalArray(d, node, x.Typ(), x); err != nil {
		return err
	}

	// merge default with matching language, if any
	var dup int = -1
	for i, v := range *x {
		if v.IsDefault {
			for j, vv := range *x {
				if i != j && v.Value == vv.Value {
					(*x)[i].Lang = vv.Lang
					dup = j
					break
				}
			}
			break
		}
	}

	// remove duplicate
	if dup > -1 {
		*x = append((*x)[:dup], (*x)[dup+1:]...)
	}

	return nil
}

type AltStringArray []*AltString

func (x AltStringArray) IsZero() bool {
	return len(x) == 0
}

func (x *AltStringArray) Add(v *AltString) error {
	if v == nil {
		return nil
	}
	*x = append(*x, v)
	return nil
}

func (x AltStringArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func NewAltStringArray(items ...*AltString) AltStringArray {
	x := make(AltStringArray, 0, len(items))
	return append(x, items...)
}

func (x AltStringArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *AltStringArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

// Unordered Integer Arrays
//
type IntArray []int

func (x IntArray) IsZero() bool {
	return len(x) == 0
}

func (x *IntArray) Add(v int) error {
	*x = append(*x, v)
	return nil
}

func (x IntArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func NewIntArray(items ...int) IntArray {
	x := make(IntArray, 0, len(items))
	return append(x, items...)
}

func (x IntArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *IntArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

// Ordered Integer Arrays
//
type IntList []int

func (x IntList) IsZero() bool {
	return len(x) == 0
}

func (x *IntList) Add(v int) error {
	*x = append(*x, v)
	return nil
}

func (x IntList) Typ() ArrayType {
	return ArrayTypeOrdered
}

func NewIntList(items ...int) IntList {
	x := make(IntList, 0, len(items))
	return append(x, items...)
}

func (x IntList) String() string {
	s := make([]string, len(x))
	for i, v := range x {
		s[i] = strconv.Itoa(v)
	}
	return strings.Join(s, " ")
}

func (x IntList) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *IntList) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

// Ordered String Arrays
//
type StringList []string

func (x StringList) IsZero() bool {
	return len(x) == 0
}

func (x *StringList) Add(v string) error {
	if v == "" {
		return nil
	}
	*x = append(*x, v)
	return nil
}

func (x *StringList) AddUnique(v string) error {
	if v == "" {
		return nil
	}
	if !x.Contains(v) {
		return x.Add(v)
	}
	return nil
}

func (x *StringList) Index(val string) int {
	if val == "" {
		return -1
	}
	for i, v := range *x {
		if v == val {
			return i
		}
	}
	return -1
}

func (x *StringList) Contains(v string) bool {
	return x.Index(v) > -1
}

func (x *StringList) Remove(v string) {
	if v == "" {
		return
	}
	if idx := x.Index(v); idx > -1 {
		*x = append((*x)[:idx], (*x)[:idx+1]...)
	}
}

func (x StringList) Typ() ArrayType {
	return ArrayTypeOrdered
}

func NewStringList(items ...string) StringList {
	if len(items) == 0 {
		return nil
	}
	x := make(StringList, 0, len(items))
	for _, v := range items {
		if v == "" {
			continue
		}
		x = append(x, v)
	}
	return x
}

func (x StringList) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *StringList) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

func (x *StringList) UnmarshalText(data []byte) error {
	list := strings.Split(strings.Replace(string(data), "\r\n", "\n", -1), "\n")
	*x = append(*x, list...)
	return nil
}

func (x StringList) MarshalText() ([]byte, error) {
	if len(x) == 0 {
		return nil, nil
	}
	return []byte(strings.Join(x, "\n")), nil
}

// Unordered String Arrays
//
type StringArray []string

func NewStringArray(items ...string) StringArray {
	if len(items) == 0 {
		return nil
	}
	x := make(StringArray, 0, len(items))
	for _, v := range items {
		if v == "" {
			continue
		}
		x = append(x, v)
	}
	return x
}

func (x StringArray) IsZero() bool {
	return len(x) == 0
}

func (x *StringArray) Add(v string) error {
	if v == "" {
		return nil
	}
	*x = append(*x, v)
	return nil
}

func (x *StringArray) AddUnique(v string) error {
	if v == "" {
		return nil
	}
	if !x.Contains(v) {
		*x = append(*x, v)
	}
	return nil
}

func (x *StringArray) Index(val string) int {
	if val == "" {
		return -1
	}
	for i, v := range *x {
		if v == val {
			return i
		}
	}
	return -1
}

func (x *StringArray) Contains(v string) bool {
	return x.Index(v) > -1
}

func (x *StringArray) Remove(v string) {
	if v == "" {
		return
	}
	if idx := x.Index(v); idx > -1 {
		*x = append((*x)[:idx], (*x)[:idx+1]...)
	}
}

func (x StringArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func (x StringArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *StringArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

func (x *StringArray) UnmarshalText(data []byte) error {
	list := strings.Split(strings.Replace(string(data), "\r\n", "\n", -1), "\n")
	*x = append(*x, list...)
	return nil
}

func (x StringArray) MarshalText() ([]byte, error) {
	if len(x) == 0 {
		return nil, nil
	}
	return []byte(strings.Join(x, "\n")), nil
}

func MarshalArray(e *Encoder, node *Node, typ ArrayType, items interface{}) error {

	val := reflect.ValueOf(items)
	kind := val.Kind()

	if kind != reflect.Slice && kind != reflect.Array {
		return fmt.Errorf("xmp: non-slice type passed to array marshal: %v %v", val.Type(), kind)
	}

	if val.Len() == 0 {
		return nil
	}

	// output enclosing array type
	arr := NewNode(xml.Name{Local: "rdf:" + string(typ)})
	node.AddNode(arr)

	// output array elements
	for i, l := 0, val.Len(); i < l; i++ {
		elem := NewNode(xml.Name{Local: "rdf:li"})
		arr.Nodes = append(arr.Nodes, elem)

		v := val.Index(i).Interface()

		// add xml:Lang attribute to alternatives
		if reflect.TypeOf(v) == reflect.TypeOf(AltItem{}) {
			ai := v.(AltItem)
			if typ == ArrayTypeAlternative {
				if ai.IsDefault || val.Len() == 1 {
					elem.AddStringAttr("xml:lang", "x-default")
					if err := e.EncodeElement(ai.Value, elem); err != nil {
						return err
					}
					// skip outputting default items twice when no lang is set
					if ai.Lang == "" {
						continue
					}
				}
				if ai.IsDefault && ai.Lang != "" {
					// add a second array node for a set language
					elem := NewNode(xml.Name{Local: "rdf:li"})
					arr.Nodes = append(arr.Nodes, elem)
					elem.AddStringAttr("xml:lang", ai.Lang)
					if err := e.EncodeElement(ai.Value, elem); err != nil {
						return err
					}
					continue
				}
				// for any non-default items, add just xml:lang
				if ai.Lang != "" {
					elem.AddStringAttr("xml:lang", ai.Lang)
				} else {
					return fmt.Errorf("xmp: language required for alternative array item '%v'", node.FullName())
				}
			}

			// write node contents to tree
			if err := e.EncodeElement(ai.Value, elem); err != nil {
				return err
			}
		} else {
			// write to stream
			if err := e.EncodeElement(v, elem); err != nil {
				return err
			}
		}
	}

	return nil
}

func UnmarshalArray(d *Decoder, node *Node, typ ArrayType, out interface{}) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(out))
	itemType := sliceValue.Type().Elem()

	//
	// LogDebugf("+++ Array start node %s\n", node.FullName())
	//
	if len(node.Nodes) != 1 {
		return fmt.Errorf("xmp: invalid array %s: contains %d nodes", node.FullName(), len(node.Nodes))
	}
	arr := node.Nodes[0]

	switch ArrayType(arr.Name()) {
	default:
		return fmt.Errorf("xmp: invalid array type %s", node.FullName())
	case ArrayTypeOrdered,
		ArrayTypeUnordered,
		ArrayTypeAlternative:
	}

	for i, n := range arr.Nodes {
		if n.FullName() != "rdf:li" {
			return fmt.Errorf("xmp: invalid array element type %s", n.FullName())
		}

		val := reflect.New(itemType)
		if itemType == reflect.TypeOf(AltItem{}) {
			// Special unmarshalling for AltItems with lang attributes
			//
			i := val.Interface().(*AltItem)
			if typ == ArrayTypeAlternative {
				for _, v := range n.GetAttr("", "lang") {
					switch v.Value {
					case "x-default":
						i.IsDefault = true
					default:
						i.Lang = v.Value
					}
				}
			}
			if err := d.DecodeElement(&i.Value, n); err != nil {
				return err
			}
		} else {
			// custom unmarshal for other types
			//
			// LogDebugf("++++ Array unmarshal custom type=%v\n", val.Type())
			//
			if err := d.unmarshal(val.Elem(), nil, n); err != nil {
				return err
			}
		}
		if sliceValue.Kind() == reflect.Array {
			if sliceValue.Type().Len() <= i {
				return fmt.Errorf("xmp: too many elements for fixed array type %s", n.FullName())
			}
			sliceValue.Index(i).Set(val.Elem())
		} else {
			sliceValue.Set(reflect.Append(sliceValue, val.Elem()))
		}
	}

	return nil
}
