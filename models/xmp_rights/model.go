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

// Package xmprights implements the XMP Rights management namespace as defined by XMP Specification Part 1.
package xmprights

import (
	"fmt"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsXmpRights = xmp.NewNamespace("xmpRights", "http://ns.adobe.com/xap/1.0/rights/", NewModel)
)

func init() {
	xmp.Register(NsXmpRights, xmp.XmpMetadata, xmp.RightsMetadata)
}

func NewModel(name string) xmp.Model {
	return &XmpRights{}
}

func MakeModel(d *xmp.Document) (*XmpRights, error) {
	m, err := d.MakeModel(NsXmpRights)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*XmpRights)
	return x, nil
}

func FindModel(d *xmp.Document) *XmpRights {
	if m := d.FindModel(NsXmpRights); m != nil {
		return m.(*XmpRights)
	}
	return nil
}

type XmpRights struct {
	Certificate  string          `xmp:"xmpRights:Certificate"`
	Marked       xmp.Bool        `xmp:"xmpRights:Marked"`
	Owner        xmp.StringArray `xmp:"xmpRights:Owner"`
	UsageTerms   xmp.AltString   `xmp:"xmpRights:UsageTerms"`
	WebStatement string          `xmp:"xmpRights:WebStatement"`
}

func (x XmpRights) Can(nsName string) bool {
	return NsXmpRights.GetName() == nsName
}

func (x XmpRights) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsXmpRights}
}

func (x *XmpRights) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *XmpRights) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x XmpRights) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *XmpRights) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *XmpRights) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsXmpRights.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *XmpRights) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsXmpRights.GetName(), err)
	}
	return nil
}
