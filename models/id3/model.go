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

// ID3 v2.2, v2.3 & v2.4
// http://id3.org/

// TODO:
// - almost all TextUnmarshalers are unimplemented due to missing samples
// - embedded XMP may be carried inside a PRIV frame with XMP owner identifier
// - multi-language text fields (requires samples)
// - multi-lang comments using CommentArray type (requires samples)

// FIXME: Version 2.4 of the specification prescribes that all text fields
// (the fields that start with a T, except for TXXX) can contain multiple
// values separated by a null character. The null character varies by
// character encoding.

// Package id3 implements the ID3v2.4 metadata standard for audio files.
package id3

import (
	"fmt"
	"strings"
	"time"

	"trimmer.io/go-xmp/models/dc"
	"trimmer.io/go-xmp/models/xmp_base"
	"trimmer.io/go-xmp/models/xmp_dm"
	"trimmer.io/go-xmp/models/xmp_rights"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsID3 = xmp.NewNamespace("id3", "http://id3.org/ns/2.4/", NewModel)
)

func init() {
	xmp.Register(NsID3, xmp.MusicMetadata)
}

func NewModel(name string) xmp.Model {
	if name == "id3" {
		return &ID3{}
	}
	return nil
}

func MakeModel(d *xmp.Document) (*ID3, error) {
	m, err := d.MakeModel(NsID3)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*ID3)
	return x, nil
}

func FindModel(d *xmp.Document) *ID3 {
	if m := d.FindModel(NsID3); m != nil {
		return m.(*ID3)
	}
	return nil
}

func (m *ID3) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsID3}
}

func (m *ID3) Can(nsName string) bool {
	return nsName == NsID3.GetName()
}

func (x *ID3) SyncFromXMP(d *xmp.Document) error {
	if m := dc.FindModel(d); m != nil {
		if len(m.Title) > 0 {
			x.TitleDescription = m.Title.Default()
		}
		if len(m.Rights) > 0 {
			x.Copyright = m.Rights.Default()
		}
	}
	if base := xmpbase.FindModel(d); base != nil {
		if !base.CreateDate.IsZero() {
			x.RecordingTime = base.CreateDate
		}
	}
	if dm := xmpdm.FindModel(d); dm != nil {
		if dm.Artist != "" {
			x.LeadPerformer = dm.Artist
		}
		if dm.Album != "" {
			x.AlbumTitle = dm.Album
		}
		if dm.LogComment != "" {
			x.Comments = dm.LogComment
		}
		if dm.Genre != "" {
			if err := x.ContentType.UnmarshalText([]byte(dm.Genre)); err != nil {
				return err
			}
		}
		if dm.TrackNumber > 0 {
			x.TrackNumber = TrackNum{
				Track: dm.TrackNumber,
			}
		}
		if dm.PartOfCompilation.Value() {
			x.IsCompilation = True
		} else {
			x.IsCompilation = False
		}
		if dm.Composer != "" {
			x.Composer = dm.Composer
		}
		if dm.Engineer != "" {
			x.ModifiedBy = dm.Engineer
		}
		if dm.DiscNumber > 0 {
			x.PartOfASet = TrackNum{
				Track: dm.DiscNumber,
			}
		}
		if dm.Lyrics != "" {
			x.UnsynchronizedLyrics = strings.Split(dm.Lyrics, "\n")
		}
	}
	if rights := xmprights.FindModel(d); rights != nil {
		if rights.WebStatement != "" {
			x.CopyrightInformation = xmp.Url(rights.WebStatement)
		}
	}

	return nil
}

