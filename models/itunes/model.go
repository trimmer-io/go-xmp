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

// iTunes-style metadata as found in .mp4, .m4a, .m4p, .m4v, .m4b files

// Package itunes implements metadata models found in Apple iTunes audio and video files.
package itunes

import (
	"fmt"

	"github.com/echa/go-xmp/xmp"
)

var (
	NsITunes = &xmp.Namespace{"iTunes", "http://ns.apple.com/itunes/1.0/", NewModel}
)

func init() {
	xmp.Register(NsITunes, xmp.MusicMetadata)
}

func NewModel(name string) xmp.Model {
	return &ITunesMetadata{}
}

func MakeModel(d *xmp.Document) (*ITunesMetadata, error) {
	m, err := d.MakeModel(NsITunes)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*ITunesMetadata)
	return x, nil
}

func FindModel(d *xmp.Document) *ITunesMetadata {
	if m := d.FindModel(NsITunes); m != nil {
		return m.(*ITunesMetadata)
	}
	return nil
}

func (m *ITunesMetadata) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsITunes}
}

func (m *ITunesMetadata) Can(nsName string) bool {
	return nsName == NsITunes.GetName()
}

func (x *ITunesMetadata) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x ITunesMetadata) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *ITunesMetadata) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *ITunesMetadata) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsITunes.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *ITunesMetadata) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsITunes.GetName(), err)
	}
	return nil
}

// Lists all non-empty tags.
func (x *ITunesMetadata) ListTags() (xmp.TagList, error) {
	if l, err := xmp.ListNativeFields(x); err != nil {
		return nil, fmt.Errorf("%s: %v", NsITunes.GetName(), err)
	} else {
		return l, nil
	}
}

