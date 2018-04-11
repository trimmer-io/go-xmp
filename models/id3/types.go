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

// http://id3.org/
// http://id3.org/id3v2.4.0-frames

// FIXME: Version 2.4 of the specification prescribes that all text fields
// (the fields that start with a T, except for TXXX) can contain multiple
// values separated by a null character. The null character varies by
// character encoding.

package id3

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"trimmer.io/go-xmp/models/xmp_dm"
	"trimmer.io/go-xmp/xmp"
)

const (
	TIME_FORMAT_23 = "1504" // HHMM
	DATE_FORMAT_23 = "0102" // DDMM
)

type Bool int // 0 no, 1 yes

const (
	False Bool = 0
	True  Bool = 1
)

func (b Bool) Value() bool {
	return b == 1
}

type Time23 time.Time

func (x Time23) Value() time.Time {
	return time.Time(x)
}

func (x Time23) IsZero() bool {
	return time.Time(x).IsZero()
}

func (x *Time23) UnmarshalText(data []byte) error {
	value := string(data)
	if t, err := time.Parse(TIME_FORMAT_23, value); err != nil {
		return fmt.Errorf("id3: invalid time value '%s'", value)
	} else {
		*x = Time23(t)
	}
	return nil
}

func (x Time23) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(time.Time(x).Format(TIME_FORMAT_23)), nil
}

type Date23 time.Time

func (x Date23) Value() time.Time {
	return time.Time(x)
}

func (x Date23) IsZero() bool {
	return time.Time(x).IsZero()
}

func (x *Date23) UnmarshalText(data []byte) error {
	value := string(data)
	if t, err := time.Parse(DATE_FORMAT_23, value); err != nil {
		return fmt.Errorf("id3: invalid date value '%s'", value)
	} else {
		*x = Date23(t)
	}
	return nil
}

func (x Date23) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(time.Time(x).Format(DATE_FORMAT_23)), nil
}

// TCON
//
type GenreV1 byte

func (x *GenreV1) UnmarshalText(data []byte) error {
	value := string(data)
	g, err := ParseGenreV1(value)
	if err != nil {
		return err
	}
	*x = g
	return nil
}

func (x GenreV1) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// See Genre List at http://id3.org/id3v2.3.0
func (x GenreV1) String() string {
	if v, ok := GenreMap[x]; ok {
		return v
	}
	return strconv.Itoa(int(x))
}

func ParseGenreV1(s string) (GenreV1, error) {
	for i, v := range GenreMap {
		if v == s {
			return i, nil
		}
	}
	if t, err := strconv.Atoi(s); err == nil {
		return GenreV1(t), nil
	}
	return 0xff, fmt.Errorf("id3: invalid genre '%s'", s)
}

type Genre []string

func (x *Genre) UnmarshalText(data []byte) error {
	value := string(data)
	for _, v := range strings.Split(value, ",") {
		v = strings.TrimSpace(v)
		if g, err := ParseGenreV1(v); err == nil {
			*x = append(*x, g.String())
			continue
		}
		switch value {
		case "RX":
			*x = append(*x, "Remix")
		case "CR":
			*x = append(*x, "Cover")
		default:
			*x = append(*x, value)
		}
	}
	return nil
}

func (x Genre) String() string {
	return strings.Join(x, ", ")
}

func (x Genre) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// TRCK
// 10
// 10/12
type TrackNum struct {
	Track int
	Total int
}

func (x TrackNum) Value() int {
	return x.Track
}

func (x TrackNum) IsZero() bool {
	return x.Track == 0 && x.Total == 0
}

func (x *TrackNum) UnmarshalText(data []byte) error {
	var err error
	v := string(data)
	n := TrackNum{}
	if strings.Contains(v, "/") {
		_, err = fmt.Sscanf(v, "%d/%d", &n.Track, &n.Total)
	} else {
		n.Track, err = strconv.Atoi(v)
	}
	if err != nil {
		return fmt.Errorf("id3: invalid track number '%s': %v", v, err)
	}
	*x = n
	return nil
}

func (x TrackNum) MarshalText() ([]byte, error) {
	if x.Total == 0 {
		return []byte(strconv.Itoa(x.Track)), nil
	}
	buf := bytes.Buffer{}
	buf.WriteString(strconv.FormatInt(int64(x.Track), 10))
	buf.WriteByte('/')
	buf.WriteString(strconv.FormatInt(int64(x.Total), 10))
	return buf.Bytes(), nil
}

// OWNE
type Owner struct {
	Currency     string   `xmp:"id3:currency"`
	Price        float32  `xmp:"id3:price"`
	PurchaseDate xmp.Date `xmp:"id3:purchaseDate"`
	Seller       string   `xmp:"id3:seller"`
}

func (x *Owner) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