func (x ID3) SyncToXMP(d *xmp.Document) error {
	// pick default language for multi-lang XMP properties
	lang := x.Language

	// DC gets title (TIT2) and copyright (TCOP)
	if x.TitleDescription != "" || x.Copyright != "" {
		m, err := dc.MakeModel(d)
		if err != nil {
			return err
		}
		if x.TitleDescription != "" && len(m.Title) == 0 {
			m.Title.AddDefault(lang, x.TitleDescription)
		}
		if x.Copyright != "" && len(m.Rights) == 0 {
			m.Rights.AddDefault(lang, x.Copyright)
		}
	}

	// XMP base gets creation date
	if x.RecordingTime.IsZero() && !x.Date_v23.IsZero() && !x.Time_v23.IsZero() {
		_, mm, dd := x.Date_v23.Value().Date()
		h, m, s := x.Time_v23.Value().Clock()
		x.RecordingTime = xmp.Date(time.Date(x.Year_v23, mm, dd, h, m, s, 0, nil))
	}
	if !x.RecordingTime.IsZero() {
		base, err := xmpbase.MakeModel(d)
		if err != nil {
			return err
		}
		base.CreateDate = x.RecordingTime
	}

	// XmpDM gets always created
	dm, err := xmpdm.MakeModel(d)
	if err != nil {
		return err
	}
	if x.Comments != "" {
		dm.LogComment = x.Comments
	}
	if x.AlbumTitle != "" {
		dm.Album = x.AlbumTitle
	}
	dm.PartOfCompilation = xmp.Bool(x.IsCompilation.Value())
	if x.Composer != "" {
		dm.Composer = x.Composer
	}
	if len(x.ContentType) > 0 {
		dm.Genre = x.ContentType.String()
	}
	if x.LeadPerformer != "" {
		dm.Artist = x.LeadPerformer
	}
	if x.ModifiedBy != "" {
		dm.Engineer = x.ModifiedBy
	}
	if !x.PartOfASet.IsZero() {
		dm.DiscNumber = x.PartOfASet.Track
	}
	if !x.TrackNumber.IsZero() {
		dm.TrackNumber = x.TrackNumber.Track
	}
	if len(x.UnsynchronizedLyrics) > 0 {
		dm.Lyrics = strings.Join(x.UnsynchronizedLyrics, "\n")
	}
	// XmpRights gets copyright
	if x.CopyrightInformation != "" {
		rights, err := xmprights.MakeModel(d)
		if err != nil {
			return err
		}
		rights.WebStatement = string(x.CopyrightInformation)
	}
	return nil
}