// iTunes specific Quicktime metadata tags
// itms: iTunes storage format atom
// itlk: mdta-style atoms Reverse DNS (com.apple.itunes)
// itsk: udta-style FourCC atoms
type ITunesMetadata struct {
	AccountKind       AppleStoreAccountType `iTunes:"akID" xmp:"iTunes:AccountKind"`
	Acknowledgement   string                `iTunes:"©cak" xmp:"iTunes:Acknowledgement"`
	Album             string                `iTunes:"©alb" xmp:"iTunes:Album"`
	AlbumArtist       string                `iTunes:"aART" xmp:"iTunes:AlbumArtist"`
	AppleID           string                `iTunes:"apID" xmp:"iTunes:AppleID"`
	Arranger          string                `iTunes:"©arg" xmp:"iTunes:Arranger"`
	ArtDirector       string                `iTunes:"©ard" xmp:"iTunes:ArtDirector"`
	Artist            string                `iTunes:"©ART" xmp:"iTunes:Artist"`
	ArtistID          string                `iTunes:"atID" xmp:"iTunes:ArtistID"`
	Author            string                `iTunes:"©aut" xmp:"iTunes:Author"`
	BeatsPerMin       int                   `iTunes:"tmpo" xmp:"iTunes:BeatsPerMin"`
	Comment           string                `iTunes:"©cmt" xmp:"iTunes:Comment"`
	Composer          string                `iTunes:"©wrt" xmp:"iTunes:Composer"`
	Conductor         string                `iTunes:"©con" xmp:"iTunes:Conductor"`
	IsExplicit        RatingCode            `iTunes:"rtng" xmp:"iTunes:IsExplicit"`
	Copyright         string                `iTunes:"cprt" xmp:"iTunes:Copyright"`
	CoverArt          string                `iTunes:"covr" xmp:"iTunes:CoverArt"`
	CoverUrl          xmp.Url               `iTunes:"cvru" xmp:"iTunes:CoverUrl"`
	Credits           string                `iTunes:"©src" xmp:"iTunes:Credits"`
	Description       string                `iTunes:"©des" xmp:"iTunes:Description"`
	Director          string                `iTunes:"©dir" xmp:"iTunes:Director"`
	DiscNumber        xmp.Rational          `iTunes:"disk" xmp:"iTunes:DiscNumber"`
	Duration          int64                 `iTunes:"dcfD" xmp:"iTunes:Duration"`
	EncodedBy         string                `iTunes:"©enc" xmp:"iTunes:EncodedBy"`
	EncodingTool      string                `iTunes:"©too" xmp:"iTunes:EncodingTool"`
	EQ                string                `iTunes:"©equ" xmp:"iTunes:EQ"`
	ExecProducer      string                `iTunes:"©xpd" xmp:"iTunes:ExecProducer"`
	GenreCode         GenreCode             `iTunes:"gnre" xmp:"iTunes:GenreCode"` // Predefined, = ID3 genres
	GenreID           GenreID               `iTunes:"geID" xmp:"iTunes:GenreID"`
	GenreName         string                `iTunes:"©gen" xmp:"iTunes:GenreName"` // user defined
	Grouping          string                `iTunes:"grup" xmp:"iTunes:Grouping"`  // like TIT1 in ID3
	IconUrl           xmp.Url               `iTunes:"icnu" xmp:"iTunes:IconUrl"`
	InfoUrl           xmp.Url               `iTunes:"infu" xmp:"iTunes:InfoUrl"`
	IsDiscCompilation Bool                  `iTunes:"cpil" xmp:"iTunes:IsDiscCompilation"`
	IsGaplessPlayback PlayGapMode           `iTunes:"pgap" xmp:"iTunes:IsGaplessPlayback"`
	IsHDVideo         Bool                  `iTunes:"hdvd" xmp:"iTunes:IsHDVideo"`
	IsiTunesU         string                `iTunes:"itnu" xmp:"iTunes:IsiTunesU"`
	IsPodcast         Bool                  `iTunes:"pcst" xmp:"iTunes:IsPodcast"`
	Keywords          string                `iTunes:"keyw" xmp:"iTunes:Keywords"`
	LinerNotes        string                `iTunes:"©lnt" xmp:"iTunes:LinerNotes"`
	Lyrics            string                `iTunes:"©lyr" xmp:"iTunes:Lyrics"`
	LyricsUrl         string                `iTunes:"lrcu" xmp:"iTunes:LyricsUrl"`
	MediaType         MediaType             `iTunes:"stik" xmp:"iTunes:MediaType"`
	Narrator          string                `iTunes:"©nrt" xmp:"iTunes:Narrator"`
	OnlineExtras      string                `iTunes:"©url" xmp:"iTunes:OnlineExtras"`
	OriginalArtist    string                `iTunes:"©ope" xmp:"iTunes:OriginalArtist"`
	Performer         string                `iTunes:"©prf" xmp:"iTunes:Performer"`
	PhonogramRights   string                `iTunes:"©phg" xmp:"iTunes:PhonogramRights"`
	PlaylistID        string                `iTunes:"plID" xmp:"iTunes:PlaylistID"`
	PodcastCategory   string                `iTunes:"catg" xmp:"iTunes:PodcastCategory"`
	PodcastGuid       string                `iTunes:"egid" xmp:"iTunes:PodcastGuid"`
	PodcastUrl        string                `iTunes:"purl" xmp:"iTunes:PodcastUrl"`
	Producer          string                `iTunes:"©prd" xmp:"iTunes:Producer"`
	ProductID         string                `iTunes:"prID" xmp:"iTunes:ProductID"`
	Publisher         string                `iTunes:"©pub" xmp:"iTunes:Publisher"`
	PurchaseDate      xmp.Date              `iTunes:"purd" xmp:"iTunes:PurchaseDate"`
	RatingPercent     string                `iTunes:"rate" xmp:"iTunes:RatingPercent"`
	RecordCompany     string                `iTunes:"©mak" xmp:"iTunes:RecordCompany"`
	ReleaseDate       xmp.Date              `iTunes:"©day" xmp:"iTunes:ReleaseDate"`
	ShowEpisodeName   string                `iTunes:"tves" xmp:"iTunes:ShowEpisodeName"`
	ShowEpisodeNum    int                   `iTunes:"tven" xmp:"iTunes:ShowEpisodeNum"`
	ShowName          string                `iTunes:"tvsh" xmp:"iTunes:ShowName"`
	ShowSeasonNum     int                   `iTunes:"tvsn" xmp:"iTunes:ShowSeasonNum"`
	Soloist           string                `iTunes:"©sol" xmp:"iTunes:Soloist"`
	SongID            string                `iTunes:"cnID" xmp:"iTunes:SongID"` // content ID, AppleStoreCatalogID
	SortAlbum         string                `iTunes:"soal" xmp:"iTunes:SortAlbum"`
	SortAlbumArtist   string                `iTunes:"soaa" xmp:"iTunes:SortAlbumArtist"`
	SortArtist        string                `iTunes:"soar" xmp:"iTunes:SortArtist"`
	SortComposer      string                `iTunes:"soco" xmp:"iTunes:SortComposer"`
	SortName          string                `iTunes:"sonm" xmp:"iTunes:SortName"`
	SortShow          string                `iTunes:"sosn" xmp:"iTunes:SortShow"`
	SoundEngineer     string                `iTunes:"©sne" xmp:"iTunes:SoundEngineer"`
	StoreFrontID      string                `iTunes:"sfID" xmp:"iTunes:StoreFrontID"` // apple store country
	Synopsis          string                `iTunes:"ldes" xmp:"iTunes:Synopsis"`
	Thanks            string                `iTunes:"©thx" xmp:"iTunes:Thanks"`
	Title             string                `iTunes:"©nam" xmp:"iTunes:Title"`
	ToolInfo          string                `iTunes:"tool" xmp:"iTunes:ToolInfo"`
	TrackNumber       int                   `iTunes:"trkn" xmp:"iTunes:TrackNumber"`
	TrackSubTitle     string                `iTunes:"©st3" xmp:"iTunes:TrackSubTitle"`
	TVNetworkName     string                `iTunes:"tvnn" xmp:"iTunes:TVNetworkName"`
	XID               string                `iTunes:"xid " xmp:"iTunes:XID"`

	ContentRating      *ContentRating `iTunes:"iTunEXTC"                xmp:"iTunes:ContentRating"`
	SoundCheck         []byte         `iTunes:"iTunNORM"                xmp:"iTunes:SoundCheck"`
	SMPB               *SMPB          `iTunes:"iTunSMPB"                xmp:"iTunes:SMPB"`
	IsGaplessPlayback2 Bool           `iTunes:"iTunPGAP"                xmp:"iTunes:IsGaplessPlayback2"`
	MovieInfo          *MovieInfo     `iTunes:"iTunMOVI"                xmp:"iTunes:MovieInfo"`
	CDDBToc            string         `iTunes:"iTunes_CDDB_1"           xmp:"iTunes:CDDBToc"`
	CDDBTrackNumber    string         `iTunes:"iTunes_CDDB_TrackNumber" xmp:"iTunes:CDDBTrackNumber"`
	CDDBMediaID        string         `iTunes:"iTunes_CDDB_IDs"         xmp:"iTunes:CDDBMediaID"`
	EncodingParams     []byte         `iTunes:"Encoding Params"         xmp:"iTunes:EncodingParams"`

	Extension xmp.TagList `iTunes:",any" xmp:"iTunes:extension"`
}
