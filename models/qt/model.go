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

// Quicktime Metadata storage options
//
// 1) Quicktime Metadata
//   - namespace, key, value triple stored in mdta atoms
//   - namespace: reverse-DNS name
//   - some quicktime tags may be exported this way as well
//   - camera manufacturers (e.g. Arri) use this to store clip metadata
// 2) Quicktime User Data
//   - Four-CC tags defined by Apple or 3GPP ISO standard
//   - stored in udta atoms
//   - flavours
//      - 3GPP: subset for ISO/MP4 files, different tag values, no multi-lang strings
//      - iTunes: specific to the iTunes store (used for music, video, apps, books)
//      - Quicktime: full set with multi-language string support
//
// References:
//
//   1) https://developer.apple.com/library/content/documentation/QuickTime/QTFF/QTFFChap2/qtff2.html#//apple_ref/doc/uid/TP40000939-CH204-25538
//   2) http://search.cpan.org/dist/MP4-Info-1.04/
//   3) http://xhelmboyx.tripod.com/formats/mp4-layout.txt
//   4) http://wiki.multimedia.cx/index.php?title=Apple_QuickTime
//   5) ISO 14496-12 (http://read.pudn.com/downloads64/ebook/226547/ISO_base_media_file_format.pdf)
//   6) ISO 14496-16 (http://www.iec-normen.de/previewpdf/info_isoiec14496-16%7Bed2.0%7Den.pdf)
//   7) http://atomicparsley.sourceforge.net/mpeg-4files.html
//   8) http://www.adobe.com/devnet/xmp/pdfs/XMPSpecificationPart3.pdf (Oct 2008)
//   9) QuickTime file format specification 2010-05-03
//   10) http://standards.iso.org/ittf/PubliclyAvailableStandards/c051533_ISO_IEC_14496-12_2008.zip
//   11) http://getid3.sourceforge.net/source/module.audio-video.quicktime.phps
//   12) http://qtra.apple.com/atoms.html
//   13) http://www.etsi.org/deliver/etsi_ts/126200_126299/126244/10.01.00_60/ts_126244v100100p.pdf
//   14) https://github.com/appsec-labs/iNalyzer/blob/master/scinfo.m
//   15) http://nah6.com/~itsme/cvs-xdadevtools/iphone/tools/decodesinf.pl
//   16) https://github.com/sannies/mp4parser
//   17) http://www.mp4ra.org/filetype.html
//   18) https://github.com/WordPress/WordPress/blob/master/wp-includes/ID3/module.audio-video.quicktime.php

// Package qt implements metadata found in Apple Quicktime (MOV) files.
package qt

import (
	"fmt"
	"strings"

	"trimmer.io/go-xmp/models/ixml"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsQuicktime = xmp.NewNamespace("qt", "http://ns.apple.com/quicktime/1.0/", NewModel)
)

func init() {
	xmp.Register(NsQuicktime, xmp.MovieMetadata, xmp.CameraMetadata)
}

func NewModel(name string) xmp.Model {
	return &QtInfo{}
}

type QtInfo struct {
	Udta    *QtUserdata `qt:"-" xmp:"qt:udta"`
	Mdta    *QtMetadata `qt:"-" xmp:"qt:mdta"`
	Player  *QtPlayer   `qt:"-" xmp:"qt:player"`
	ProApps *QtProApps  `qt:"-" xmp:"qt:proapps"`

	// external structs
	IXML *ixml.IXML    `qt:"-" xmp:"-"`
	XMP  *xmp.Document `qt:"XMP_" xmp:"-"`

	// unknown 3rd party tags
	Extension xmp.TagList `qt:",any" xmp:"qt:extension"`
}

func (m *QtInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsQuicktime}
}

func (m *QtInfo) Can(nsName string) bool {
	return nsName == NsQuicktime.GetName()
}

func (x *QtInfo) SyncFromXMP(d *xmp.Document) error {
	if x.Udta != nil {
		if err := x.Udta.SyncFromXMP(d); err != nil {
			return err
		}
	}
	if x.Mdta != nil {
		if err := x.Mdta.SyncFromXMP(d); err != nil {
			return err
		}
	}
	if x.Player != nil {
		if err := x.Player.SyncFromXMP(d); err != nil {
			return err
		}
	}
	if x.ProApps != nil {
		if err := x.ProApps.SyncFromXMP(d); err != nil {
			return err
		}
	}
	return nil
}