type ID3 struct {
	AudioEncryption                      *AudioEncryption        `id3:"AENC"        xmp:"id3:audioEncryption"`                      // AENC Audio encryption
	AttachedPicture                      AttachedPictureArray    `id3:"APIC"        xmp:"id3:attachedPicture"`                      // APIC Attached picture
	AudioSeekPointIndex                  *AudioSeekIndex         `id3:"ASPI,v2.4+"  xmp:"id3:audioSeekPointIndex,v2.4+"`            // ASPI Audio seek point index
	Comments                             string                  `id3:"COMM"        xmp:"id3:comments"`                             // COMM Comments
	Commercial                           *Commercial             `id3:"COMR"        xmp:"id3:commercial"`                           // COMR Commercial frame
	Encryption                           EncryptionMethodArray   `id3:"ENCR"        xmp:"id3:encryption"`                           // ENCR Encryption method registration
	Equalization_v23                     string                  `id3:"EQUA,v2.3-"  xmp:"id3:equalization_v23,v2.3-"`               // EQUA Equalization
	Equalization                         EqualizationList        `id3:"EQU2,v2.4+"  xmp:"id3:equalization,v2.4+"`                   // EQU2 Equalisation (2)
	EventTimingCodes                     MarkerList              `id3:"ETCO"        xmp:"id3:eventTimingCodes"`                     // ETCO Event timing codes
	GeneralEncapsulatedObject            EncapsulatedObjectArray `id3:"GEOB"        xmp:"id3:generalEncapsulatedObject"`            // GEOB General encapsulated object
	GroupIdentifier                      GroupArray              `id3:"GRID"        xmp:"id3:groupIdentifier"`                      // GRID Group identification registration
	InvolvedPeopleList_v23               string                  `id3:"IPLS,v2.3-"  xmp:"id3:involvedPeopleList_v23,v2.3-"`         // IPLS Involved people list
	Link                                 LinkArray               `id3:"LINK"        xmp:"id3:link"`                                 // LINK Linked information
	MusicCDIdentifier                    []byte                  `id3:"MCDI"        xmp:"id3:musicCDIdentifier"`                    // MCDI Music CD identifier
	MPEGLocationLookupTable              []byte                  `id3:"MLLT"        xmp:"id3:mpegLocationLookupTable"`              // MLLT MPEG location lookup table
	Ownership                            *Owner                  `id3:"OWNE"        xmp:"id3:ownership"`                            // OWNE Ownership frame
	Private                              PrivateDataArray        `id3:"PRIV"        xmp:"id3:private"`                              // PRIV Private frame
	PlayCounter                          int64                   `id3:"PCNT"        xmp:"id3:playCounter"`                          // PCNT Play counter
	Popularimeter                        PopularimeterArray      `id3:"POPM"        xmp:"id3:popularimeter"`                        // POPM Popularimeter
	PositionSynchronization              *PositionSync           `id3:"POSS"        xmp:"id3:positionSynchronization"`              // POSS Position synchronisation frame
	RecommendedBufferSize                *BufferSize             `id3:"RBUF"        xmp:"id3:recommendedBufferSize"`                // RBUF Recommended buffer size
	RelativeVolumeAdjustment_v23         string                  `id3:"RVAD,v2.3-"  xmp:"id3:relativeVolumeAdjustment_v23,v2.3-"`   // RVAD Relative volume adjustment
	RelativeVolumeAdjustment             VolumeAdjustArray       `id3:"RVA2,v2.4+"  xmp:"id3:relativeVolumeAdjustment,v2.4+"`       // RVA2 Relative volume adjustment (2)
	Reverb                               *Reverb                 `id3:"RVRB"        xmp:"id3:reverb"`                               // RVRB Reverb
	Seek                                 int64                   `id3:"SEEK,v2.4+"  xmp:"id3:seek,v2.4+"`                           // SEEK Seek frame
	Signature                            SignatureArray          `id3:"SIGN,v2.4+"  xmp:"id3:signature,v2.4+"`                      // SIGN Signature frame
	SynchronizedLyrics                   TimedLyricsArray        `id3:"SYLT"        xmp:"id3:synchronizedLyrics"`                   // SYLT Synchronized lyric/text
	SynchronizedTempoCodes               []byte                  `id3:"SYTC"        xmp:"id3:synchronizedTempoCodes"`               // SYTC Synchronized tempo codes
	AlbumTitle                           string                  `id3:"TALB"        xmp:"id3:albumTitle"`                           // TALB Album/Movie/Show title
	BeatsPerMinute                       int                     `id3:"TBPM"        xmp:"id3:beatsPerMinute"`                       // TBPM BPM (beats per minute)
	Composer                             string                  `id3:"TCOM"        xmp:"id3:composer"`                             // TCOM Composer
	ContentType                          Genre                   `id3:"TCON"        xmp:"id3:contentType"`                          // TCON Content type
	Copyright                            string                  `id3:"TCOP"        xmp:"id3:copyright"`                            // TCOP Copyright message
	Date_v23                             Date23                  `id3:"TDAT,v2.3-"  xmp:"id3:date_v23,v2.3-"`                       // TDAT Date
	EncodingTime                         xmp.Date                `id3:"TDEN,v2.4+"  xmp:"id3:encodingTime,v2.4+"`                   // TDEN Encoding time
	PlaylistDelay                        int64                   `id3:"TDLY"        xmp:"id3:playlistDelay"`                        // TDLY Playlist delay
	OriginalReleaseTime                  xmp.Date                `id3:"TDOR,v2.4+"  xmp:"id3:originalReleaseTime,v2.4+"`            // TDOR Original release time
	RecordingTime                        xmp.Date                `id3:"TDRC,v2.4+"  xmp:"id3:recordingTime,v2.4+"`                  // TDRC Recording time
	ReleaseTime                          xmp.Date                `id3:"TDRL,v2.4+"  xmp:"id3:releaseTime,v2.4+"`                    // TDRL Release time
	TaggingTime                          xmp.Date                `id3:"TDTG,v2.4+"  xmp:"id3:taggingTime,v2.4+"`                    // TDTG Tagging time
	EncodedBy                            string                  `id3:"TENC"        xmp:"id3:encodedBy"`                            // TENC Encoded by
	Lyricist                             string                  `id3:"TEXT"        xmp:"id3:lyricist"`                             // TEXT Lyricist/Text writer
	FileType                             string                  `id3:"TFLT"        xmp:"id3:fileType"`                             // TFLT File type
	Time_v23                             Time23                  `id3:"TIME,v2.3-"  xmp:"id3:time_v23,v2.3-"`                       // TIME Time
	InvolvedPeopleList                   string                  `id3:"TIPL,v2.4+"  xmp:"id3:involvedPeopleList,v2.4+"`             // TIPL Involved people list
	ContentGroupDescription              string                  `id3:"TIT1"        xmp:"id3:contentGroupDescription"`              // TIT1 Content group description
	TitleDescription                     string                  `id3:"TIT2"        xmp:"id3:titleDescription"`                     // TIT2 Title/songname/content description
	SubTitle                             string                  `id3:"TIT3"        xmp:"id3:subTitle"`                             // TIT3 Subtitle/Description refinement
	InitialKey                           string                  `id3:"TKEY"        xmp:"id3:initialKey"`                           // TKEY Initial key
	Language                             string                  `id3:"TLAN"        xmp:"id3:language"`                             // TLAN Language(s)
	Length                               string                  `id3:"TLEN"        xmp:"id3:length"`                               // TLEN Length
	MusicianCreditsList                  string                  `id3:"TMCL,v2.4+"  xmp:"id3:musicianCreditsList,v2.4+"`            // TMCL Musician credits list
	MediaType                            string                  `id3:"TMED"        xmp:"id3:mediaType"`                            // TMED Media type
	Mood                                 string                  `id3:"TMOO,v2.4+"  xmp:"id3:mood,v2.4+"`                           // TMOO Mood
	OriginalAlbumTitle                   string                  `id3:"TOAL"        xmp:"id3:originalAlbumTitle"`                   // TOAL Original album/movie/show title
	OriginalFilename                     string                  `id3:"TOFN"        xmp:"id3:originalFilename"`                     // TOFN Original filename
	OriginalLyricist                     string                  `id3:"TOLY"        xmp:"id3:originalLyricist"`                     // TOLY Original lyricist(s)/text writer(s)
	OriginalArtist                       string                  `id3:"TOPE"        xmp:"id3:originalArtist"`                       // TOPE Original artist(s)/performer(s)
	OriginalReleaseYear_v23              int                     `id3:"TORY,v2.3-"  xmp:"id3:originalReleaseYear_v23,v2.3-"`        // TORY Original release year
	FileOwner                            string                  `id3:"TOWN"        xmp:"id3:fileOwner"`                            // TOWN File owner/licensee
	LeadPerformer                        string                  `id3:"TPE1"        xmp:"id3:leadPerformer"`                        // TPE1 Lead performer(s)/Soloist(s)
	Band                                 string                  `id3:"TPE2"        xmp:"id3:band"`                                 // TPE2 Band/orchestra/accompaniment
	Conductor                            string                  `id3:"TPE3"        xmp:"id3:conductor"`                            // TPE3 Conductor/performer refinement
	ModifiedBy                           string                  `id3:"TPE4"        xmp:"id3:modifiedBy"`                           // TPE4 Interpreted, remixed, or otherwise modified by
	PartOfASet                           TrackNum                `id3:"TPOS"        xmp:"id3:partOfASet"`                           // TPOS Part of a set
	ProducedNotice                       string                  `id3:"TPRO,v2.4+"  xmp:"id3:producedNotice,v2.4+"`                 // TPRO Produced notice
	Publisher                            string                  `id3:"TPUB"        xmp:"id3:publisher"`                            // TPUB Publisher
	TrackNumber                          TrackNum                `id3:"TRCK"        xmp:"id3:trackNumber"`                          // TRCK Track number/Position in set
	RecordingDates                       string                  `id3:"TRDA"        xmp:"id3:recordingDates"`                       // TRDA Recording dates
	InternetRadioStationName             string                  `id3:"TRSN"        xmp:"id3:internetRadioStationName"`             // TRSN Internet radio station name
	InternetRadioStationOwner            string                  `id3:"TRSO"        xmp:"id3:internetRadioStationOwner"`            // TRSO Internet radio station owner
	Size_v23                             string                  `id3:"TSIZ,v2.3-"  xmp:"id3:size_v23,v2.3-"`                       // TSIZ Size
	AlbumSortOrder                       string                  `id3:"TSOA,v2.4+"  xmp:"id3:albumSortOrder,v2.4+"`                 // TSOA Album sort order
	PerformerSortOrder                   string                  `id3:"TSOP,v2.4+"  xmp:"id3:performerSortOrder,v2.4+"`             // TSOP Performer sort order
	TitleSortOrder                       string                  `id3:"TSOT,v2.4+"  xmp:"id3:titleSortOrder,v2.4+"`                 // TSOT Title sort order
	InternationalStandardRecordingCode   string                  `id3:"TSRC"        xmp:"id3:internationalStandardRecordingCode"`   // TSRC ISRC (international standard recording code)
	EncodedWith                          string                  `id3:"TSSE"        xmp:"id3:encodedWith"`                          // TSSE Software/Hardware and settings used for encoding
	SetSubtitle                          string                  `id3:"TSST,v2.4+"  xmp:"id3:setSubtitle,v2.4+"`                    // TSST Set subtitle
	Year_v23                             int                     `id3:"TYER,v2.3-"  xmp:"id3:year_v23,v2.3-"`                       // TYER Year
	UserText                             string                  `id3:"TXXX"        xmp:"id3:userText"`                             // TXXX User defined text information frame
	UniqueFileIdentifier                 xmp.Uri                 `id3:"UFID"        xmp:"id3:uniqueFileIdentifier"`                 // UFID Unique file identifier
	TermsOfUse                           xmp.StringList          `id3:"USER"        xmp:"id3:termsOfUse"`                           // USER Terms of use
	UnsynchronizedLyrics                 xmp.StringList          `id3:"USLT"        xmp:"id3:unsynchronizedLyrics"`                 // USLT Unsynchronized lyric/text transcription
	CommercialInformation                xmp.Url                 `id3:"WCOM"        xmp:"id3:commercialInformation"`                // WCOM Commercial information
	CopyrightInformation                 xmp.Url                 `id3:"WCOP"        xmp:"id3:copyrightInformation"`                 // WCOP Copyright/Legal information
	OfficialAudioFileWebpage             xmp.Url                 `id3:"WOAF"        xmp:"id3:officialAudioFileWebpage"`             // WOAF Official audio file webpage
	OfficialArtistWebpage                xmp.Url                 `id3:"WOAR"        xmp:"id3:officialArtistWebpage"`                // WOAR Official artist/performer webpage
	OfficialAudioSourceWebpage           xmp.Url                 `id3:"WOAS"        xmp:"id3:officialAudioSourceWebpage"`           // WOAS Official audio source webpage
	OfficialInternetRadioStationHomepage xmp.Url                 `id3:"WORS"        xmp:"id3:officialInternetRadioStationHomepage"` // WORS Official Internet radio station homepage
	Payment                              xmp.Url                 `id3:"WPAY"        xmp:"id3:payment"`                              // WPAY Payment
	OfficialPublisherWebpage             xmp.Url                 `id3:"WPUB"        xmp:"id3:officialPublisherWebpage"`             // WPUB Publishers official webpage
	UserURL                              xmp.Url                 `id3:"WXXX"        xmp:"id3:userURL"`                              // WXXX User defined URL link frame
	AccessibilityText                    string                  `id3:"ATXT"        xmp:"id3:accessibilityText"`                    // accessibility addendum http://id3.org/id3v2-accessibility-1.0
	Chapters                             ChapterList             `id3:"CHAP"        xmp:"id3:chapters"`                             // chapters addedum http://id3.org/id3v2-chapters-1.0
	TableOfContents                      TocEntryList            `id3:"CTOC"        xmp:"id3:tableOfContents"`                      // chapters addedum http://id3.org/id3v2-chapters-1.0
	AlbumArtistSortOrder                 string                  `id3:"TSO2"        xmp:"id3:albumArtistSortOrder"`                 // special iTunes tags
	ComposerSortOrder                    string                  `id3:"TSOC"        xmp:"id3:composerSortOrder"`
	IsCompilation                        Bool                    `id3:"TCMP"        xmp:"id3:isCompilation"`
	ITunesU                              Bool                    `id3:"ITNU"        xmp:"id3:iTunesU"`
	IsPodcast                            Bool                    `id3:"PCST"        xmp:"id3:isPodcast"`
	PodcastDescription                   string                  `id3:"TDES"        xmp:"id3:podcastDescription"`
	PodcastID                            string                  `id3:"TGID"        xmp:"id3:podcastID"`
	PodcastURL                           string                  `id3:"WFED"        xmp:"id3:podcastURL"`
	PodcastKeywords                      string                  `id3:"TKWD"        xmp:"id3:podcastKeywords"`
	PodcastCategory                      string                  `id3:"TCAT"        xmp:"id3:podcastCategory"`
	Extension                            xmp.TagList             `id3:",any"        xmp:"id3:extension"`
}

