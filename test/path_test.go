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
	"testing"

	_ "github.com/echa/go-xmp/models"
	"github.com/echa/go-xmp/models/dc"
	"github.com/echa/go-xmp/models/xmp_mm"
	"github.com/echa/go-xmp/xmp"
)

func TestPathGet(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Format: "teststring",
	}
	d.AddModel(c)
	if v, err := d.GetPath(xmp.Path("dc:format")); err != nil {
		T.Error(err)
	} else if v != "teststring" {
		T.Errorf("invalid result: expected=teststring got=%s", v)
	}
}

func TestPathGetSlice(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if v, err := d.GetPath(xmp.Path("dc:type")); err != nil {
		T.Error(err)
	} else if v != "one" {
		T.Errorf("invalid result: expected=one got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:type[]")); err != nil {
		T.Error(err)
	} else if v != "one" {
		T.Errorf("invalid result: expected=one got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:type[0]")); err != nil {
		T.Error(err)
	} else if v != "one" {
		T.Errorf("invalid result: expected=one got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:type[1]")); err != nil {
		T.Error(err)
	} else if v != "two" {
		T.Errorf("invalid result: expected=two got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:type[3]")); err != nil {
		T.Error(err)
	} else if v != "" {
		T.Errorf("invalid result: expected=<empty> got=%s", v)
	}
}

func TestPathGetAltString(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Description: xmp.NewAltString(
			xmp.AltItem{Value: "german", Lang: "de", IsDefault: false},
			xmp.AltItem{Value: "english", Lang: "en", IsDefault: true},
		),
	}
	d.AddModel(c)
	if v, err := d.GetPath(xmp.Path("dc:description")); err != nil {
		T.Error(err)
	} else if v != "english" {
		T.Errorf("invalid result: expected=english got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:description[]")); err != nil {
		T.Error(err)
	} else if v != "english" {
		T.Errorf("invalid result: expected=english got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:description[en]")); err != nil {
		T.Error(err)
	} else if v != "english" {
		T.Errorf("invalid result: expected=english got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:description[de]")); err != nil {
		T.Error(err)
	} else if v != "german" {
		T.Errorf("invalid result: expected=german got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("dc:description[fr]")); err != nil {
		T.Error(err)
	} else if v != "" {
		T.Errorf("invalid result: expected=<empty> got=%s", v)
	}
	// not sure if this is a nuonce of the implementation
	if v, err := d.GetPath(xmp.Path("dc:description[0]")); err != nil {
		T.Error(err)
	} else if v != "english" {
		T.Errorf("invalid result: expected=english got=%s", v)
	}
}

