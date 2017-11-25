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

package main

import (
	"fmt"
	"testing"

	_ "trimmer.io/go-xmp/models"
	"trimmer.io/go-xmp/xmp"
)

type TestModel struct {
	S1    string                  `test:"s1"      xmp:"test:s1"`
	T1    TestType1               `test:"t1"      xmp:"test:t1"`      // Text (Un)Marshaler
	T2    TestType2               `test:"t2"      xmp:"test:t2"`      // XMP (Un)Marshaler
	T3    TestType3               `test:"t3"      xmp:"test:t3"`      // Binary Unmarshaler
	A1    TestType4               `test:"a1"      xmp:"test:a1,attr"` // XMP Attr Marshaler
	None  string                  `test:"-"       xmp:"-"`            // unexported
	Empty string                  `test:"e"       xmp:"test:e,empty"` // export even if empty
	Ext   xmp.NamedExtensionArray `test:"ext,any" xmp:"test:ext,any"` // any flag for ext
}

var NsTest = xmp.NewNamespace("test", "http://ns.example.com/test/1.0/", NewTestModel)

func init() {
	xmp.Register(NsTest)
}

func NewTestModel(name string) xmp.Model {
	return &TestModel{}
}

func (m *TestModel) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsTest}
}

func (m *TestModel) Can(nsName string) bool {
	return nsName == NsTest.GetName()
}

func (x *TestModel) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *TestModel) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x TestModel) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *TestModel) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *TestModel) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsTest.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *TestModel) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsTest.GetName(), err)
	}
	return nil
}

type TestType1 string
type TestType2 string
type TestType3 string
type TestType4 string

func (t *TestType1) UnmarshalText(b []byte) error {
	s := string(b)
	switch s {
	case "error":
		return fmt.Errorf("unmarshal text error")
	case "x":
		*t = "y"
	default:
		*t = TestType1(s)
	}
	return nil
}

func (t *TestType2) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return nil
}

func (t *TestType3) UnmarshalBinary(b []byte) error {
	return nil
}

func (t *TestType4) UnmarshalAttr(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return nil
}

// use `any` field of a model that supports it
func TestPathSetAny(T *testing.T) {
}

func TestPathGetAny(T *testing.T) {
}

func TestPathSetNamedExtension(T *testing.T) {
}

func TestPathGetNamedExtension(T *testing.T) {
}

func TestStructTagMinus(T *testing.T) {}
func TestStructTagEmpty(T *testing.T) {}
func TestStructTagAttr(T *testing.T)  {}

func TestSetNativeTag(T *testing.T) {}
func TestGetNativeTag(T *testing.T) {}

func TestTextMarshaler(T *testing.T)      {}
func TestTextUnmarshaler(T *testing.T)    {}
func TestXmpAttrMarshaler(T *testing.T)   {}
func TestXmpAttrUnmarshaler(T *testing.T) {}
func TestXmpMarshaler(T *testing.T)       {}
func TestXmpUnmarshaler(T *testing.T)     {}
func TestBinaryUnmarshaler(T *testing.T) {
	// m.SetTag()
}