var V11_TO_V24_TAGS map[string]string = map[string]string{
	"0x0003": "TIT2", // title
	"0x0021": "TPE1", // artist
	"0x003f": "TALB", // album
	"0x005d": "TYER", // year
	"0x0061": "COMM", // comment
	"0x007d": "TRCK", // track
	"0x007f": "TCON", // genre
}

var V22_TO_V24_TAGS map[string]string = map[string]string{
	"BUF": "RBUF", //   Recommended buffer size
	"CNT": "PCNT", //   Play counter
	"COM": "COMM", //   Comments
	"CRA": "AENC", //   Audio encryption
	"CRM": "ENCR", //   Encrypted meta frame
	"ETC": "ETCO", //   Event timing codes
	"EQU": "EQU2", //   Equalization
	"GEO": "GEOB", //   General encapsulated object
	"IPL": "TIPL", //   Involved people list
	"LNK": "LINK", //   Linked information
	"MCI": "MCDI", //   Music CD Identifier
	"MLL": "MLLT", //   MPEG location lookup table
	"PIC": "APIC", //   Attached picture
	"POP": "POPM", //   Popularimeter
	"REV": "RVRB", //   Reverb
	"RVA": "RVA2", //   Relative volume adjustment
	"SLT": "SYLT", //   Synchronized lyric/text
	"STC": "SYTC", //   Synced tempo codes
	"TAL": "TALB", //   Album/Movie/Show title
	"TBP": "TBPM", //   BPM (Beats Per Minute)
	"TCM": "TCOM", //   Composer
	"TCO": "TMED", //   Content type
	"TCR": "TCOP", //   Copyright message
	"TDA": "TDRC", // # Date, arguably this could also be release date TDRL
	"TDY": "TDLY", //   Playlist delay
	"TEN": "TENC", //   Encoded by
	"TFT": "TFLT", //   File type
	"TIM": "TDRC", // # Time, arguably this could also be release date TDRL
	"TKE": "TKEY", //   Initial key
	"TLA": "TLAN", //   Language(s)
	"TLE": "TLEN", //   Length
	"TMT": "TMED", //   Media type
	"TOA": "TOPE", //   Original artist(s)/performer(s)
	"TOF": "TOFN", //   Original filename
	"TOL": "TOLY", //   Original Lyricist(s)/text writer(s)
	"TOR": "TORY", //   Original release year (v2.3-)
	"TOT": "TOAL", //   Original album/Movie/Show title
	"TP1": "TPE1", //   Lead artist(s)/Lead performer(s)/Soloist(s)/Performing group
	"TP2": "TPE2", //   Band/Orchestra/Accompaniment
	"TP3": "TPE3", //   Conductor/Performer refinement
	"TP4": "TPE4", //   Interpreted, remixed, or otherwise modified by
	"TPA": "TPOS", //   Part of a set
	"TPB": "TPUB", //   Publisher
	"TRC": "TSRC", //   ISRC (International Standard Recording Code)
	"TRD": "TDRC", //   Recording dates
	"TRK": "TRCK", //   Track number/Position in set
	"TSI": "TSIZ", //   Size (v2.3-)
	"TSS": "TSSE", //   Software/hardware and settings used for encoding
	"TT1": "TIT1", //   Content group description
	"TT2": "TIT2", //   Title/Songname/Content description
	"TT3": "TIT3", //   Subtitle/Description refinement
	"TXT": "TEXT", //   Lyricist/text writer
	"TXX": "TXXX", //   User defined text information frame
	"TYE": "TYER", //   Year (v2.3-)
	"UFI": "UFID", //   Unique file identifier
	"ULT": "USLT", //   Unsychronized lyric/text transcription
	"WAF": "WOAF", //   Official audio file webpage
	"WAR": "WOAR", //   Official artist/performer webpage
	"WAS": "WOAS", //   Official audio source webpage
	"WCM": "WCOM", //   Commercial information
	"WCP": "WCOP", //   Copyright/Legal information
	"WPB": "WPUB", //   Publishers official webpage
	"WXX": "WXXX", //   User defined URL link frame
	"TCP": "TCMP", //   iTunes compilation
	"TST": "TSOT", //   iTunes title sort order
	"TS2": "TSO2", //   iTunes album_artist_sort
	"TSA": "TSOA", //   iTunes album_sort
	"TSP": "TSOP", //   iTunes artist_sort
	"TSC": "TSOC", //   iTunes composer_sort
}

