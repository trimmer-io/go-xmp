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

// Package tiff implements metadata for TIFF and JPEG image files as defined
// by XMP Specification Part 1.
package tiff

import (
	"fmt"
	"strings"

	"trimmer.io/go-xmp/models/dc"
	"trimmer.io/go-xmp/models/xmp_base"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsTiff = xmp.NewNamespace("tiff", "http://ns.adobe.com/tiff/1.0/", NewModel)
)

func init() {
	xmp.Register(NsTiff, xmp.ImageMetadata)
}

func NewModel(name string) xmp.Model {
	return &TiffInfo{}
}

func MakeModel(d *xmp.Document) (*TiffInfo, error) {
	m, err := d.MakeModel(NsTiff)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*TiffInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *TiffInfo {
	if m := d.FindModel(NsTiff); m != nil {
		return m.(*TiffInfo)
	}
	return nil
}

type TiffInfo struct {
	Artist                    xmp.StringList    `xmp:"dc:creator"`
	BitsPerSample             xmp.IntList       `xmp:"tiff:BitsPerSample"` // 3 components
	Compression               CompressionType   `xmp:"tiff:Compression"`
	DateTime                  xmp.Date          `xmp:"xmp:ModifyDate"`
	ImageLength               int               `xmp:"tiff:ImageLength"`
	ImageWidth                int               `xmp:"tiff:ImageWidth"` // A. Tags relating to image data structure
	Make                      string            `xmp:"tiff:Make"`
	Model                     string            `xmp:"tiff:Model"`
	Software                  string            `xmp:"xmp:CreatorTool"`
	ImageDescription          xmp.AltString     `xmp:"dc:description"`
	Copyright                 xmp.AltString     `xmp:"dc:rights"`
	Orientation               OrientationType   `xmp:"tiff:Orientation"`
	PhotometricInterpretation ColorModel        `xmp:"tiff:PhotometricInterpretation"`
	PlanarConfiguration       PlanarType        `xmp:"tiff:PlanarConfiguration"`
	PrimaryChromaticities     xmp.RationalArray `xmp:"tiff:PrimaryChromaticities"` // 6 components
	ReferenceBlackWhite       xmp.RationalArray `xmp:"tiff:ReferenceBlackWhite"`   // 6 components
	ResolutionUnit            ResolutionUnit    `xmp:"tiff:ResolutionUnit"`        // 2 = inches, 3 = centimeters
	SamplesPerPixel           int               `xmp:"tiff:SamplesPerPixel"`
	TransferFunction          xmp.IntList       `xmp:"tiff:TransferFunction"` // C. Tags relating to image data characteristics
	WhitePoint                xmp.RationalArray `xmp:"tiff:WhitePoint"`
	XResolution               xmp.Rational      `xmp:"tiff:XResolution"`
	YCbCrCoefficients         xmp.RationalArray `xmp:"tiff:YCbCrCoefficients"` // 3 components
	YCbCrPositioning          YCbCrPosition     `xmp:"tiff:YCbCrPositioning"`
	YCbCrSubSampling          YCbCrSubSampling  `xmp:"tiff:YCbCrSubSampling"`
	YResolution               xmp.Rational      `xmp:"tiff:YResolution"`
	NativeDigest              string            `xmp:"tiff:NativeDigest,omit"` // ignore according to spec

	// be resilient to broken writers: read tags erroneously
	// mapped to the wrong XMP properties, but do not output them
	// when writing ourselves
	X_Artist           string        `xmp:"tiff:Artist,omit"`
	X_DateTime         xmp.Date      `xmp:"tiff:DateTime,omit"`
	X_Software         string        `xmp:"tiff:Software,omit"`
	X_ImageDescription xmp.AltString `xmp:"tiff:ImageDescription,omit"`
	X_Copyright        xmp.AltString `xmp:"tiff:Copyright,omit"`
}

func (x TiffInfo) Can(nsName string) bool {
	return NsTiff.GetName() == nsName
}

func (x TiffInfo) Namespaces() xmp.NamespaceList {
	return []*xmp.Namespace{NsTiff}
}

func (x *TiffInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *TiffInfo) SyncFromXMP(d *xmp.Document) error {
	if m := dc.FindModel(d); m != nil {
		x.Artist = m.Creator
		x.ImageDescription = m.Description
		x.Copyright = m.Rights
	}
	if m := xmpbase.FindModel(d); m != nil {
		x.DateTime = m.ModifyDate
		if m.CreatorTool != "" {
			x.Software = m.CreatorTool.String()
		}
	}
	return nil
}

// also remap X_* attributes to the correct standard positions, but don't overwrite
func (x TiffInfo) SyncToXMP(d *xmp.Document) error {
	m, err := dc.MakeModel(d)
	if err != nil {
		return err
	}
	if len(m.Creator) == 0 && len(x.Artist) > 0 {
		m.Creator = x.Artist
	}
	if len(m.Creator) == 0 && len(x.X_Artist) > 0 {
		m.Creator = strings.Split(x.X_Artist, ",")
	}
	if len(m.Description) == 0 {
		m.Description = x.ImageDescription
	}
	if len(m.Description) == 0 {
		m.Description = x.X_ImageDescription
	}
	if len(m.Rights) == 0 {
		m.Rights = x.Copyright
	}
	if len(m.Rights) == 0 {
		m.Rights = x.X_Copyright
	}

	// XMP base
	base, err := xmpbase.MakeModel(d)
	if err != nil {
		return err
	}
	if base.ModifyDate.IsZero() {
		base.ModifyDate = x.DateTime
	}
	if base.ModifyDate.IsZero() {
		base.ModifyDate = x.X_DateTime
	}
	if base.CreatorTool.IsZero() && x.Software != "" {
		base.CreatorTool = xmp.AgentName(x.Software)
	}
	return nil
}

func (x *TiffInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *TiffInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsTiff.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *TiffInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsTiff.GetName(), err)
	}
	return nil
}
