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
	"fmt"
	"strings"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsXmp = xmp.NewNamespace("xmp", "http://ns.adobe.com/xap/1.0/", NewModel)

	nsXmpIdq  = xmp.NewNamespace("xmpidq", "http://ns.adobe.com/xmp/Identifier/qual/1.0/", nil)
	nsXmpG    = xmp.NewNamespace("xmpG", "http://ns.adobe.com/xap/1.0/g/", nil)
	nsXmpGImg = xmp.NewNamespace("xmpGImg", "http://ns.adobe.com/xap/1.0/g/img/", nil)
)

func init() {
	xmp.Register(NsXmp, xmp.XmpMetadata)
	xmp.Register(nsXmpIdq)
	xmp.Register(nsXmpG)
	xmp.Register(nsXmpGImg)
}

func NewModel(name string) xmp.Model {
	return &XmpBase{}
}

func MakeModel(d *xmp.Document) (*XmpBase, error) {
	m, err := d.MakeModel(NsXmp)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*XmpBase)
	return x, nil
}

func FindModel(d *xmp.Document) *XmpBase {
	if m := d.FindModel(NsXmp); m != nil {
		return m.(*XmpBase)
	}
	return nil
}

type XmpBase struct {
	Advisory     xmp.StringArray         `xmp:"xmp:Advisory"`
	BaseURL      xmp.Uri                 `xmp:"xmp:BaseURL"`
	CreateDate   xmp.Date                `xmp:"xmp:CreateDate"`
	CreatorTool  xmp.AgentName           `xmp:"xmp:CreatorTool"`
	Identifier   IdentifierArray         `xmp:"xmp:Identifier"`
	Label        string                  `xmp:"xmp:Label"`
	MetadataDate xmp.Date                `xmp:"xmp:MetadataDate"`
	ModifyDate   xmp.Date                `xmp:"xmp:ModifyDate"`
	Nickname     string                  `xmp:"xmp:Nickname"`
	Rating       Rating                  `xmp:"xmp:Rating"`
	Thumbnails   ThumbnailArray          `xmp:"xmp:Thumbnails"`
	Extensions   xmp.NamedExtensionArray `xmp:"xmp:extension"`
}

func (x XmpBase) Can(nsName string) bool {
	return NsXmp.GetName() == nsName
}

func (x XmpBase) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsXmp}
}

func prefixer(prefix string) xmp.ConverterFunc {
	return func(val string) string {
		return strings.Join([]string{prefix, val}, ":")
	}
}

var identifierDesc = xmp.SyncDescList{
	&xmp.SyncDesc{"trim:Asset/UUID", "xmp:Identifier", xmp.MERGE, nil},
	&xmp.SyncDesc{"exif:ImageUniqueID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:exif")},
	&xmp.SyncDesc{"arri:UUID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:arri")},
	&xmp.SyncDesc{"arri:SMPTE_UMID", "xmp:Identifier", xmp.MERGE, prefixer("umid:smpte")},
	&xmp.SyncDesc{"iXML:fileUid", "xmp:Identifier", xmp.MERGE, prefixer("uuid:ixml")},
	&xmp.SyncDesc{"qt:mdta/ContentIdentifier", "xmp:Identifier", xmp.MERGE, prefixer("uuid:cid")},
	&xmp.SyncDesc{"iTunes:StoreFrontID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:sfid")},
	&xmp.SyncDesc{"qt:GUID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:guid")},
	&xmp.SyncDesc{"qt:ISRCCode", "xmp:Identifier", xmp.MERGE, prefixer("uuid:isrc")},
	&xmp.SyncDesc{"qt:ContentID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:cid")},
	&xmp.SyncDesc{"qt:ClipID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:clipid")},
	&xmp.SyncDesc{"qt:proapps/ClipID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:clipid")},
	&xmp.SyncDesc{"bext:umid", "xmp:Identifier", xmp.MERGE, prefixer("umid:smpte")},
	&xmp.SyncDesc{"id3:podcastID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:podcast")},
	&xmp.SyncDesc{"id3:uniqueFileIdentifier", "xmp:Identifier", xmp.MERGE, prefixer("uuid:id3")},
	&xmp.SyncDesc{"GettyImagesGIFT:AssetID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:getty")},
	&xmp.SyncDesc{"Iptc4xmpExt:DigImageGUID", "xmp:Identifier", xmp.MERGE, prefixer("uuid:iptc")},
	&xmp.SyncDesc{"mxf:PackageID", "xmp:Identifier", xmp.MERGE, prefixer("umid:smpte")},

	// TODO: map more id schemes here
	// FCPX: Asset UID from fcpxml          uuid:fcpx:
	// Panasonic: GlobalClipID (SMPTE UMID) umid:smpte:
}

func (x *XmpBase) SyncFromXMP(d *xmp.Document) error {
	// sync list
	if err := d.SyncMulti(identifierDesc, x); err != nil {
		return err
	}
	return nil
}

func (x XmpBase) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *XmpBase) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *XmpBase) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsXmp.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *XmpBase) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsXmp.GetName(), err)
	}
	return nil
}

func (x *XmpBase) AddID(id string) {
	x.Identifier = append(x.Identifier, Identifier{ID: id})
}

func (x *XmpBase) AddIDWithScheme(id, scheme string) {
	x.Identifier = append(x.Identifier, Identifier{ID: id, Scheme: scheme})
}