// tag can be a four letter v2.3/v2.4 code or a 3 letter v2.2 code
func mapTagToV24(tag string) (string, error) {
	if len(tag) > 4 {
		if t, ok := V11_TO_V24_TAGS[tag]; !ok {
			return tag, fmt.Errorf("id3: unknown v1 tag '%s'", tag)
		} else {
			tag = t
		}
	} else if len(tag) < 4 {
		if t, ok := V22_TO_V24_TAGS[tag]; !ok {
			return tag, fmt.Errorf("id3: unknown v2.2 tag '%s'", tag)
		} else {
			tag = t
		}
	}
	return tag, nil
}

func (x *ID3) CanTag(tag string) bool {
	t, err := mapTagToV24(tag)
	if err != nil {
		return false
	}
	_, err = xmp.GetNativeField(x, t)
	return err == nil
}

func (x *ID3) GetTag(tag string) (string, error) {
	t, err := mapTagToV24(tag)
	if err != nil {
		return "", err
	}
	if v, err := xmp.GetNativeField(x, t); err != nil {
		return "", fmt.Errorf("%s: %v", NsID3.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *ID3) SetTag(tag, value string) error {
	t, err := mapTagToV24(tag)
	if err != nil {
		return err
	}
	if err := xmp.SetNativeField(x, t, value); err != nil {
		return fmt.Errorf("%s: %v", NsID3.GetName(), err)
	}
	return nil
}

// Returns tag value as string for the selected language. Falls back to
// default language when lang is empty.
func (x *ID3) GetLocaleTag(lang string, tag string) (string, error) {
	t, err := mapTagToV24(tag)
	if err != nil {
		return "", err
	}
	if val, err := xmp.GetLocaleField(x, lang, t); err != nil {
		return "", fmt.Errorf("%s: %v", NsID3.GetName(), err)
	} else {
		return val, nil
	}
}

func (x *ID3) SetLocaleTag(lang string, tag, value string) error {
	if err := xmp.SetLocaleField(x, lang, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsID3.GetName(), err)
	}
	return nil
}

// Lists all non-empty tags.
func (x *ID3) ListTags() (xmp.TagList, error) {
	if l, err := xmp.ListNativeFields(x); err != nil {
		return nil, fmt.Errorf("%s: %v", NsID3.GetName(), err)
	} else {
		return l, nil
	}
}
