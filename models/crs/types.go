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

package crs

import (
	"bytes"
	"fmt"
	"strconv"
	"trimmer.io/go-xmp/xmp"
)

// Point "x, y"
//
type Point struct {
	X int64
	Y int64
}

func (x Point) IsZero() bool {
	return x.X == 0 && x.Y == 0
}

func (x Point) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(strconv.FormatInt(x.X, 10))
	buf.WriteString(", ")
	buf.WriteString(strconv.FormatInt(x.Y, 10))
	return buf.String()
}

func (x Point) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *Point) UnmarshalText(data []byte) error {
	v := string(data)
	r := Point{}
	if _, err := fmt.Sscanf(v, "%d, %d", &r.X, &r.Y); err != nil {
		return fmt.Errorf("xmp: invalid point '%s': %v", v, err)
	}
	*x = r
	return nil
}

type PointList []Point

func (x PointList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x PointList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *PointList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type Correction struct {
	What                  string             `xmp:"crs:What,attr"`
	CorrectionAmount      float32            `xmp:"crs:CorrectionAmount,attr"`
	CorrectionActive      xmp.Bool           `xmp:"crs:CorrectionActive,attr"`
	LocalExposure         float32            `xmp:"crs:LocalExposure,attr"`
	LocalSaturation       float32            `xmp:"crs:LocalSaturation,attr"`
	LocalContrast         float32            `xmp:"crs:LocalContrast,attr"`
	LocalBrightness       float32            `xmp:"crs:LocalBrightness,attr"`
	LocalClarity          float32            `xmp:"crs:LocalClarity,attr"`
	LocalSharpness        float32            `xmp:"crs:LocalSharpness,attr"`
	LocalToningHue        float32            `xmp:"crs:LocalToningHue,attr"`
	LocalToningSaturation float32            `xmp:"crs:LocalToningSaturation,attr"`
	LocalExposure2012     float32            `xmp:"crs:LocalExposure2012,attr"`
	LocalContrast2012     float32            `xmp:"crs:LocalContrast2012,attr"`
	LocalHighlights2012   float32            `xmp:"crs:LocalHighlights2012,attr"`
	LocalShadows2012      float32            `xmp:"crs:LocalShadows2012,attr"`
	LocalClarity2012      float32            `xmp:"crs:LocalClarity2012,attr"`
	LocalLuminanceNoise   float32            `xmp:"crs:LocalLuminanceNoise,attr"`
	LocalMoire            float32            `xmp:"crs:LocalMoire,attr"`
	LocalDefringe         float32            `xmp:"crs:LocalDefringe,attr"`
	LocalTemperature      float32            `xmp:"crs:LocalTemperature,attr"`
	LocalTint             float32            `xmp:"crs:LocalTint,attr"`
	CorrectionMasks       CorrectionMaskList `xmp:"crs:CorrectionMasks"`
}

type CorrectionList []Correction

func (x CorrectionList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x CorrectionList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *CorrectionList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type CorrectionMask struct {
	What           string         `xmp:"crs:What,attr"`
	MaskValue      float32        `xmp:"crs:MaskValue,attr"`
	Radius         float32        `xmp:"crs:Radius,attr"`
	Flow           float32        `xmp:"crs:Flow,attr"`
	CenterWeight   float32        `xmp:"crs:CenterWeight,attr"`
	Dabs           xmp.StringList `xmp:"crs:Dabs"`
	ZeroX          float32        `xmp:"crs:ZeroX,attr"`
	ZeroY          float32        `xmp:"crs:ZeroY,attr"`
	FullX          float32        `xmp:"crs:FullX,attr"`
	FullY          float32        `xmp:"crs:FullY,attr"`
	Top            float32        `xmp:"crs:Top,attr"`
	Left           float32        `xmp:"crs:Left,attr"`
	Bottom         float32        `xmp:"crs:Bottom,attr"`
	Right          float32        `xmp:"crs:Right,attr"`
	Angle          float32        `xmp:"crs:Angle,attr"`
	Midpoint       float32        `xmp:"crs:Midpoint,attr"`
	Roundness      float32        `xmp:"crs:Roundness,attr"`
	Feather        float32        `xmp:"crs:Feather,attr"`
	Flipped        xmp.Bool       `xmp:"crs:Flipped,attr"`
	Version        int64          `xmp:"crs:Version,attr"`
	SizeX          float32        `xmp:"crs:SizeX,attr"`
	SizeY          float32        `xmp:"crs:SizeY,attr"`
	X              float32        `xmp:"crs:X,attr"`
	Y              float32        `xmp:"crs:Y,attr"`
	Alpha          float32        `xmp:"crs:Alpha,attr"`
	CenterValue    float32        `xmp:"crs:CenterValue,attr"`
	PerimeterValue float32        `xmp:"crs:PerimeterValue,attr"`
}

type CorrectionMaskList []CorrectionMask

func (x CorrectionMaskList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x CorrectionMaskList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *CorrectionMaskList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type RetouchArea struct {
	SpotType    string             `xmp:"crs:SpotType,attr"`    // heal
	SourceState string             `xmp:"crs:SourceState,attr"` // sourceSetExplicitly
	Method      string             `xmp:"crs:Method,attr"`      // gaussian
	SourceX     float32            `xmp:"crs:SourceX,attr"`
	OffsetY     float32            `xmp:"crs:OffsetY,attr"`
	Opacity     float32            `xmp:"crs:Opacity,attr"`
	Feather     float32            `xmp:"crs:Feather,attr"`
	Seed        int64              `xmp:"crs:Seed,attr"`
	Masks       CorrectionMaskList `xmp:"crs:Masks"`
}

type RetouchAreaList []RetouchArea

func (x RetouchAreaList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x RetouchAreaList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *RetouchAreaList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// TODO: (Un)MarshalText
// <rdf:li>centerX = 0.290216, centerY = 0.190122, radius = 0.005845, sourceState = sourceSetExplicitly, sourceX = 0.260958, sourceY = 0.189526, spotType = heal</rdf:li>
type RetouchInfo struct {
	SpotType    string
	SourceState string
	CenterX     float32
	CenterY     float32
	Radius      float32
	SourceX     float32
	SourceY     float32
}