func TestPathGetNestedSlice(T *testing.T) {
	d := xmp.NewDocument()
	mm := &xmpmm.XmpMM{
		Ingredients: xmpmm.ResourceRefArray{
			&xmpmm.ResourceRef{AlternatePaths: xmp.UriArray{xmp.NewUri("x")}},
		},
	}
	d.AddModel(mm)
	if v, err := d.GetPath(xmp.Path("xmpMM:Ingredients[0]/alternatePaths")); err != nil {
		T.Error(err)
	} else if v != "x" {
		T.Errorf("invalid result: expected=x got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("xmpMM:Ingredients[0]/alternatePaths[]")); err != nil {
		T.Error(err)
	} else if v != "x" {
		T.Errorf("invalid result: expected=x got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("xmpMM:Ingredients[0]/alternatePaths[0]")); err != nil {
		T.Error(err)
	} else if v != "x" {
		T.Errorf("invalid result: expected=x got=%s", v)
	}
	if v, err := d.GetPath(xmp.Path("xmpMM:Ingredients[0]/alternatePaths[1]")); err != nil {
		T.Error(err)
	} else if v != "" {
		T.Errorf("invalid result: expected=<empty> got=%s", v)
	}
}

func TestPathGetExtension(T *testing.T) {
	d1 := xmp.NewDocument()
	c := &dc.DublinCore{
		Format: "x",
	}
	d1.AddModel(c)
	d2 := xmp.NewDocument()
	mm, err := xmpmm.MakeModel(d2)
	if err != nil {
		T.Errorf("error creating xmpMM model: %v", err)
		return
	}
	mm.AddPantry(d1)
	if l := len(mm.Pantry); l != 1 {
		T.Errorf("invalid pantry len: expected=1 got=%d", l)
	}
	if v, err := d2.GetPath(xmp.Path("xmpMM:Pantry[0]/dc:format")); err != nil {
		T.Error(err)
	} else if v != "x" {
		T.Errorf("invalid result: expected=x got=%s", v)
	}
}

func TestPathSet(T *testing.T) {
	d := xmp.NewDocument()
	p := xmp.Path("dc:format")
	if err := d.SetPath(xmp.PathValue{
		Path:  p,
		Value: "y",
		Flags: xmp.DEFAULT,
	}); err != nil {
		T.Errorf("%v: set failed: %v", p, err)
	}
	if v, err := d.GetPath(p); err != nil {
		T.Error(err)
	} else if v != "y" {
		T.Errorf("invalid result: expected=x got=%s", v)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("missing dc model")
	} else if v.Format != "y" {
		T.Errorf("invalid result: expected=x got=%s", v.Format)
	}
}

func TestPathSetSlice(T *testing.T) {
	d := xmp.NewDocument()
	p := xmp.Path("dc:type")
	if err := d.SetPath(xmp.PathValue{
		Path:  p,
		Value: "y",
		Flags: xmp.DEFAULT,
	}); err != nil {
		T.Errorf("%v: set failed: %v", p, err)
	}
	if v, err := d.GetPath(p); err != nil {
		T.Error(err)
	} else if v != "y" {
		T.Errorf("invalid result: expected=y got=%s", v)
	}
	v := dc.FindModel(d)
	if v == nil {
		T.Errorf("missing dc model")
		return
	}
	if l := len(v.Type); l != 1 {
		T.Errorf("invalid slice length: expected=1 got=%d", l)
	}
	if s := v.Type[0]; s != "y" {
		T.Errorf("invalid result: expected=y got=%s", s)
	}
}

func TestPathSetAltString(T *testing.T) {
	d := xmp.NewDocument()
	p := xmp.Path("dc:description[en]")
	if err := d.SetPath(xmp.PathValue{
		Path:  p,
		Value: "y",
		Flags: xmp.DEFAULT,
	}); err != nil {
		T.Errorf("%v: set failed: %v", p, err)
	}
	if v, err := d.GetPath(p); err != nil {
		T.Error(err)
	} else if v != "y" {
		T.Errorf("invalid result: expected=y got=%s", v)
	}
	v := dc.FindModel(d)
	if v == nil {
		T.Errorf("missing dc model")
		return
	}
	if l := len(v.Description); l != 1 {
		T.Errorf("invalid slice length: expected=1 got=%d", l)
	}
	if s := v.Description.Get("en"); s != "y" {
		T.Errorf("invalid result: expected=y got=%s", s)
	}
	if s := v.Description.Default(); s != "y" {
		T.Errorf("invalid result: expected=y got=%s", s)
	}
}

func TestPathSetNestedSlice(T *testing.T) {
	d := xmp.NewDocument()
	p := xmp.Path("xmpMM:Ingredients[0]/alternatePaths")
	if err := d.SetPath(xmp.PathValue{
		Path:  p,
		Value: "y",
		Flags: xmp.DEFAULT,
	}); err != nil {
		T.Errorf("%v: set failed: %v", p, err)
	}
	if v, err := d.GetPath(p); err != nil {
		T.Error(err)
	} else if v != "y" {
		T.Errorf("invalid result: expected=y got=%s", v)
	}
}

func TestPathSetExtension(T *testing.T) {
	d1 := xmp.NewDocument()
	c := &dc.DublinCore{
		Format: "x",
	}
	d1.AddModel(c)
	d2 := xmp.NewDocument()
	mm, err := xmpmm.MakeModel(d2)
	if err != nil {
		T.Errorf("error creating xmpMM model: %v", err)
		return
	}
	mm.AddPantry(d1)
	p := xmp.Path("xmpMM:Pantry[0]/dc:format")
	if err := d2.SetPath(xmp.PathValue{
		Path:  p,
		Value: "y",
		Flags: xmp.REPLACE,
	}); err != nil {
		T.Errorf("%v: set failed: %v", p, err)
	}
	if val, err := d2.GetPath(p); err != nil {
		T.Errorf("%v: get failed: %v", p, err)
	} else if val != "y" {
		T.Errorf("%v: invalid contents, expected=y got=%s", p, val)
	}
}

func TestPathCreateExtension(T *testing.T) {
	d := xmp.NewDocument()
	p := xmp.Path("xmpMM:Pantry[0]/dc:format")
	if err := d.SetPath(xmp.PathValue{
		Path:  p,
		Value: "y",
		Flags: xmp.CREATE,
	}); err != nil {
		T.Errorf("%v: set failed: %v", p, err)
	}
	if val, err := d.GetPath(p); err != nil {
		T.Errorf("%v: get failed: %v", p, err)
	} else if val != "y" {
		T.Errorf("%v: invalid contents, expected=y got=%s", p, val)
	}
}

func TestPathCreate(T *testing.T) {
	d1 := xmp.NewDocument()
	for n, v := range SyncMergeTestcases {
		if err := d1.SetPath(xmp.PathValue{
			Path:      n,
			Value:     v,
			Namespace: "exotic:ns:42",
			Flags:     xmp.CREATE,
		}); err != nil {
			T.Errorf("%v: %v", n, err)
		}
	}
	d2 := xmp.NewDocument()
	if err := d2.Merge(d1, xmp.MERGE); err != nil {
		T.Errorf("merge failed: %v", err)
	}
	for n, v := range SyncMergeTestcases {
		if val, err := d2.GetPath(n); err != nil {
			T.Errorf("%v: get failed: %v", n, err)
		} else if v != val {
			T.Errorf("%v: invalid contents, expected=%s got=%s", n, v, val)
		}
	}
}

func TestPathCreateModelNamespace(T *testing.T) {
	d := xmp.NewDocument()
	if v := dc.FindModel(d); v != nil {
		T.Errorf("expected no model to exist")
	}
	if err := d.SetPath(xmp.PathValue{Path: xmp.Path("dc:"),
		Value: "-",
		Flags: xmp.CREATE,
	}); err != nil {
		T.Errorf("error creating model: %v", err)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("expected model to exist")
	}
}

func TestPathCreateCustomNamespace(T *testing.T) {
	d := xmp.NewDocument()
	p := xmp.Path("myCustomNs:")
	if _, err := d.GetPath(p); err == nil {
		T.Errorf("expected no path to exist")
	}
	if err := d.SetPath(xmp.PathValue{
		Path:      p,
		Namespace: "ns.example.com/",
		Value:     "-",
		Flags:     xmp.CREATE,
	}); err != nil {
		T.Errorf("error creating non-model path: %v", err)
	}
	if _, err := d.GetPath(p); err != nil {
		T.Errorf("expected path to exist")
	}
}

func TestPathReplace(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Format: "old",
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:format"),
		Value: "new",
		Flags: xmp.REPLACE,
	}); err != nil {
		T.Errorf("path set failed: %v", err)
	}
	if c.Format != "new" {
		T.Errorf("invalid result: expected=new got=%s", c.Format)
	}
}

