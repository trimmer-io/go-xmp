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

package ixml

import (
	"encoding/xml"
	"fmt"
	"strings"

	"trimmer.io/go-xmp/models/xmp_dm"
	"trimmer.io/go-xmp/xmp"
)

type Bool bool

func (x Bool) Value() bool {
	return bool(x)
}

func (x Bool) MarshalText() ([]byte, error) {
	if x {
		return []byte("TRUE"), nil
	}
	return []byte("FALSE"), nil
}

func (x *Bool) UnmarshalText(data []byte) error {
	switch strings.ToLower(string(data)) {
	case "true":
		*x = true
	case "false":
		*x = false
	default:
		return fmt.Errorf("ixml: invalid bool value '%s'", string(data))
	}
	return nil
}

type SyncPoint struct {
	XMLName       xml.Name              `xml:"SYNC_POINT"                            xmp:"-"`
	Type          SyncPointType         `xml:"SYNC_POINT_TYPE"                       xmp:"iXML:syncPointType"`
	Function      SyncPointFunctionType `xml:"SYNC_POINT_FUNCTION"                   xmp:"iXML:syncPointFunction"`
	Comment       string                `xml:"SYNC_POINT_COMMENT,omitempty"          xmp:"iXML:syncPointComment"`
	Low           int                   `xml:"SYNC_POINT_LOW,omitempty"              xmp:"iXML:syncPointLow"`
	High          int                   `xml:"SYNC_POINT_HIGH,omitempty"             xmp:"iXML:syncPointHigh"`
	EventDuration int64                 `xml:"SYNC_POINT_EVENT_DURATION,omitempty"   xmp:"iXML:syncPointEventDuration"`
}

type SyncPointList []SyncPoint

func (x SyncPointList) ContainsFunc(f SyncPointFunctionType) (int, bool) {
	for i, v := range x {
		if v.Function == f {
			return i, true
		}
	}
	return -1, false
}

func (x *SyncPointList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	depth := 1
	for depth > 0 {
		t, err := d.Token()
		if err != nil {
			return fmt.Errorf("%s: %v", NsIXML.GetName(), err)
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local != "SYNC_POINT" {
				d.Skip()
				continue
			}
			sp := SyncPoint{}
			if err := d.DecodeElement(&sp, &t); err != nil {
				return fmt.Errorf("%s: %v", NsIXML.GetName(), err)
			}
			*x = append(*x, sp)

		case xml.EndElement:
			depth--
		}
	}
	return nil
}

func (x SyncPointList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(x) == 0 {
		return nil
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	cToken := xml.StartElement{Name: xml.Name{Local: "SYNC_POINT_COUNT"}}
	if err := e.EncodeElement(len(x), cToken); err != nil {
		return err
	}
	for _, v := range x {
		if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "SYNC_POINT"}}); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return e.Flush()
}

