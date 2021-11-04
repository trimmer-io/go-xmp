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

// see also Photoshop Specification at
// http://www.adobe.com/devnet-apps/photoshop/fileformatashtml

// Package ps implements Adobe Photoshop metadata as defined by XMP Specification Part 2 Chapter 3.2.
package ps

import (
	"fmt"
	"github.com/trimmer-io/go-xmp/xmp"
)

var (
	NsPhotoshop = xmp.NewNamespace("photoshop", "http://ns.adobe.com/photoshop/1.0/", NewModel)
)

func init() {
	xmp.Register(NsPhotoshop, xmp.ImageMetadata)
}

func NewModel(name string) xmp.Model {
	return &PhotoshopInfo{}
}

func MakeModel(d *xmp.Document) (*PhotoshopInfo, error) {
	m, err := d.MakeModel(NsPhotoshop)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*PhotoshopInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *PhotoshopInfo {
	if m := d.FindModel(NsPhotoshop); m != nil {
		return m.(*PhotoshopInfo)
	}
	return nil
}

type PhotoshopInfo struct {
	AuthorsPosition        string          `xmp:"photoshop:AuthorsPosition"`
	CaptionWriter          string          `xmp:"photoshop:CaptionWriter"`
	Category               string          `xmp:"photoshop:Category"`
	City                   string          `xmp:"photoshop:City"`
	ColorMode              ColorMode       `xmp:"photoshop:ColorMode"`
	Country                string          `xmp:"photoshop:Country"`
	Credit                 string          `xmp:"photoshop:Credit"`
	DateCreated            xmp.Date        `xmp:"photoshop:DateCreated"`
	DocumentAncestors      AnchestorArray  `xmp:"photoshop:DocumentAncestors"`
	Headline               string          `xmp:"photoshop:Headline"`
	History                string          `xmp:"photoshop:History"`
	ICCProfile             string          `xmp:"photoshop:ICCProfile"`
	Instructions           string          `xmp:"photoshop:Instructions"`
	Layer                  *Layer          `xmp:"photoshop:Layer"`
	SidecarForExtension    string          `xmp:"photoshop:SidecarForExtension"` // "NEF"
	Source                 string          `xmp:"photoshop:Source"`
	State                  string          `xmp:"photoshop:State"`
	SupplementalCategories xmp.StringArray `xmp:"photoshop:SupplementalCategories"`
	TextLayers             LayerList       `xmp:"photoshop:TextLayers"`
	TransmissionReference  string          `xmp:"photoshop:TransmissionReference"`
	Urgency                int             `xmp:"photoshop:Urgency"` // 1 - 8

	EmbeddedXMPDigest string `xmp:"photoshop:EmbeddedXMPDigest,omit"` // "00000000000000000000000000000000"
	LegacyIPTCDigest  string `xmp:"photoshop:LegacyIPTCDigest,omit"`  // "AA5133A9479EA0F732E6A7414060A81F"
}

func (x PhotoshopInfo) Can(nsName string) bool {
	return NsPhotoshop.GetName() == nsName
}

func (x PhotoshopInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsPhotoshop}
}

func (x *PhotoshopInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *PhotoshopInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x PhotoshopInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *PhotoshopInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *PhotoshopInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsPhotoshop.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *PhotoshopInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsPhotoshop.GetName(), err)
	}
	return nil
}

// 3.2.1.1 Ancestor
type Anchestor struct {
	AncestorID xmp.Uri `xmp:"photoshop:AncestorID"`
}

type AnchestorArray []Anchestor

func (x AnchestorArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x AnchestorArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *AnchestorArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// 3.2.1.2 Layer
type Layer struct {
	LayerName string `xmp:"photoshop:LayerName"`
	LayerText string `xmp:"photoshop:LayerText"`
}

func (x Layer) IsZero() bool {
	return x.LayerName == "" && x.LayerText == ""
}

type LayerList []Layer

func (x LayerList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x LayerList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *LayerList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type ColorMode int

const (
	ColorModeBitmap       = 0 // 0 = Bitmap
	ColorModeGrayscale    = 1 // 1 = Gray scale
	ColorModeIndexed      = 2 // 2 = Indexed colour
	ColorModeRGB          = 3 // 3 = RGB colour
	ColorModeCMYK         = 4 // 4 = CMYK colour
	ColorModeMultiChannel = 7 // 7 = Multi-channel
	ColorModeDuotone      = 8 // 8 = Duotone
	ColorModeLAB          = 9 // 9 = LAB colour
)
