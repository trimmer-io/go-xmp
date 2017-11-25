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

// Package pdf implements metadata for PDF files as defined by XMP Specification Part 2.
package pdf

import (
	"fmt"
	"trimmer.io/go-xmp/models/dc"
	"trimmer.io/go-xmp/models/xmp_base"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsPDF = xmp.NewNamespace("pdf", "http://ns.adobe.com/pdf/1.3/", NewModel)
)

func init() {
	xmp.Register(NsPDF, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &PDFInfo{}
}

func MakeModel(d *xmp.Document) (*PDFInfo, error) {
	m, err := d.MakeModel(NsPDF)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*PDFInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *PDFInfo {
	if m := d.FindModel(NsPDF); m != nil {
		return m.(*PDFInfo)
	}
	return nil
}

// Part 2: 3.1 Adobe PDF namespace
type PDFInfo struct {
	Title        xmp.AltString  `xmp:"dc:title"`
	Author       xmp.StringList `xmp:"dc:creator"`
	Subject      xmp.AltString  `xmp:"dc:description"`
	Keywords     string         `xmp:"pdf:Keywords"`
	Creator      xmp.AgentName  `xmp:"xmp:CreatorTool"`
	Producer     xmp.AgentName  `xmp:"pdf:Producer"`
	PDFVersion   string         `xmp:"pdf:PDFVersion"`
	CreationDate xmp.Date       `xmp:"xmp:CreateDate"`
	ModifyDate   xmp.Date       `xmp:"xmp:ModifyDate"`
	Trapped      xmp.Bool       `xmp:"pdf:Trapped"`
	Copyright    string         `xmp:"pdf:Copyright"`
	Marked       xmp.Bool       `xmp:"pdf:Marked"`
}

func (x PDFInfo) Can(nsName string) bool {
	return NsPDF.GetName() == nsName
}

func (x PDFInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsPDF}
}

func (x *PDFInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *PDFInfo) SyncFromXMP(d *xmp.Document) error {
	if m := dc.FindModel(d); m != nil {
		x.Title = m.Title
		x.Author = m.Creator
		x.Subject = m.Description
		if s := m.Rights.Default(); s != "" {
			x.Copyright = s
		}
	}
	if base := xmpbase.FindModel(d); base != nil {
		x.Creator = base.CreatorTool
		x.CreationDate = base.CreateDate
		x.ModifyDate = base.ModifyDate
	}
	return nil
}

func (x PDFInfo) SyncToXMP(d *xmp.Document) error {
	m, err := dc.MakeModel(d)
	if err != nil {
		return err
	}
	if len(m.Title) == 0 {
		m.Title = x.Title
	}
	if len(m.Creator) == 0 {
		m.Creator = x.Author
	}
	if len(m.Description) == 0 {
		m.Description = x.Subject
	}
	if len(m.Rights) == 0 {
		m.Rights.AddDefault("", x.Copyright)
	}

	// XMP base
	base, err := xmpbase.MakeModel(d)
	if err != nil {
		return err
	}
	if base.CreateDate.IsZero() {
		base.CreateDate = x.CreationDate
	}
	if base.ModifyDate.IsZero() {
		base.ModifyDate = x.ModifyDate
	}
	if base.CreatorTool == "" {
		base.CreatorTool = x.Creator
	}
	return nil
}

func (x *PDFInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *PDFInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsPDF.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *PDFInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsPDF.GetName(), err)
	}
	return nil
}
