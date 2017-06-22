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

// Package xmpbase implements the XMP namespace as defined by XMP Specification Part 2.
package xmpbase

import (
	"bytes"
	"github.com/echa/go-xmp/xmp"
)

type Rating int

const (
	RatingRejected Rating = -1
	RatingUnrated  Rating = 0
	Rating1        Rating = 1
	Rating2        Rating = 2
	Rating3        Rating = 3
	Rating4        Rating = 4
	Rating5        Rating = 5
)

// 8.7 xmpidq namespace
//
// A qualifier providing the name of the formal identification scheme
// used for an item in the xmp:Identifier array.
//
// Form 1: <rdf:li>http://www.example.com/</rdf:li>
// Form 2:
//   <rdf:li>
//     <rdf:Description>
//       <rdf:value>http://www.example.com/xmp_example/xmp/identifier3</rdf:value>
//       <xmpidq:Scheme>myscheme</xmpidq:Scheme>
//     </rdf:Description>
//   </rdf:li>
// Form 3:
//   <rdf:li xmpidq:Scheme="myscheme">http://www.example.com/</rdf:li>

type Identifier struct {
	ID     string `xmp:"rdf:value"`
	Scheme string `xmp:"xmpidq:Scheme"`
}

func (x Identifier) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.Scheme == "" {
		return e.EncodeElement(x.ID, node)
	} else {
		id := struct {
			ID     string `xmp:"rdf:value"`
			Scheme string `xmp:"xmpidq:Scheme"`
		}{
			ID:     x.ID,
			Scheme: x.Scheme,
		}
		return e.EncodeElement(id, node)
	}
}

func (x *Identifier) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	id := struct {
		ID     string `xmp:"rdf:value"`
		Scheme string `xmp:"xmpidq:Scheme"`
	}{
		ID:     x.ID,
		Scheme: x.Scheme,
	}
	var err error
	if len(node.Nodes) == 0 && len(node.Attr) == 0 {
		err = d.DecodeElement(&id.ID, node) // value-only from node.Value
	} else if len(node.Nodes) == 0 {
		err = d.DecodeElement(&id.ID, node) // value from node.Value
		err = d.DecodeElement(&id, node)    // scheme from attr
	} else {
		err = d.DecodeElement(&id, node.Nodes[0]) // both from child nodes
	}
	if err != nil {
		return err
	}
	x.ID = id.ID
	x.Scheme = id.Scheme
	return nil
}

func (x Identifier) MarshalText() ([]byte, error) {
	buf := bytes.Buffer{}
	if x.Scheme != "" {
		buf.WriteString(x.Scheme)
		buf.WriteByte(':')
	}
	buf.WriteString(x.ID)
	return buf.Bytes(), nil
}

func (x *Identifier) UnmarshalText(data []byte) error {
	x.ID = string(data)
	return nil
}

type IdentifierArray []Identifier

func NewIdentifierArray(items ...interface{}) IdentifierArray {
	x := make(IdentifierArray, 0)
	for _, v := range items {
		if s := xmp.ToString(v); s != "" {
			x = append(x, Identifier{ID: s})
		}
	}
	return x
}

func (x IdentifierArray) IsZero() bool {
	return len(x) == 0
}

func (x *IdentifierArray) AddUnique(v string) error {
	if !x.Contains(v) {
		x.Add(v)
	}
	return nil
}

func (x *IdentifierArray) Add(value string) {
	*x = append(*x, Identifier{ID: value})
}

func (x *IdentifierArray) Contains(v string) bool {
	return x.Index(v) > -1
}

func (x *IdentifierArray) Index(val string) int {
	for i, v := range *x {
		if v.ID == val {
			return i
		}
	}
	return -1
}

func (x *IdentifierArray) Remove(v string) {
	if idx := x.Index(v); idx > -1 {
		*x = append((*x)[:idx], (*x)[:idx+1]...)
	}
}

func (x IdentifierArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x IdentifierArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *IdentifierArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// Part 2: 1.2.2.4 Thumbnail
type Thumbnail struct {
	Format string `xmp:"xmpGImg:format"`
	Width  int64  `xmp:"xmpGImg:height"`
	Height int64  `xmp:"xmpGImg:width"`
	Image  []byte `xmp:"xmpGImg:image"`
}

func (x Thumbnail) IsZero() bool {
	return x.Format == "" && x.Width == 0 && x.Height == 0 && len(x.Image) == 0
}

func (x Thumbnail) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t Thumbnail
	return e.EncodeElement(_t(x), node)
}

type ThumbnailArray []Thumbnail

func (x ThumbnailArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeAlternative
}

func (x ThumbnailArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *ThumbnailArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// Metadata Workinggroup Guidelines 2.0, 2010
// stArea
// http://ns.adobe.com/xmp/sType/Area#
type Area struct {
	X    float64  `xmp:"stArea:x,attr"`
	Y    float64  `xmp:"stArea:y,attr"`
	W    float64  `xmp:"stArea:w,attr"`
	H    float64  `xmp:"stArea:h,attr"`
	D    float64  `xmp:"stArea:d,attr"`
	Unit AreaUnit `xmp:"stArea:unit,attr"`
}

type AreaUnit string

const (
	AreaUnitRatio AreaUnit = "normalized"
	AreaUnitPixel AreaUnit = "pixel"
)

func (x Area) IsZero() bool {
	return x.X == 0 && x.Y == 0 && x.W == 0 && x.H == 0 && x.D == 0 && x.Unit == ""
}

func (x Area) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t Area
	return e.EncodeElement(_t(x), node)
}
