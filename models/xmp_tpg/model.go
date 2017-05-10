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
//

// Package xmptpg implements the XMP Paged Text namespace as defined by XMP Specification Part 2.
package xmptpg

import (
	"fmt"
	"github.com/echa/go-xmp/xmp"
)

var (
	NsXmpTPg = xmp.NewNamespace("xmpTPg", "http://ns.adobe.com/xap/1.0/t/pg/", NewModel)
	nsStDim  = xmp.NewNamespace("stDim", "http://ns.adobe.com/xap/1.0/sType/Dimensions#", nil)
	nsStFnt  = xmp.NewNamespace("stFnt", "http://ns.adobe.com/xap/1.0/sType/Font#", nil)
)

func init() {
	xmp.Register(NsXmpTPg, xmp.XmpMetadata)
	xmp.Register(nsStDim)
	xmp.Register(nsStFnt)
}

func NewModel(name string) xmp.Model {
	return &PagedText{}
}

func MakeModel(d *xmp.Document) (*PagedText, error) {
	m, err := d.MakeModel(NsXmpTPg)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*PagedText)
	return x, nil
}

func FindModel(d *xmp.Document) *PagedText {
	if m := d.FindModel(NsXmpTPg); m != nil {
		return m.(*PagedText)
	}
	return nil
}

type PagedText struct {
	Colorants   ColorantArray   `xmp:"xmpTPg:Colorants"`
	Fonts       FontArray       `xmp:"xmpTPg:Fonts"`
	MaxPageSize Dimensions      `xmp:"xmpTPg:MaxPageSize"`
	NPages      int64           `xmp:"xmpTPg:NPages,attr"`
	PlateNames  xmp.StringArray `xmp:"xmpTPg:PlateNames"`
}

func (x PagedText) Can(nsName string) bool {
	return NsXmpTPg.GetName() == nsName
}

func (x PagedText) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsXmpTPg}
}

func (x *PagedText) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x PagedText) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *PagedText) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *PagedText) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsXmpTPg.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *PagedText) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsXmpTPg.GetName(), err)
	}
	return nil
}
