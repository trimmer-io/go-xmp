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

package qt

import (
	"fmt"
	"github.com/echa/go-xmp/xmp"
)

// QuickTime User Data as written by the "udta" handler using FourCC atom names.
// © \251 \xA9 Tags are multi-language versions
//
// Tag ID's beginning with the copyright symbol (hex 0xa9) are multi-language text.
//
// https://developer.apple.com/library/content/documentation/QuickTime/QTFF/QTFFChap2/qtff2.html
type QtUserdata struct {

	// Quicktime tags exported on MacOS 10.11 SDK (some also used by iTunes)
	Album                string    `qt:"©alb" xmp:"qt:Album"`                // QT, iTunes
	Arranger             string    `qt:"©arg" xmp:"qt:Arranger"`             // QT, iTunes
	Artist               string    `qt:"©ART" xmp:"qt:Artist"`               // QT, iTunes
	Author               string    `qt:"©aut" xmp:"qt:Author"`               // QT, iTunes
	Chapter              string    `qt:"©chp" xmp:"qt:Chapter"`              // QT-only
	Comment              string    `qt:"©cmt" xmp:"qt:Comment"`              // QT, iTunes
	Composer             string    `qt:"©com" xmp:"qt:Composer"`             // QT-only
	Copyright            string    `qt:"©cpy" xmp:"qt:Copyright"`            // QT-only
	ReleaseDate          xmp.Date  `qt:"©day" xmp:"qt:ReleaseDate"`          // QT, iTunes
	Description          string    `qt:"©des" xmp:"qt:Description"`          // QT, iTunes
	Director             string    `qt:"©dir" xmp:"qt:Director"`             // QT, iTunes
	Disclaimer           string    `qt:"©dis" xmp:"qt:Disclaimer"`           // QT-only
	EncodedBy            string    `qt:"©enc" xmp:"qt:EncodedBy"`            // QT, iTunes
	Title                string    `qt:"©nam" xmp:"qt:Title"`                // QT, iTunes
	GenreName            string    `qt:"©gen" xmp:"qt:GenreName"`            // QT, iTunes
	HostComputer         string    `qt:"©hst" xmp:"qt:HostComputer"`         // QT-only
	Information          string    `qt:"©inf" xmp:"qt:Information"`          // QT-only
	Keywords             string    `qt:"©key" xmp:"qt:Keywords"`             // QT-only
	FileCreator          string    `qt:"©mak" xmp:"qt:FileCreator"`          // QT, iTunes (RecordCompany)
	FileCreatorModel     string    `qt:"©mod" xmp:"qt:FileCreatorModel"`     // QT-only
	OriginalArtist       string    `qt:"©ope" xmp:"qt:OriginalArtist"`       // iTunes
	FileFormat           string    `qt:"©fmt" xmp:"qt:FileFormat"`           // QT-only
	Credits              string    `qt:"©src" xmp:"qt:Credits"`              // QT, iTunes
	Performer            string    `qt:"©prf" xmp:"qt:Performer"`            // QT, iTunes
	Producer             string    `qt:"©prd" xmp:"qt:Producer"`             // QT, iTunes
	Publisher            string    `qt:"©pub" xmp:"qt:Publisher"`            // QT, iTunes
	Product              string    `qt:"©PRD" xmp:"qt:Product"`              // QT-only
	FileCreatorSoftware  string    `qt:"©swr" xmp:"qt:FileCreatorSoftware"`  // QT-only
	PlaybackRequirements string    `qt:"©req" xmp:"qt:PlaybackRequirements"` // QT-only
	Track                string    `qt:"©trk" xmp:"qt:Track"`                // QT-only
	CopyWarning          string    `qt:"©wrn" xmp:"qt:CopyWarning"`          // QT-only
	Writer               string    `qt:"©wrt" xmp:"qt:Writer"`               // QT, iTunes (Composer)
	Url                  xmp.Url   `qt:"©url" xmp:"qt:Url"`                  // QT, iTunes (Online Extras)
	LocationGPS          *Location `qt:"©xyz" xmp:"qt:LocationGPS"`          // GPS latitude+longitude+altitude
	TrackName            string    `qt:"tnam" xmp:"qt:TrackName"`            // QT-only
	PhonogramRights      string    `qt:"©phg" xmp:"qt:PhonogramRights"`      // QT, iTunes
	DisplayName          string    `qt:"name" xmp:"qt:DisplayName"`          // QT-only
	TaggedCharacteristic string    `qt:"tagc" xmp:"qt:TaggedCharacteristic"` // QT, ISO

	// QT only (administrative metadata)
	AudioBookReleaseDate xmp.Date `qt:"rldt" xmp:"qt:AudioBookReleaseDate"`
	ClipFileName         string   `qt:"clfn" xmp:"qt:ClipFileName"`
	ClipID               string   `qt:"clid" xmp:"qt:ClipID"`
	ContentDistributorID string   `qt:"cdis" xmp:"qt:ContentDistributorID"`
	ContentID            string   `qt:"ccid" xmp:"qt:ContentID"`
	CreationDate         xmp.Date `qt:"date" xmp:"qt:CreationDate"`
	Grouping             string   `qt:"©grp" xmp:"qt:Grouping"`
	GUID                 string   `qt:"GUID" xmp:"qt:GUID"`
	ISRCCode             string   `qt:"©isr" xmp:"qt:ISRCCode"`

	// QT only (descriptive metadata)
	AlbumArtist        string  `qt:"albr" xmp:"qt:AlbumArtist"`
	Angle              string  `qt:"angl" xmp:"qt:Angle"`
	ArrangerKeywords   string  `qt:"©ark" xmp:"qt:ArrangerKeywords"`
	AudibleTags        string  `qt:"tags" xmp:"qt:AudibleTags"`
	CameraID           string  `qt:"cmid" xmp:"qt:CameraID"`
	CameraManufacturer string  `qt:"manu" xmp:"qt:CameraManufacturer"`
	CameraModel        string  `qt:"modl" xmp:"qt:CameraModel"`
	CameraName         string  `qt:"cmnm" xmp:"qt:CameraName"`
	CameraSerialNumber string  `qt:"slno" xmp:"qt:CameraSerialNumber"`
	ChapterList        string  `qt:"chpl" xmp:"qt:ChapterList"`
	ComposerKeywords   string  `qt:"©cok" xmp:"qt:ComposerKeywords"`
	Edit1              string  `qt:"©ed1" xmp:"qt:Edit1"`
	Edit2              string  `qt:"©ed2" xmp:"qt:Edit2"`
	Edit3              string  `qt:"©ed3" xmp:"qt:Edit3"`
	Edit4              string  `qt:"©ed4" xmp:"qt:Edit4"`
	Edit5              string  `qt:"©ed5" xmp:"qt:Edit5"`
	Edit6              string  `qt:"©ed6" xmp:"qt:Edit6"`
	Edit7              string  `qt:"©ed7" xmp:"qt:Edit7"`
	Edit8              string  `qt:"©ed8" xmp:"qt:Edit8"`
	Edit9              string  `qt:"©ed9" xmp:"qt:Edit9"`
	FileCreatorModel2  string  `qt:"©mdl" xmp:"qt:FileCreatorModel2"`
	FileCreatorUrl     xmp.Url `qt:"©mal" xmp:"qt:FileCreatorUrl"`
	Lyrics             string  `qt:"©lyr" xmp:"qt:Lyrics"`    // QT, iTunes
	LyricsUrl          xmp.Url `qt:"lrcu" xmp:"qt:LyricsUrl"` // iTunes
	OriginalFormat     string  `qt:"orif" xmp:"qt:OriginalFormat"`
	OriginalSource     string  `qt:"oris" xmp:"qt:OriginalSource"`
	PerformerKeywords  string  `qt:"©prk" xmp:"qt:PerformerKeywords"`
	PerformerUrl       xmp.Url `qt:"©prl" xmp:"qt:PerformerUrl"`
	ProducerKeywords   string  `qt:"©pdk" xmp:"qt:ProducerKeywords"`
	ProductVersion     string  `qt:"VERS" xmp:"qt:ProductVersion"`
	RecordLabelName    string  `qt:"©lab" xmp:"qt:RecordLabelName"`
	RecordLabelUrl     xmp.Url `qt:"©lal" xmp:"qt:RecordLabelUrl"`
	Reel               string  `qt:"reel" xmp:"qt:Reel"`
	Scene              string  `qt:"scen" xmp:"qt:Scene"`
	Shot               string  `qt:"shot" xmp:"qt:Shot"`
	SongWriter         string  `qt:"©swf" xmp:"qt:SongWriter"`
	SongWriterKeywords string  `qt:"©swk" xmp:"qt:SongWriterKeywords"`
	Subtitle           string  `qt:"©snm" xmp:"qt:Subtitle"`
	SubtitleKeywords   string  `qt:"©snk" xmp:"qt:SubtitleKeywords"`
	Synopsis           string  `qt:"ldes" xmp:"qt:Synopsis"`
	TitleKeywords      string  `qt:"©nak" xmp:"qt:TitleKeywords"`

	// QT mdta (not in official udta, but has FourCC code)
	CoverArt       string `qt:"covr" xmp:"qt:CoverArt"`
	CollectionUser string `qt:"coll" xmp:"qt:CollectionUser"`
	UserRating     string `qt:"rtng" xmp:"qt:UserRating"`
	RecordingYear  int    `qt:"yrrc" xmp:"qt:RecordingYear"`

	// QT only (technical metadata)
	ApertureMode   string  `qt:"apmd" xmp:"qt:ApertureMode"`
	FlightPitch    float64 `qt:"©fpt" xmp:"qt:FlightPitch"`
	FlightRoll     float64 `qt:"©frl" xmp:"qt:FlightRoll"`
	FlightSpeedX   float64 `qt:"©xsp" xmp:"qt:FlightSpeedX"`
	FlightSpeedY   float64 `qt:"©ysp" xmp:"qt:FlightSpeedY"`
	FlightSpeedZ   float64 `qt:"©zsp" xmp:"qt:FlightSpeedZ"`
	FlightYaw      float64 `qt:"©fyw" xmp:"qt:FlightYaw"`
	GimbalPitch    float64 `qt:"©gpt" xmp:"qt:GimbalPitch"`
	GimbalRoll     float64 `qt:"©grl" xmp:"qt:GimbalRoll"`
	GimbalYaw      float64 `qt:"©gyw" xmp:"qt:GimbalYaw"`
	HintInfo       string  `qt:"hnti" xmp:"qt:HintInfo"`
	HintTrackInfo  string  `qt:"hinf" xmp:"qt:HintTrackInfo"`
	HintVersion    string  `qt:"hinv" xmp:"qt:HintVersion"`
	PrintToVideo   string  `qt:"ptv " xmp:"qt:PrintToVideo"`
	TrackType      string  `qt:"kgtt" xmp:"qt:TrackType"`
	WindowLocation Point   `qt:"WLOC" xmp:"qt:WindowLocation"`

	// embedded XML documents
	// IXML string `qt:"iXML" xmp:"qt:ixml" json:"qt:iXML,omitempty"`
	// XMP  string `qt:"XMP_" xmp:"-" json:"-"`

	//
	// Vendor-specific tags
	//
	// Google transcoded videos (YouTube)
	GoogleHostHeader    string `qt:"gshh" xmp:"qt:GoogleHostHeader"`
	GooglePingMessage   string `qt:"gspm" xmp:"qt:GooglePingMessage"`
	GooglePingURL       string `qt:"gspu" xmp:"qt:GooglePingURL"`
	GoogleSourceData    string `qt:"gssd" xmp:"qt:GoogleSourceData"`
	GoogleStartTime     string `qt:"gsst" xmp:"qt:GoogleStartTime"`
	GoogleTrackDuration string `qt:"gstd" xmp:"qt:GoogleTrackDuration"`

	// Microsoft http://www.sno.phy.queensu.ca/~phil/exiftool/TagNames/Microsoft.html
	// MicrosoftExtra   MicrosoftExtra `qt:"Xtra" xmp:"msft:Extra" json:"-"`

	// Canon Cameras
	CanonCodec     string `qt:"CNCV" xmp:"qt:CanonCodec"`    // "CanonAVC0002"
	CanonModel     string `qt:"CNMN" xmp:"qt:CanonModel"`    // "Canon EOS 5D Mark II"
	CanonFirmware  string `qt:"CNFV" xmp:"qt:CanonFirmware"` // "Firmware Version 2.1.1"
	CanonThumbnail []byte `qt:"CNDA" xmp:"qt:CanonThumbnail"`
}

func (m *QtUserdata) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsQuicktime}
}

func (m *QtUserdata) Can(nsName string) bool {
	return nsName == NsQuicktime.GetName()
}

func (x *QtUserdata) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x QtUserdata) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *QtUserdata) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *QtUserdata) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *QtUserdata) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	}
	return nil
}