func TestPathNotReplace(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Format: "old",
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:format"),
		Value: "new",
		Flags: xmp.CREATE,
	}); err == nil {
		T.Errorf("unexpected nil error when overwriting existing path")
	}
	if c.Format != "old" {
		T.Errorf("invalid result: expected=old got=%s", c.Format)
	}
}

func TestPathReplaceSlice(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type"),
		Value: "four",
		Flags: xmp.REPLACE,
	}); err != nil {
		T.Errorf("unexpected error on append to existing path: %v", err)
	}
	if l := len(c.Type); l != 1 {
		T.Errorf("invalid len: expected=1 got=%d", l)
	}
	if c.Type[0] != "four" {
		T.Errorf("invalid result: expected=four got=%s", c.Type[0])
	}
}

func TestPathAppendExist(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type"),
		Value: "three",
		Flags: xmp.APPEND,
	}); err != nil {
		T.Errorf("unexpected error on append to existing path: %v", err)
	}
	l := len(c.Type)
	if l != 4 {
		T.Errorf("invalid len: expected=4 got=%d", l)
	}
	if c.Type[l-1] != "three" {
		T.Errorf("invalid result: expected=three got=%s", c.Type[l-1])
	}
}

func TestPathAppendNotExist(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type"),
		Value: "three",
		Flags: xmp.APPEND | xmp.CREATE,
	}); err != nil {
		T.Errorf("expected no error on append+create to non-existing path: %v", err)
	}
	l := len(c.Type)
	if l != 1 {
		T.Errorf("invalid len: expected=1 got=%d", l)
	}
	if c.Type[0] != "three" {
		T.Errorf("invalid result: expected=three got=%s", c.Type[0])
	}
}