func (x QtInfo) SyncToXMP(d *xmp.Document) error {
	if x.XMP != nil {
		if err := d.Merge(x.XMP, xmp.MERGE); err != nil {
			return err
		}
	}
	if x.Udta != nil {
		if err := x.Udta.SyncToXMP(d); err != nil {
			return err
		}
	}
	if x.Mdta != nil {
		if err := x.Mdta.SyncToXMP(d); err != nil {
			return err
		}
	}
	if x.Player != nil {
		if err := x.Player.SyncToXMP(d); err != nil {
			return err
		}
	}
	if x.ProApps != nil {
		if err := x.ProApps.SyncToXMP(d); err != nil {
			return err
		}
	}
	if x.IXML != nil {
		d.AddModel(x.IXML)
	}
	return nil
}

func (x *QtInfo) CanTag(tag string) bool {
	switch {
	case len(tag) == 4:
		v := &QtUserdata{}
		_, err := xmp.GetNativeField(v, tag)
		return err == nil
	case strings.HasPrefix("com.apple.quicktime.player", tag):
		v := &QtPlayer{}
		return v.CanTag(tag)
	case strings.HasPrefix("com.apple.quicktime", tag):
		v := &QtMetadata{}
		return v.CanTag(tag)
	case strings.HasPrefix("com.apple.proapps", tag):
		v := &QtProApps{}
		return v.CanTag(tag)
	case tag == "XMP_" || tag == "iXML" || tag == "info.ixml.xml" || tag == "info.ixml.metadata" || tag == "info.ixml.info":
		return true
	}
	return false
}

func (x *QtInfo) GetTag(tag string) (string, error) {
	switch {
	case len(tag) == 4:
		if x.Udta == nil {
			return "", nil
		}
		if v, err := xmp.GetNativeField(x.Udta, tag); err != nil {
			return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		} else {
			return v, nil
		}
	case strings.HasPrefix(tag, "com.apple.quicktime.player"):
		if x.Player == nil {
			return "", nil
		}
		if v, err := xmp.GetNativeField(x.Player, tag); err != nil {
			return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		} else {
			return v, nil
		}
	case strings.HasPrefix(tag, "com.apple.quicktime"):
		if x.Mdta == nil {
			return "", nil
		}
		if v, err := xmp.GetNativeField(x.Mdta, tag); err != nil {
			return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		} else {
			return v, nil
		}
	case strings.HasPrefix(tag, "com.apple.proapps"):
		if x.ProApps == nil {
			return "", nil
		}
		if v, err := xmp.GetNativeField(x.ProApps, tag); err != nil {
			return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		} else {
			return v, nil
		}
	}
	return "", nil
}

func (x *QtInfo) SetTag(tag, value string) error {
	switch {
	case len(tag) == 4:
		if x.Udta == nil {
			x.Udta = &QtUserdata{}
		}
		if err := xmp.SetNativeField(x.Udta, tag, value); err != nil {
			return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		}
	case strings.HasPrefix(tag, "com.apple.quicktime.player"):
		if x.Player == nil {
			x.Player = &QtPlayer{}
		}
		if err := xmp.SetNativeField(x.Player, tag, value); err != nil {
			return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		}
	case strings.HasPrefix(tag, "com.apple.quicktime"):
		if x.Mdta == nil {
			x.Mdta = &QtMetadata{}
		}
		if err := xmp.SetNativeField(x.Mdta, tag, value); err != nil {
			return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		}
	case strings.HasPrefix(tag, "com.apple.proapps"):
		if x.ProApps == nil {
			x.ProApps = &QtProApps{}
		}
		if err := xmp.SetNativeField(x.ProApps, tag, value); err != nil {
			return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
		}
	case tag == "iXML" || tag == "info.ixml.xml" || tag == "info.ixml.metadata" || tag == "info.ixml.info":
		i := &ixml.IXML{}
		if err := i.ParseXML([]byte(value)); err != nil {
			return fmt.Errorf("%s: parsing ixml: %v", NsQuicktime.GetName(), err)
		} else {
			x.IXML = i
		}
	case tag == "XMP_":
		v := xmp.NewDocument()
		if err := xmp.Unmarshal([]byte(value), v); err != nil {
			return fmt.Errorf("%s: parsing embedded xmp: %v", NsQuicktime.GetName(), err)
		} else {
			x.XMP = v
		}
	default:
		// silently ignore all other sorts of tags
	}
	return nil
}

func (x *QtInfo) ListTags() (xmp.TagList, error) {
	if l, err := xmp.ListNativeFields(x); err != nil {
		return nil, fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	} else {
		return l, nil
	}
}
