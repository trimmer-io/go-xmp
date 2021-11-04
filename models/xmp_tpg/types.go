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
//

// Package xmptpg implements the XMP Paged Text namespace as defined by XMP Specification Part 2.
package xmptpg

import (
	"github.com/trimmer-io/go-xmp/xmp"
)

// Part 2: 1.2.2.2 Dimensions
//
// Default dimensions
// Letter       612x792
// LetterSmall  612x792
// Tabloid      792x1224
// Ledger       1224x792
// Legal        612x1008
// Statement    396x612
// Executive    540x720
// A0           2384x3371
// A1           1685x2384
// A2           1190x1684
// A3           842x1190
// A4           595x842
// A4Small      595x842
// A5           420x595
// B4           729x1032
// B5           516x729
// Envelope     ???x???
// Folio        612x936
// Quarto       610x780
// 10x14        720x1008
type Dimensions struct {
	H    float32 `xmp:"stDim:h,attr"`
	W    float32 `xmp:"stDim:w,attr"`
	Unit Unit    `xmp:"stDim:unit,attr"`
}

func (x Dimensions) IsZero() bool {
	return x.H == 0 && x.W == 0 && x.Unit == ""
}

func (x Dimensions) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t Dimensions
	return e.EncodeElement(_t(x), node)
}

type Unit string

const (
	UnitInch       Unit = "inch"
	UnitMillimeter Unit = "mm"
	UnitPixel      Unit = "pixel"
	UnitPica       Unit = "pica"
	UnitPoint      Unit = "point"
)

// Part 2: 1.2.2.3 Font
type Font struct {
	ChildFontFiles xmp.StringArray `xmp:"stFnt:childFontFiles"`
	Composite      xmp.Bool        `xmp:"stFnt:composite,attr"`
	FontFace       string          `xmp:"stFnt:fontFace,attr"`
	FontFamily     string          `xmp:"stFnt:fontFamily,attr"`
	FontFileName   string          `xmp:"stFnt:fontFileName,attr"`
	FontName       string          `xmp:"stFnt:fontName,attr"`
	FontType       FontType        `xmp:"stFnt:fontType,attr"`
	VersionString  string          `xmp:"stFnt:versionString,attr"`
}

type FontType string

const (
	FontTypeTrueType FontType = "TrueType"
	FontTypeType1    FontType = "Type 1"
	FontTypeOpenType FontType = "Open Type"
)

type FontArray []Font

func (x FontArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x FontArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *FontArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// Part 2: 1.2.2.1 Colorant
type Colorant struct {
	A          int64        `xmp:"xmpG:A,attr"`
	B          int64        `xmp:"xmpG:B,attr"`
	L          float64      `xmp:"xmpG:L,attr"`
	Black      float64      `xmp:"xmpG:black,attr"`
	Cyan       float64      `xmp:"xmpG:cyan,attr"`
	Magenta    float64      `xmp:"xmpG:magenta,attr"`
	Yellow     float64      `xmp:"xmpG:yellow,attr"`
	Blue       int64        `xmp:"xmpG:blue,attr"`
	Green      int64        `xmp:"xmpG:green,attr"`
	Red        int64        `xmp:"xmpG:red,attr"`
	Mode       ColorantMode `xmp:"xmpG:mode,attr"`
	SwatchName string       `xmp:"xmpG:swatchName,attr"`
	Type       ColorType    `xmp:"xmpG:type,attr"`
}

func (x Colorant) IsZero() bool {
	return x.A == 0 && x.B == 0 && x.L == 0 &&
		x.Black == 0 && x.Cyan == 0 && x.Magenta == 0 && x.Yellow == 0 &&
		x.Blue == 0 && x.Green == 0 && x.Red == 0 && x.Mode == "" && x.Type == "" &&
		x.SwatchName == ""
}

func (x Colorant) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t Colorant
	return e.EncodeElement(_t(x), node)
}

type ColorantMode string

const (
	ColorantModeCMYK ColorantMode = "CMYK"
	ColorantModeRGB  ColorantMode = "RGB"
	ColorantModeLAB  ColorantMode = "LAB"
)

type ColorType string

const (
	ColorTypeProcess ColorType = "PROCESS"
	ColorTypeSpot    ColorType = "SPOT"
)

type ColorantList []Colorant

func (x ColorantList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x ColorantList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *ColorantList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type SwatchGroup struct {
	GroupName string       `xmp:"xmpG:groupName"`
	GroupType int          `xmp:"xmpG:groupType"`
	Colorants ColorantList `xmp:"xmpG:Colorants"`
}

type SwatchGroupList []SwatchGroup

func (x SwatchGroupList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x SwatchGroupList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *SwatchGroupList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}