func TestPathAppendAltExist(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Description: xmp.NewAltString(
			xmp.AltItem{Value: "german", Lang: "de", IsDefault: false},
			xmp.AltItem{Value: "english", Lang: "en", IsDefault: true},
		),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:description[fi]"),
		Value: "finnish",
		Flags: xmp.APPEND,
	}); err != nil {
		T.Errorf("unexpected error on append to existing path: %v", err)
	}
	l := len(c.Description)
	if l != 3 {
		T.Errorf("invalid len: expected=3 got=%d", l)
	}
	if s := c.Description.Get("fi"); s != "finnish" {
		T.Errorf("invalid result: expected=finnish got=%s", s)
	}
	if s := c.Description.Default(); s == "" {
		T.Errorf("invalid result: expected not to be default")
	}
}

func TestPathAppendAltNotExist(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:description[fi]"),
		Value: "finnish",
		Flags: xmp.APPEND | xmp.CREATE,
	}); err != nil {
		T.Errorf("unexpected error on append+create to non-existing path: %v", err)
	}
	l := len(c.Description)
	if l != 1 {
		T.Errorf("invalid len: expected=1 got=%d", l)
	}
	if s := c.Description.Get("fi"); s != "finnish" {
		T.Errorf("invalid result: expected=finnish got=%s", s)
	}
	if s := c.Description.Default(); s == "" {
		T.Errorf("invalid result: expected default")
	}
}

func TestPathAppendUnique(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type"),
		Value: "one",
		Flags: xmp.APPEND | xmp.UNIQUE,
	}); err != nil {
		T.Errorf("unexpected error on append unique to existing path: %v", err)
	}
	l := len(c.Type)
	if l != 3 {
		T.Errorf("invalid len: expected=3 got=%d", l)
	}
	if c.Type[l-1] != "three" {
		T.Errorf("invalid result: expected=three got=%s", c.Type[l-1])
	}
}

func TestPathAppendUniqueAlt(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Description: xmp.NewAltString(
			xmp.AltItem{Value: "german", Lang: "de", IsDefault: false},
			xmp.AltItem{Value: "english", Lang: "en", IsDefault: true},
		),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:description[de]"),
		Value: "german",
		Flags: xmp.APPEND | xmp.UNIQUE,
	}); err != nil {
		T.Errorf("unexpected error on append unique to existing path: %v", err)
	}
	l := len(c.Description)
	if l != 2 {
		T.Errorf("invalid len: expected=2 got=%d", l)
	}
}

func TestPathDeleteWithoutFlag(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Format: "old",
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:format"),
		Value: "",
		Flags: xmp.CREATE,
	}); err == nil {
		T.Errorf("unexpected nil error when deleting without flag")
	}
	if c.Format == "" {
		T.Errorf("invalid result: expected=old got=%s", c.Format)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
}

func TestPathDeleteNonExist(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Description: xmp.NewAltString(
			xmp.AltItem{Value: "german", Lang: "de", IsDefault: false},
			xmp.AltItem{Value: "english", Lang: "en", IsDefault: true},
		),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:format"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing path: %v", err)
	}
	if len(c.Description) == 0 {
		T.Errorf("unexpected empty AltString")
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
}

func TestPathDelete(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Format: "old",
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:format"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing path: %v", err)
	}
	if c.Format != "" {
		T.Errorf("invalid result: expected=<empty> got=%s", c.Format)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
}

func TestPathDeleteSlice(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing path: %v", err)
	}
	if len(c.Type) != 0 {
		T.Errorf("unexpected non-empty array after delete")
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
}

func TestPathDeleteAltString(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Description: xmp.NewAltString(
			xmp.AltItem{Value: "german", Lang: "de", IsDefault: false},
			xmp.AltItem{Value: "english", Lang: "en", IsDefault: true},
		),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:description"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing path: %v", err)
	}
	if len(c.Description) != 0 {
		T.Errorf("unexpected non-empty array after delete")
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
}

