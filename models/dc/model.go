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

// Package dc implements Dublin Core metadata as defined in XMP Specification Part 1.
package dc

import (
	"fmt"
	"github.com/echa/go-xmp/xmp"
)

var (
	NsDc = xmp.NewNamespace("dc", "http://purl.org/dc/elements/1.1/", NewModel)
)

func init() {
	xmp.Register(NsDc, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &DublinCore{}
}

func MakeModel(d *xmp.Document) (*DublinCore, error) {
	m, err := d.MakeModel(NsDc)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*DublinCore)
	return x, nil
}

func FindModel(d *xmp.Document) *DublinCore {
	if m := d.FindModel(NsDc); m != nil {
		return m.(*DublinCore)
	}
	return nil
}

// Note: differences from vanilla Dublin Core are
// - Date is an array instead of a simple string
// - Kanguage is an array instead of a simple string
type DublinCore struct {
	Contributor xmp.StringArray `xmp:"dc:contributor"`
	Coverage    string          `xmp:"dc:coverage"`
	Creator     xmp.StringList  `xmp:"dc:creator"`
	Date        xmp.DateList    `xmp:"dc:date"`
	Description xmp.AltString   `xmp:"dc:description"`
	Format      string          `xmp:"dc:format"`
	Identifier  string          `xmp:"dc:identifier"`
	Language    LocaleArray     `xmp:"dc:language"`
	Publisher   xmp.StringArray `xmp:"dc:publisher"`
	Relation    xmp.StringArray `xmp:"dc:relation"`
	Rights      xmp.AltString   `xmp:"dc:rights"`
	Source      string          `xmp:"dc:source"`
	Subject     xmp.StringArray `xmp:"dc:subject"`
	Title       xmp.AltString   `xmp:"dc:title"`
	Type        xmp.StringArray `xmp:"dc:type"`
}

func (x DublinCore) Can(nsName string) bool {
	return NsDc.GetName() == nsName
}

func (x DublinCore) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsDc}
}

func (x *DublinCore) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x DublinCore) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *DublinCore) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *DublinCore) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsDc.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *DublinCore) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsDc.GetName(), err)
	}
	return nil
}
