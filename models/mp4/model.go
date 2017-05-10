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

// Package mp4 implements metadata found in ISO/IEC 14496-14:2003 (3GPP/ISO MP4) files.
package mp4

import (
	"fmt"
	"github.com/echa/go-xmp/models/qt"
	"github.com/echa/go-xmp/xmp"
)

var (
	NsMP4 *xmp.Namespace = &xmp.Namespace{"mp4", "http://ns.apple.com/quicktime/mp4/1.0/", NewModel}
)

func init() {
	xmp.Register(NsMP4, xmp.MovieMetadata)
}

func NewModel(name string) xmp.Model {
	return &MP4Info{}
}

// Note: all tags are single-language versions
type MP4Info struct {
	AlbumAndTrack        string         `mp4:"albm" xmp:"mp4:AlbumAndTrack"`
	Author               string         `mp4:"auth" xmp:"mp4:Author"`
	Classification       string         `mp4:"clsf" xmp:"mp4:Classification"`
	CollectionName       string         `mp4:"coll" xmp:"mp4:CollectionName"`
	Copyright            string         `mp4:"cprt" xmp:"mp4:Copyright"`
	CreationDate         xmp.Date       `mp4:"date" xmp:"mp4:CreationDate"`
	Description          string         `mp4:"dscp" xmp:"mp4:Description"`
	GenreName            string         `mp4:"gnre" xmp:"mp4:GenreName"`
	KeywordList          xmp.StringList `mp4:"kywd" xmp:"mp4:KeywordList"`
	Location             *qt.Location   `mp4:"loci" xmp:"mp4:Location"`
	MediaRatingText      string         `mp4:"rtng" xmp:"mp4:MediaRatingText"`
	Performer            string         `mp4:"perf" xmp:"mp4:Performer"`
	RecordingYear        int            `mp4:"yrrc" xmp:"mp4:RecordingYear"`
	TaggedCharacteristic string         `mp4:"tagc" xmp:"mp4:TaggedCharacteristic"`
	Thumbnail            []byte         `mp4:"thmb" xmp:"mp4:Thumbnail"`
	Title                string         `mp4:"titl" xmp:"mp4:Title"`
	UserRating           int            `mp4:"urat" xmp:"mp4:UserRating"`
}

func (m *MP4Info) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsMP4}
}

func (m *MP4Info) Can(nsName string) bool {
	return nsName == NsMP4.GetName()
}

func (x *MP4Info) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x MP4Info) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *MP4Info) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *MP4Info) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsMP4.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *MP4Info) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsMP4.GetName(), err)
	}
	return nil
}

func (x *MP4Info) ListTags() (xmp.TagList, error) {
	if l, err := xmp.ListNativeFields(x); err != nil {
		return nil, fmt.Errorf("%s: %v", NsMP4.GetName(), err)
	} else {
		return l, nil
	}
}
