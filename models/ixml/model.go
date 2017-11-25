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

// TODO:
// - io.Writer interface with maxLen and padding to combine with file writers
// - UID generator
//   <Manufacturer><Device><DeviceSN><Date><Time><Random><FileSetIndex>
//   - 3 chars for the Manufacturer ID: AAT, AVI, DIG, FOS, GAL, HHB, MTI, NAG, ZAX, etc...
//   - 3 chars for the Device ID: CAN, PMX, POR, etc...
//   - 5 chars for the Device Serial Number: 12345
//   - 8 chars for the date: YYYYMMDD
//   - 6 chars for the time: HHMMSS
//   - 3 chars for pure random digits that can very well be the millisecond in the given time second, to ensure that two generations occurring in the same second have a different UID.
//   - 4 chars for the file set index. These 4 digits are set to zero in the <FAMILY_UID> and to the <FILE_SET_INDEX> value in both the <FILE_UID> and BEXT Originator Reference of every files in the set.

// Package ixml implements the iXML audio chunk metadata standard for broadcast wave audio.
package ixml

import (
	"encoding/xml"
	"fmt"
	"strings"

	"trimmer.io/go-xmp/models/xmp_dm"
	"trimmer.io/go-xmp/xmp"
)

const ixmlVersion = "2.0"

var (
	NsIXML = xmp.NewNamespace("iXML", "http://ns.adobe.com/ixml/1.0/", NewModel)
)

func init() {
	xmp.Register(NsIXML, xmp.SoundMetadata)
}

func NewModel(name string) xmp.Model {
	return &IXML{
		Version: ixmlVersion,
	}
}

func MakeModel(d *xmp.Document) (*IXML, error) {
	m, err := d.MakeModel(NsIXML)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*IXML)
	return x, nil
}

func FindModel(d *xmp.Document) *IXML {
	if m := d.FindModel(NsIXML); m != nil {
		return m.(*IXML)
	}
	return nil
}

type IXML struct {
	XMLName              xml.Name                `xml:"BWFXML"                            xmp:"-"`
	Version              string                  `xml:"IXML_VERSION"                      xmp:"iXML:version"`
	Project              string                  `xml:"PROJECT,omitempty"                 xmp:"iXML:project"`
	SceneName            string                  `xml:"SCENE,omitempty"                   xmp:"iXML:sceneName"`
	SoundRoll            string                  `xml:"TAPE,omitempty"                    xmp:"iXML:soundRoll"`
	Take                 int                     `xml:"TAKE,omitempty"                    xmp:"iXML:take"`
	TakeType             TakeTypeList            `xml:"TAKE_TYPE,omitempty"               xmp:"iXML:takeType"` // v2.0+
	IsNotGood            Bool                    `xml:"NO_GOOD,omitempty"                 xmp:"-"`             // v1.9-
	IsFalseStart         Bool                    `xml:"FALSE_START,omitempty"             xmp:"-"`             // v1.9-
	IsWildTrack          Bool                    `xml:"WILD_TRACK,omitempty"              xmp:"-"`             // v1.9-
	PreRecordSamplecount int                     `xml:"PRE_RECORD_SAMPLECOUNT,omitempty"  xmp:"-"`             // v1.3-
	IsCircled            Bool                    `xml:"CIRCLED,omitempty"                 xmp:"iXML:circled"`
	FileUID              string                  `xml:"FILE_UID,omitempty"                xmp:"iXML:fileUid"`
	UserBits             string                  `xml:"UBITS,omitempty"                   xmp:"iXML:userBits"`
	Note                 string                  `xml:"NOTE,omitempty"                    xmp:"iXML:note"`
	SyncPoints           SyncPointList           `xml:"SYNC_POINT_LIST,omitempty"         xmp:"iXML:syncPoints"`
	Speed                *Speed                  `xml:"SPEED,omitempty"                   xmp:"iXML:speed"`
	History              *History                `xml:"HISTORY,omitempty"                 xmp:"iXML:history"`
	FileSet              *FileSet                `xml:"FILE_SET,omitempty"                xmp:"iXML:fileSet"`
	TrackList            TrackList               `xml:"TRACK_LIST,omitempty"              xmp:"iXML:trackList"`
	UserData             *UserData               `xml:"USER,omitempty"                    xmp:"iXML:userData"`
	Location             *Location               `xml:"LOCATION,omitempty"                xmp:"iXML:location"`
	Extension            xmp.NamedExtensionArray `xml:",any"                              xmp:"iXML:extension,any"`
}

func (m *IXML) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsIXML}
}

func (m *IXML) Can(nsName string) bool {
	return nsName == NsIXML.GetName()
}

func (x *IXML) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *IXML) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x *IXML) SyncToXMP(d *xmp.Document) error {
	dm, err := xmpdm.MakeModel(d)
	if err != nil {
		return err
	}
	dm.TapeName = x.SoundRoll
	dm.Scene = x.SceneName
	dm.TrackNumber = x.Take
	dm.LogComment = x.Note
	dm.ProjectName = x.Project
	dm.Good = xmp.Bool(x.IsCircled.Value())
	if x.Speed != nil {
		dm.StartTimecode = x.Speed.XmpTimecode()
		dm.AudioSampleRate = x.Speed.FileSampleRate.Value()
	}
	return nil
}

func (x *IXML) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *IXML) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsIXML.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *IXML) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsIXML.GetName(), err)
	}
	return nil
}

func (x *IXML) ParseXML(data []byte) error {
	// must strip type funcs to avoid recursion
	type _t IXML
	m := &IXML{}
	if err := xml.Unmarshal(data, (*_t)(m)); err != nil {
		return fmt.Errorf("ixml: parse error: %v", err)
	}
	if err := m.upgrade(); err != nil {
		return err
	}
	*x = *m
	return nil
}

func (x *IXML) UnmarshalText(data []byte) error {
	return x.ParseXML(data)
}

func (x *IXML) upgrade() error {
	if err := x.v13_to_v14(); err != nil {
		return err
	}
	if err := x.v1x_to_v20(); err != nil {
		return err
	}
	return nil
}

func (x *IXML) v1x_to_v20() error {
	if x.IsNotGood {
		x.TakeType.Add(TakeTypeNoGood)
		x.IsNotGood = false
	}
	if x.IsFalseStart {
		x.TakeType.Add(TakeTypeFalseStart)
		x.IsFalseStart = false
	}
	if x.IsWildTrack {
		x.TakeType.Add(TakeTypeWildTrack)
		x.IsWildTrack = false
	}
	if x.UserData != nil {
		x.UserData.Comment = strings.TrimSpace(x.UserData.Comment)
	}
	return nil
}

func (x *IXML) v13_to_v14() error {
	if x.PreRecordSamplecount > 0 {
		if i, ok := x.SyncPoints.ContainsFunc(SyncPointPreRecordSamplecount); ok {
			x.SyncPoints[i].Low = x.PreRecordSamplecount
		} else {
			x.SyncPoints = append(x.SyncPoints, SyncPoint{
				Type:     SyncPointRelative,
				Function: SyncPointPreRecordSamplecount,
				Comment:  "upgrade from deprecated property root.PRE_RECORD_SAMPLECOUNT",
				Low:      x.PreRecordSamplecount,
			})
		}
		x.PreRecordSamplecount = 0
	}
	return nil
}
