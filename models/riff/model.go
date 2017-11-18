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

// Package riff implements metadata for AVI and WAV files as defined by XMP Specification Part 3.
package riff

import (
	"fmt"
	"trimmer.io/go-xmp/models/dc"
	"trimmer.io/go-xmp/models/xmp_base"
	"trimmer.io/go-xmp/models/xmp_dm"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsRiff = xmp.NewNamespace("riffinfo", "http://ns.adobe.com/riff/info", NewModel)
)

func init() {
	xmp.Register(NsRiff, xmp.SoundMetadata)
}

func NewModel(name string) xmp.Model {
	return &RiffInfo{}
}

func MakeModel(d *xmp.Document) (*RiffInfo, error) {
	m, err := d.MakeModel(NsRiff)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*RiffInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *RiffInfo {
	if m := d.FindModel(NsRiff); m != nil {
		return m.(*RiffInfo)
	}
	return nil
}

type RiffInfo struct {
	ArchiveLocation string        `riffinfo:"IARL" xmp:"riffinfo:archivalLocation"`
	Artist          string        `riffinfo:"IART" xmp:"xmpDM:artist"`
	CommissionedBy  StringArray   `riffinfo:"ICMS" xmp:"riffinfo:commissioned"` // list: semicolon+blank separated
	Comments        string        `riffinfo:"ICMT" xmp:"xmpDM:logComment"`
	Copyright       AltString     `riffinfo:"ICOP" xmp:"dc:rights"`      // list: semicolon+blank separated
	CreateDate      xmp.Date      `riffinfo:"ICRD" xmp:"xmp:CreateDate"` // YYYY-MM-DD
	Engineeer       string        `riffinfo:"IENG" xmp:"xmpDM:engineer"` // list: semicolon+blank separated
	Genre           string        `riffinfo:"IGNR" xmp:"xmpDM:genre"`
	Keywords        StringArray   `riffinfo:"IKEY" xmp:"dc:subject"` // list: semicolon+blank separated
	SourceMedium    string        `riffinfo:"IMED" xmp:"dc:source"`
	Title           AltString     `riffinfo:"INAM" xmp:"riffinfo:name"`
	Product         string        `riffinfo:"IPRD" xmp:"riffinfo:product"`
	Description     AltString     `riffinfo:"ISBJ" xmp:"dc:description"`
	Software        xmp.AgentName `riffinfo:"ISFT" xmp:"xmp:CreatorTool"`
	SourceCredit    string        `riffinfo:"ISRC" xmp:"riffinfo:source"`
	SourceType      StringArray   `riffinfo:"ISRF" xmp:"dc:type"`
	Technician      string        `riffinfo:"ITCH" xmp:"riffinfo:technician"`

	// other tags found in the wild
	Rated               string `riffinfo:"AGES" xmp:"riffinfo:Rated"`
	Comment             string `riffinfo:"CMNT" xmp:"riffinfo:Comment"`
	EncodedBy           string `riffinfo:"CODE" xmp:"riffinfo:EncodedBy"`
	Comments2           string `riffinfo:"COMM" xmp:"riffinfo:Comments"`
	Directory           string `riffinfo:"DIRC" xmp:"riffinfo:Directory"`
	SoundSchemeTitle    string `riffinfo:"DISP" xmp:"riffinfo:SoundSchemeTitle"`
	DateTimeOriginal    string `riffinfo:"DTIM" xmp:"riffinfo:DateTimeOriginal"`
	Genre2              string `riffinfo:"GENR" xmp:"riffinfo:Genre"`
	ArchivalLocation    string `riffinfo:"IARL" xmp:"riffinfo:ArchivalLocation"`
	FirstLanguage       string `riffinfo:"IAS1" xmp:"riffinfo:FirstLanguage"`
	SecondLanguage      string `riffinfo:"IAS2" xmp:"riffinfo:SecondLanguage"`
	ThirdLanguage       string `riffinfo:"IAS3" xmp:"riffinfo:ThirdLanguage"`
	FourthLanguage      string `riffinfo:"IAS4" xmp:"riffinfo:FourthLanguage"`
	FifthLanguage       string `riffinfo:"IAS5" xmp:"riffinfo:FifthLanguage"`
	SixthLanguage       string `riffinfo:"IAS6" xmp:"riffinfo:SixthLanguage"`
	SeventhLanguage     string `riffinfo:"IAS7" xmp:"riffinfo:SeventhLanguage"`
	EighthLanguage      string `riffinfo:"IAS8" xmp:"riffinfo:EighthLanguage"`
	NinthLanguage       string `riffinfo:"IAS9" xmp:"riffinfo:NinthLanguage"`
	BaseURL             string `riffinfo:"IBSU" xmp:"riffinfo:BaseURL"`
	DefaultAudioStream  string `riffinfo:"ICAS" xmp:"riffinfo:DefaultAudioStream"`
	CostumeDesigner     string `riffinfo:"ICDS" xmp:"riffinfo:CostumeDesigner"`
	Commissioned        string `riffinfo:"ICMS" xmp:"riffinfo:Commissioned"`
	Cinematographer     string `riffinfo:"ICNM" xmp:"riffinfo:Cinematographer"`
	Country             string `riffinfo:"ICNT" xmp:"riffinfo:Country"`
	Cropped             string `riffinfo:"ICRP" xmp:"riffinfo:Cropped"`
	Dimensions          string `riffinfo:"IDIM" xmp:"riffinfo:Dimensions"`
	DateTimeOriginal2   string `riffinfo:"IDIT" xmp:"-"`
	DotsPerInch         string `riffinfo:"IDPI" xmp:"riffinfo:DotsPerInch"`
	DistributedBy       string `riffinfo:"IDST" xmp:"riffinfo:DistributedBy"`
	EditedBy            string `riffinfo:"IEDT" xmp:"riffinfo:EditedBy"`
	EncodedBy2          string `riffinfo:"IENC" xmp:"-"`
	Lightness           string `riffinfo:"ILGT" xmp:"riffinfo:Lightness"`
	LogoURL             string `riffinfo:"ILGU" xmp:"riffinfo:LogoURL"`
	LogoIconURL         string `riffinfo:"ILIU" xmp:"riffinfo:LogoIconURL"`
	Language            string `riffinfo:"ILNG" xmp:"riffinfo:Language"`
	MoreInfoBannerImage string `riffinfo:"IMBI" xmp:"riffinfo:MoreInfoBannerImage"`
	MoreInfoBannerURL   string `riffinfo:"IMBU" xmp:"riffinfo:MoreInfoBannerURL"`
	MoreInfoText        string `riffinfo:"IMIT" xmp:"riffinfo:MoreInfoText"`
	MoreInfoURL         string `riffinfo:"IMIU" xmp:"riffinfo:MoreInfoURL"`
	MusicBy             string `riffinfo:"IMUS" xmp:"riffinfo:MusicBy"`
	ProductionDesigner  string `riffinfo:"IPDS" xmp:"riffinfo:ProductionDesigner"`
	NumColors           string `riffinfo:"IPLT" xmp:"riffinfo:NumColors"`
	ProducedBy          string `riffinfo:"IPRO" xmp:"riffinfo:ProducedBy"`
	RippedBy            string `riffinfo:"IRIP" xmp:"riffinfo:RippedBy"`
	Rating              string `riffinfo:"IRTD" xmp:"riffinfo:Rating"`
	SecondaryGenre      string `riffinfo:"ISGN" xmp:"riffinfo:SecondaryGenre"`
	Sharpness           string `riffinfo:"ISHP" xmp:"riffinfo:Sharpness"`
	TimeCode            string `riffinfo:"ISMP" xmp:"riffinfo:TimeCode"`
	ProductionStudio    string `riffinfo:"ISTD" xmp:"riffinfo:ProductionStudio"`
	Starring            string `riffinfo:"ISTR" xmp:"riffinfo:Starring"`
	WatermarkURL        string `riffinfo:"IWMU" xmp:"riffinfo:WatermarkURL"`
	WrittenBy           string `riffinfo:"IWRI" xmp:"riffinfo:WrittenBy"`
	Language2           string `riffinfo:"LANG" xmp:"riffinfo:Language2"`
	Location            string `riffinfo:"LOCA" xmp:"riffinfo:Location"`
	Part                string `riffinfo:"PRT1" xmp:"riffinfo:Part"`
	NumberOfParts       string `riffinfo:"PRT2" xmp:"riffinfo:NumberOfParts"`
	Rate                string `riffinfo:"RATE" xmp:"riffinfo:Rate"`
	Starring2           string `riffinfo:"STAR" xmp:"-"`
	Statistics          string `riffinfo:"STAT" xmp:"riffinfo:Statistics"`
	TapeName            string `riffinfo:"TAPE" xmp:"riffinfo:TapeName"`
	EndTimecode         string `riffinfo:"TCDO" xmp:"riffinfo:EndTimecode"`
	StartTimecode       string `riffinfo:"TCOD" xmp:"riffinfo:StartTimecode"`
	Title2              string `riffinfo:"TITL" xmp:"riffinfo:Title"`
	Length              string `riffinfo:"TLEN" xmp:"riffinfo:Length"`
	Organization        string `riffinfo:"TORG" xmp:"riffinfo:Organization"`
	TrackNumber         string `riffinfo:"TRCK" xmp:"riffinfo:TrackNumber"`
	URL                 string `riffinfo:"TURL" xmp:"riffinfo:URL"`
	Version             string `riffinfo:"TVER" xmp:"riffinfo:Version"`
	VegasVersionMajor   string `riffinfo:"VMAJ" xmp:"riffinfo:VegasVersionMajor"`
	VegasVersionMinor   string `riffinfo:"VMIN" xmp:"riffinfo:VegasVersionMinor"`
	Year                string `riffinfo:"YEAR" xmp:"riffinfo:Year"`
}

func (m *RiffInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsRiff}
}