func (x SyncPointList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x SyncPointList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *SyncPointList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type Speed struct {
	XMLName                         xml.Name     `xml:"SPEED"                                          xmp:"-"`
	Note                            string       `xml:"NOTE,omitempty"                                 xmp:"iXML:note"`
	MasterSpeed                     xmp.Rational `xml:"MASTER_SPEED,omitempty"                         xmp:"iXML:masterSpeed"`
	CurrentSpeed                    xmp.Rational `xml:"CURRENT_SPEED,omitempty"                        xmp:"iXML:currentSpeed"`
	TimecodeRate                    xmp.Rational `xml:"TIMECODE_RATE,omitempty"                        xmp:"iXML:timecodeRate"`
	TimecodeFlag                    TimecodeFlag `xml:"TIMECODE_FLAG,omitempty"                        xmp:"iXML:timecodeFlag"`
	FileSampleRate                  xmp.Rational `xml:"FILE_SAMPLE_RATE,omitempty"                     xmp:"iXML:fileSampleRate"`
	AudioBitDepth                   int          `xml:"AUDIO_BIT_DEPTH,omitempty"                      xmp:"iXML:audioBitDepth"`
	DigitizerSampleRate             int          `xml:"DIGITIZER_SAMPLE_RATE,omitempty"                xmp:"iXML:digitizerSampleRate"`
	TimestampSamplesSinceMidnightHI int          `xml:"TIMESTAMP_SAMPLES_SINCE_MIDNIGHT_HI,omitempty"  xmp:"iXML:timestampSamplesSinceMidnightHI"`
	TimestampSamplesSinceMidnightLO int          `xml:"TIMESTAMP_SAMPLES_SINCE_MIDNIGHT_LO,omitempty"  xmp:"iXML:timestampSamplesSinceMidnightLO"`
	TimestampSampleRate             int          `xml:"TIMESTAMP_SAMPLE_RATE,omitempty"                xmp:"iXML:timestampSampleRate"`
}

func (x Speed) XmpTimecode() xmpdm.Timecode {
	samples := int64(x.TimestampSamplesSinceMidnightHI)<<32 | int64(x.TimestampSamplesSinceMidnightLO)
	samples = samples * 1000
	rate := int64(x.TimestampSampleRate) * 1000
	if rate == 0 {
		rate = int64(x.TimecodeRate.Value() * 1000)
	}
	sec := samples / rate
	h := sec / 3600
	m := (sec - h*3600) / 60
	s := sec - h*3600 - m*60
	f := int((samples % rate) * x.TimecodeRate.Num / rate / 1000)
	isDrop := x.TimecodeFlag == TimecodeFlagDF
	return xmpdm.Timecode{
		Format:      xmpdm.ParseTimecodeFormat(x.TimecodeRate.Value(), isDrop),
		H:           int(h),
		M:           int(m),
		S:           int(s),
		F:           f,
		IsDropFrame: isDrop,
	}
}

type History struct {
	XMLName          xml.Name `xml:"HISTORY"                      xmp:"-"`
	OriginalFilename string   `xml:"ORIGINAL_FILENAME,omitempty"  xmp:"iXML:originalFilename"`
	ParentFilename   string   `xml:"PARENT_FILENAME,omitempty"    xmp:"iXML:parentFilename"`
	ParentUID        string   `xml:"PARENT_UID,omitempty"         xmp:"iXML:parentUID"`
}

type FileSet struct {
	XMLName      xml.Name `xml:"FILE_SET"                  xmp:"-"`
	TotalFiles   int      `xml:"TOTAL_FILES,omitempty"     xmp:"iXML:totalFiles"`
	FamilyUID    string   `xml:"FAMILY_UID,omitempty"      xmp:"iXML:familyUID"`
	FamilyName   string   `xml:"FAMILY_NAME,omitempty"     xmp:"iXML:familyName"`
	FileSetIndex string   `xml:"FILE_SET_INDEX,omitempty"  xmp:"iXML:fileSetIndex"`
}

type Track struct {
	XMLName         xml.Name     `xml:"TRACK"                       xmp:"-"`
	ChannelIndex    int          `xml:"CHANNEL_INDEX,omitempty"     xmp:"iXML:channelIndex"`
	InterleaveIndex int          `xml:"INTERLEAVE_INDEX,omitempty"  xmp:"iXML:interleaveIndex"`
	Name            string       `xml:"NAME,omitempty"              xmp:"iXML:name"`
	Function        FunctionType `xml:"FUNCTION,omitempty"          xmp:"iXML:function"`
	DefaultMix      DefaultMix   `xml:"DEFAULT_MIX,omitempty"       xmp:"iXML:defaultMix"`
}

type TrackList []Track

func (x *TrackList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	depth := 1
	for depth > 0 {
		t, err := d.Token()
		if err != nil {
			return fmt.Errorf("%s: %v", NsIXML.GetName(), err)
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local != "TRACK" {
				d.Skip()
				continue
			}
			tt := Track{}
			if err := d.DecodeElement(&tt, &t); err != nil {
				return fmt.Errorf("%s: %v", NsIXML.GetName(), err)
			}
			*x = append(*x, tt)

		case xml.EndElement:
			depth--
		}
	}
	return nil
}

func (x TrackList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(x) == 0 {
		return nil
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	cToken := xml.StartElement{Name: xml.Name{Local: "TRACK_COUNT"}}
	if err := e.EncodeElement(len(x), cToken); err != nil {
		return err
	}
	for _, v := range x {
		if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "TRACK"}}); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return e.Flush()
}