// ETCO
type Marker struct {
	MarkerType MarkerType `xmp:"id3:type,attr"`
	Timestamp  int64      `xmp:"id3:timestamp,attr"`
	Unit       UnitType   `xmp:"id3:unit,attr"`
}

type MarkerList []Marker

func (x *MarkerList) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x MarkerList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x MarkerList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *MarkerList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// SYLT
//
type TimedLyrics struct {
	Lang   string        `xmp:"id3:lang,attr"`
	Type   LyricsType    `xmp:"id3:type,attr"`
	Unit   PositionType  `xmp:"id3:unit,attr"`
	Lyrics TimedTextList `xmp:"id3:lyrics,attr"`
}

type TimedLyricsArray []TimedLyrics

func (x TimedLyricsArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x TimedLyricsArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *TimedLyricsArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// see id3 v.24 spec 4.9
// - one id3 frame per language and content type
// - null-terminated text followed by a timestamp and optionally newline
//
//    "Strang" $00 xx xx "ers" $00 xx xx " in" $00 xx xx " the" $00 xx xx
//    " night" $00 xx xx 0A "Ex" $00 xx xx "chang" $00 xx xx "ing" $00 xx
//    xx "glan" $00 xx xx "ces" $00 xx xx
//
func (x *TimedLyricsArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

type TimedText struct {
	Text      string `xmp:"id3:text,attr"`
	Timestamp int64  `xmp:"id3:timestamp,attr"`
}

type TimedTextList []TimedText

func (x TimedTextList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x TimedTextList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *TimedTextList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// COMM
//
type Comment struct {
	Lang string `xmp:"id3:lang,attr"`
	Type string `xmp:"id3:type,attr"`
	Text string `xmp:"id3:text,attr"`
}

type CommentArray []Comment

// TODO: need samples
func (x *CommentArray) UnmarshalText(data []byte) error {
	*x = append(*x, Comment{
		Text: string(data),
	})
	return nil
}

func (x CommentArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x CommentArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *CommentArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// RVA2
//
type VolumeAdjust struct {
	Channel    ChannelType `xmp:"id3:channel,attr"`
	Adjustment float32     `xmp:"id3:adjust,attr"`
	Bits       byte        `xmp:"id3:bits,attr"`
	PeakVolume float32     `xmp:"id3:peak,attr"`
}

type VolumeAdjustArray []VolumeAdjust

func (x *VolumeAdjustArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x VolumeAdjustArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x VolumeAdjustArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *VolumeAdjustArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// EQU2
//
type Equalization struct {
	Method         EqualizationMethod  `xmp:"id3:method,attr"`
	Identification string              `xmp:"id3:id,attr"`
	Adjustments    AdjustmentPointList `xmp:"id3:points"`
}

type EqualizationList []Equalization

func (x *EqualizationList) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x EqualizationList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x EqualizationList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *EqualizationList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type AdjustmentPoint struct {
	Frequency  int     `xmp:"id3:frequency,attr"`
	Adjustment float32 `xmp:"id3:adjust,attr"`
}

type AdjustmentPointList []AdjustmentPoint

func (x AdjustmentPointList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x AdjustmentPointList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *AdjustmentPointList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// RVRB
type Reverb struct {
	ReverbLeftMs       int  `xmp:"id3:reverbLeftMs"`
	ReverbRightMs      int  `xmp:"id3:reverbRightMs"`
	ReverbBouncesLeft  byte `xmp:"id3:reverbBouncesLeft"`
	ReverbBouncesRight byte `xmp:"id3:reverbBouncesRight"`
	ReverbFeedbackLtoL byte `xmp:"id3:reverbFeedbackLtoL"`
	ReverbFeedbackLToR byte `xmp:"id3:reverbFeedbackLToR"`
	ReverbFeedbackRtoR byte `xmp:"id3:reverbFeedbackRtoR"`
	ReverbFeedbackRtoL byte `xmp:"id3:reverbFeedbackRtoL"`
	PremixLtoR         byte `xmp:"id3:premixLtoR"`
	PremixRtoL         byte `xmp:"id3:premixRtoL"`
}

func (x *Reverb) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

// APIC
//
type AttachedPicture struct {
	Mimetype    string      `xmp:"id3:mimetype,attr"`
	Type        PictureType `xmp:"id3:pictureType,attr"`
	Description string      `xmp:"id3:description,attr"`
	Data        []byte      `xmp:"id3:image"`
}

type AttachedPictureArray []AttachedPicture

func (x *AttachedPictureArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x AttachedPictureArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x AttachedPictureArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *AttachedPictureArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// GEOB
//
type EncapsulatedObject struct {
	Mimetype    string `xmp:"id3:mimetype,attr"`
	Filename    string `xmp:"id3:filename,attr"`
	Description string `xmp:"id3:description,attr"`
	Data        []byte `xmp:"id3:data"`
}

type EncapsulatedObjectArray []EncapsulatedObject

func (x *EncapsulatedObjectArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x EncapsulatedObjectArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x EncapsulatedObjectArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *EncapsulatedObjectArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// POPM
//
type Popularimeter struct {
	Email   string `xmp:"id3:email,attr"`
	Rating  byte   `xmp:"id3:rating,attr"`
	Counter int64  `xmp:"id3:counter,attr"`
}

type PopularimeterArray []Popularimeter

func (x *PopularimeterArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x PopularimeterArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x PopularimeterArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *PopularimeterArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// RBUF
//
type BufferSize struct {
	BufferSize int64 `xmp:"id3:bufferSize,attr"`
	Flag       byte  `xmp:"id3:flag,attr"`
	NextOffset int64 `xmp:"id3:next,attr"`
}

func (x *BufferSize) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

// AENC
//
type AudioEncryption struct {
	Owner          string `xmp:"id3:owner,attr"`
	PreviewStart   int    `xmp:"id3:previewStart,attr"`
	PreviewLength  int    `xmp:"id3:previewLength,attr"`
	EncryptionInfo []byte `xmp:"id3:encryptionInfo"`
}

func (x *AudioEncryption) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

// LINK
//
type Link struct {
	LinkedID string         `xmp:"id3:linkedId,attr"`
	Url      string         `xmp:"id3:url,attr"`
	Extra    xmp.StringList `xmp:"id3:extra"`
}

type LinkArray []Link

func (x *LinkArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x LinkArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x LinkArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *LinkArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// POSS
//
type PositionSync struct {
	Unit     PositionType `xmp:"id3:unit,attr"`
	Position int64        `xmp:"id3:position,attr"`
}

// COMR
//
type Commercial struct {
	Prices       PriceList      `xmp:"id3:prices"`
	ValidUntil   Date23         `xmp:"id3:validUntil,attr"`
	ContactUrl   xmp.Url        `xmp:"id3:contactUrl"`
	ReceivedAs   DeliveryMethod `xmp:"id3:receivedAs,attr"`
	SellerName   string         `xmp:"id3:sellerName,attr"`
	Description  string         `xmp:"id3:description,attr"`
	LogoMimetype string         `xmp:"id3:logoMimetype,attr"`
	LogoImage    []byte         `xmp:"id3:logoImage"`
}

func (x *Commercial) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

type Price struct {
	Currency string
	Amount   float32
}

type PriceList []Price

func (x PriceList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x PriceList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *PriceList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// ENCR
//
type EncryptionMethod struct {
	OwnerUrl xmp.Url `xmp:"id3:ownerUrl,attr"`
	Method   byte    `xmp:"id3:method,attr"`
	Data     []byte  `xmp:"id3:data"`
}

type EncryptionMethodArray []EncryptionMethod

func (x *EncryptionMethodArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x EncryptionMethodArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x EncryptionMethodArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *EncryptionMethodArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// GRID
//
type Group struct {
	OwnerUrl xmp.Url `xmp:"id3:ownerUrl,attr"`
	Symbol   byte    `xmp:"id3:symbol,attr"` // 0x80-0xF0
	Data     []byte  `xmp:"id3:data"`
}

type GroupArray []Group

func (x *GroupArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x GroupArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x GroupArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *GroupArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// PRIV
//
type PrivateData struct {
	Owner string `xmp:"id3:owner,attr"`
	Data  []byte `xmp:"id3:data"`
}

type PrivateDataArray []PrivateData

func (x *PrivateDataArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x PrivateDataArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x PrivateDataArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *PrivateDataArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// SIGN
//
type Signature struct {
	GroupSymbol byte   `xmp:"id3:groupSymbol,attr"` // 0x80-0xF0
	Signature   []byte `xmp:"id3:signature"`
}

type SignatureArray []Signature

func (x *SignatureArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x SignatureArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x SignatureArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *SignatureArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// ASPI
//
type AudioSeekIndex struct {
	Start          int64       `xmp:"id3:start,attr"`
	Length         int64       `xmp:"id3:length,attr"`
	NumberOfPoints int         `xmp:"id3:numberOfPoints,attr"`
	BitsPerPoint   byte        `xmp:"id3:bitsPerPoint,attr"`
	Points         xmp.IntList `xmp:"id3:points"`
}

func (x *AudioSeekIndex) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

// CHAP
//
type Chapter xmpdm.Marker

type ChapterList []Chapter

func (x *ChapterList) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x ChapterList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x ChapterList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *ChapterList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// CTOC
//
type TocEntry struct {
	Flags  int            `xmp:"id3:flags"`
	IDList xmp.StringList `xmp:"id3:idlist"`
	Info   string         `xmp:"id3:info"`
}

type TocEntryList []TocEntry

func (x *TocEntryList) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x TocEntryList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x TocEntryList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *TocEntryList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}