func (m *RiffInfo) Can(nsName string) bool {
	return nsName == NsRiff.GetName()
}

func (x *RiffInfo) SyncFromXMP(d *xmp.Document) error {
	if m := dc.FindModel(d); m != nil {
		if len(m.Rights) > 0 {
			x.Copyright = AltString(m.Rights)
		}
		if len(m.Subject) > 0 {
			x.Keywords = StringArray(m.Subject)
		}
		if len(m.Type) > 0 {
			x.SourceType = StringArray(m.Type)
		}
		if len(m.Description) > 0 {
			x.Description = AltString(m.Description)
		}
		if len(m.Title) > 0 {
			x.Title = AltString(m.Title)
		}
		x.SourceMedium = m.Source
	}
	if base := xmpbase.FindModel(d); base != nil {
		if !base.CreateDate.IsZero() {
			x.CreateDate = base.CreateDate
		}
		x.Software = base.CreatorTool
	}
	if dm := xmpdm.FindModel(d); dm != nil {
		if dm.Artist != "" {
			x.Artist = dm.Artist
		}
		if dm.LogComment != "" {
			x.Comments = dm.LogComment
		}
		if dm.Genre != "" {
			x.Genre = dm.Genre
		}
		if dm.Engineer != "" {
			x.Engineeer = dm.Engineer
		}
	}
	return nil
}

func (x *RiffInfo) SyncToXMP(d *xmp.Document) error {
	if x.DateTimeOriginal2 != "" && x.DateTimeOriginal == "" {
		x.DateTimeOriginal = x.DateTimeOriginal2
	}
	if x.EncodedBy2 != "" && x.EncodedBy == "" {
		x.EncodedBy = x.EncodedBy2
	}
	if x.Starring2 != "" && x.Starring == "" {
		x.Starring = x.Starring2
	}
	// manually add INAM to dc:title["x-default"] if no title exists
	m, err := dc.MakeModel(d)
	if err != nil {
		return err
	}
	if len(m.Title) == 0 {
		m.Title = xmp.AltString(x.Title)
	}
	return nil
}

func (x *RiffInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *RiffInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsRiff.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *RiffInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsRiff.GetName(), err)
	}
	return nil
}

// Lists all non-empty tags.
func (x *RiffInfo) ListTags() (xmp.TagList, error) {
	if l, err := xmp.ListNativeFields(x); err != nil {
		return nil, fmt.Errorf("%s: %v", NsRiff.GetName(), err)
	} else {
		return l, nil
	}
}