func (x TrackList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x TrackList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *TrackList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type UserData struct {
	XMLName                   xml.Name `xml:"USER"                                    xmp:"-"`
	Comment                   string   `xml:",chardata"                               xmp:"iXML:comment"`
	FullTitle                 string   `xml:"FULL_TITLE,omitempty"                    xmp:"iXML:fullTitle"`
	DirectorName              string   `xml:"DIRECTOR_NAME,omitempty"                 xmp:"iXML:directorName"`
	ProductionName            string   `xml:"PRODUCTION_NAME,omitempty"               xmp:"iXML:productionName"`
	ProductionAddress         string   `xml:"PRODUCTION_ADDRESS,omitempty"            xmp:"iXML:productionAddress"`
	ProductionEmail           string   `xml:"PRODUCTION_EMAIL,omitempty"              xmp:"iXML:productionEmail"`
	ProductionPhone           string   `xml:"PRODUCTION_PHONE,omitempty"              xmp:"iXML:productionPhone"`
	ProductionNote            string   `xml:"PRODUCTION_NOTE,omitempty"               xmp:"iXML:productionNote"`
	SoundMixerName            string   `xml:"SOUND_MIXER_NAME,omitempty"              xmp:"iXML:soundMixerName"`
	SoundMixerAddress         string   `xml:"SOUND_MIXER_ADDRESS,omitempty"           xmp:"iXML:soundMixerAddress"`
	SoundMixerEmail           string   `xml:"SOUND_MIXER_EMAIL,omitempty"             xmp:"iXML:soundMixerEmail"`
	SoundMixerPhone           string   `xml:"SOUND_MIXER_PHONE,omitempty"             xmp:"iXML:soundMixerPhone"`
	SoundMixerNote            string   `xml:"SOUND_MIXER_NOTE,omitempty"              xmp:"iXML:soundMixerNote"`
	AudioRecorderModel        string   `xml:"AUDIO_RECORDER_MODEL,omitempty"          xmp:"iXML:audioRecorderModel"`
	AudioRecorderSerialNumber string   `xml:"AUDIO_RECORDER_SERIAL_NUMBER,omitempty"  xmp:"iXML:audioRecorderSerialNumber"`
	AudioRecorderFirmware     string   `xml:"AUDIO_RECORDER_FIRMWARE,omitempty"       xmp:"iXML:audioRecorderFirmware"`
}

type Location struct {
	XMLName          xml.Name         `xml:"LOCATION"                     xmp:"-"`
	LocationName     string           `xml:"LOCATION_NAME,omitempty"      xmp:"iXML:locationName"`     //  Human readable description of location
	LocationGps      string           `xml:"LOCATION_GPS,omitempty"       xmp:"iXML:locationGps"`      //  47.756787, -123.729977
	LocationAltitude string           `xml:"LOCATION_ALTITUDE,omitempty"  xmp:"iXML:locationAltitude"` //
	LocationType     LocationTypeList `xml:"LOCATION_TYPE,omitempty"      xmp:"iXML:locationType"`     //  [dictionary]
	LocationTime     LocationTimeList `xml:"LOCATION_TIME,omitempty"      xmp:"iXML:locationTime"`     //  [dictionary]
}

// extended attribute for tracks: http://www.gallery.co.uk/ixml/defaultmix.html
type DefaultMix struct {
	XMLName xml.Name `xml:"DEFAULT_MIX"      xmp:"-"`
	Level   string   `xml:"LEVEL,omitempty"  xmp:"iXML:level"` // Gain in dB as float, 'OFF' is an allowed value (e.g. -3.5)
	Pan     string   `xml:"PAN,omitempty"    xmp:"iXML:pan"`   // Pan angle in degrees L or R of centre (e.g. 23.5R)
}