func TestPathDeleteFromSlice1(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type[0]"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing slice element: %v", err)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
	if l := len(c.Type); l != 2 {
		T.Errorf("unexpected slice length after delete: expected=2 got=%d", l)
		return
	}
	if s := c.Type[0]; s != "two" {
		T.Errorf("unexpected first slice element after delete: expected=two got=%s", s)
	}
	if s := c.Type[1]; s != "three" {
		T.Errorf("unexpected last slice element after delete: expected=three got=%s", s)
	}
}

func TestPathDeleteFromSlice2(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type[1]"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing slice element: %v", err)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
	if l := len(c.Type); l != 2 {
		T.Errorf("unexpected slice length after delete: expected=2 got=%d", l)
		return
	}
	if s := c.Type[0]; s != "one" {
		T.Errorf("unexpected first slice element after delete: expected=one got=%s", s)
	}
	if s := c.Type[1]; s != "three" {
		T.Errorf("unexpected last slice element after delete: expected=three got=%s", s)
	}
}

func TestPathDeleteFromSlice3(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type[2]"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing slice element: %v", err)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
	if l := len(c.Type); l != 2 {
		T.Errorf("unexpected slice length after delete: expected=2 got=%d", l)
		return
	}
	if s := c.Type[0]; s != "one" {
		T.Errorf("unexpected first slice element after delete: expected=one got=%s", s)
	}
	if s := c.Type[1]; s != "two" {
		T.Errorf("unexpected last slice element after delete: expected=two got=%s", s)
	}
}

func TestPathDeleteFromSliceNotExist(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Type: xmp.NewStringArray("one", "two", "three"),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:type[3]"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting non-existing slice element: %v", err)
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
	if l := len(c.Type); l != 3 {
		T.Errorf("unexpected slice length after delete: expected=3 got=%d", l)
		return
	}
	if s := c.Type[0]; s != "one" {
		T.Errorf("unexpected first slice element after delete: expected=one got=%s", s)
	}
	if s := c.Type[1]; s != "two" {
		T.Errorf("unexpected first slice element after delete: expected=two got=%s", s)
	}
	if s := c.Type[2]; s != "three" {
		T.Errorf("unexpected last slice element after delete: expected=three got=%s", s)
	}
}

func TestPathDeleteFromAltString(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Description: xmp.NewAltString(
			xmp.AltItem{Value: "german", Lang: "de", IsDefault: false},
			xmp.AltItem{Value: "english", Lang: "en", IsDefault: true},
		),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:description[de]"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing path: %v", err)
	}
	if len(c.Description) != 1 {
		T.Errorf("unexpected non-empty AltString after delete")
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
	if s := c.Description.Get("de"); s != "" {
		T.Errorf("invalid result: expected=<empty> got=%s", s)
	}
	if s := c.Description.Get("en"); s != "english" {
		T.Errorf("invalid result: expected=english got=%s", s)
	}
	if s := c.Description.Default(); s != "english" {
		T.Errorf("invalid default: expected=english got=%s", s)
	}
}

func TestPathDeleteFromAltStringDefault(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{
		Description: xmp.NewAltString(
			xmp.AltItem{Value: "german", Lang: "de", IsDefault: false},
			xmp.AltItem{Value: "english", Lang: "en", IsDefault: true},
		),
	}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:description[en]"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error when deleting existing path: %v", err)
	}
	if len(c.Description) != 1 {
		T.Errorf("unexpected non-empty AltString after delete")
	}
	if v := dc.FindModel(d); v == nil {
		T.Errorf("unexpected empty model after delete")
	}
	if s := c.Description.Get("en"); s != "" {
		T.Errorf("invalid result: expected=<empty> got=%s", s)
	}
	if s := c.Description.Get("de"); s != "german" {
		T.Errorf("invalid result: expected=german got=%s", s)
	}
	if s := c.Description.Default(); s != "german" {
		T.Errorf("invalid default: expected=german got=%s", s)
	}
}

func TestPathDeleteNamespace(T *testing.T) {
	d := xmp.NewDocument()
	c := &dc.DublinCore{}
	d.AddModel(c)
	if err := d.SetPath(xmp.PathValue{
		Path:  xmp.Path("dc:"),
		Value: "",
		Flags: xmp.DELETE,
	}); err != nil {
		T.Errorf("unexpected error on delete existing namespace: %v", err)
	}
	if v := dc.FindModel(d); v != nil {
		T.Errorf("expected no model to exist")
	}
}
