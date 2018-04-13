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

// Types as defined in ISO 16684-1:2011(E) 8.2.1 (Core value types)
// - Bool
// - Date
// - AgentName
// - GUID
// - Uri
// - Url
// - Rational
// - Rating (enum)
// Derived Types as defined in ISO 16684-1:2011(E) 8.2.2 (Derived value types)
// - Locale
// - GPSCoord
// Array Types
// - DateList (ordered)
// - GUIDList (ordered)
// - UriList (ordered)
// - UriArray (unordered)
// - UrlList (ordered)
// - UrlArray (unordered)
// - RationalArray (unordered)
// - LocaleArray

package xmp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Zero interface {
	IsZero() bool
}

// 8.2.1.1 Boolean
type Bool bool

const (
	True  Bool = true
	False Bool = false
)

func (x Bool) Value() bool {
	return bool(x)
}

func (x Bool) MarshalText() ([]byte, error) {
	if x {
		return []byte("True"), nil
	}
	return []byte("False"), nil
}

func (x *Bool) UnmarshalText(data []byte) error {
	s := string(data)
	switch s {
	case "True", "true", "TRUE":
		*x = true
	case "False", "false", "FALSE":
		*x = false
	default:
		return fmt.Errorf("xmp: invalid bool value '%s'", s)
	}
	return nil
}

func (x Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Value())
}

// 8.2.1.2 Date
type Date time.Time

func NewDate(t time.Time) Date {
	return Date(t)
}

func Now() Date {
	return Date(time.Now())
}

func (x Date) Time() Time {
	return Time(x.Value())
}

func (x Date) Value() time.Time {
	return time.Time(x)
}

func (x Date) String() string {
	return time.Time(x).Format(time.RFC3339)
}

func (x Date) IsZero() bool {
	return time.Time(x).IsZero()
}

func (x Date) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(x.String()), nil
}

var dateFormats []string = []string{
	"2006-01-02T15:04:05.999999999",       // XMP
	time.RFC3339,                          // XMP
	"2006-01-02T15:04-07:00",              // EXIF
	"2006-01-02",                          // EXIF
	"2006-01-02 15:04:05",                 // EXR
	"2006-01-02T15:04:05.999999999Z07:00", // XMP
	"2006-01-02T15:04:05.999999999Z",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-0700",
	"2006:01:02 15:04:05.999",   // MXF
	"2006/01/02T15:04:05-07:00", // Arri CSV
	"06/01/02T15:04:05-07:00",   // Arri CSV
	"20060102T15h04m05-07:00",   // Arri QT
	"20060102T15h04m05s-07:00",  // Arri XML in MXF
	"2006-01-02T15:04:05",
	"2006-01-02T15:04Z",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
	"2006:01:02", // ID3 date
	"2006-01",
	"2006",
	"15:04:05-07:00",                // time with timezone (IPTC)
	"15:04:05",                      // time without timezone (IPTC)
	"150405-0700",                   // time with timezone (Getty)
	"2006-01-02T00:00:00.000000000", // zero filler to catch potential bad date strings
	"2006-01-00T00:00:00.000000000", // zero filler to catch potential bad date strings
	"2006-00-00T00:00:00.000000000", // zero filler to catch potential bad date strings
	"2006-01-02T00:00:00Z",          // zero filler to catch potential bad date strings
	"2006-01-00T00:00:00Z",          // zero filler to catch potential bad date strings
	"2006-00-00T00:00:00Z",          // zero filler to catch potential bad date strings
}

var illegalZero StringList = StringList{
	"--", // ARRI undefined
	"00/00/00T00:00:00+00:00", // ARRI zero time
}

// repair single-digit hours in timezones
// "00/00/00T00:00:00+00:00", // ARRI zero time
// "2011-02-15T10:15:14+1:00"
// "2011-02-15T10:15:14+1"
// "2017-09-15T20:17:41+00200", // seen from iPhone5s iOS10 video, ffprobe
func repairTZ(value string) string {
	if illegalZero.Contains(value) {
		return "0001-01-01T00:00:00Z"
	}
	l := len(value)
	if l < 20 {
		return value
	}
	a, b, c := value[l-5], value[l-2], value[l-6]
	switch a {
	case '+', '-', 'Z':
		return value[:l-4] + "0" + value[l-4:]
	}
	switch b {
	case '+', '-', 'Z':
		return value[:l-1] + "0" + value[l-1:] + ":00"
	}
	switch c {
	case '+', '-', 'Z':
		if !strings.Contains(value[l-6:], ":") {
			return value[:l-5] + value[l-4:]
		}
	}
	return value
}

func ParseDate(value string) (Date, error) {
	if value != "" {
		value = repairTZ(value)
		for _, f := range dateFormats {
			if t, err := time.Parse(f, value); err == nil {
				return Date(t), nil
			}
		}
	}
	return Date{}, fmt.Errorf("xmp: invalid datetime value '%s'", value)
}

func (x *Date) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*x = Date{}
		return nil
	}
	if d, err := ParseDate(string(data)); err != nil {
		return err
	} else {
		*x = d
	}
	return nil
}

type DateList []Date

func (x DateList) Typ() ArrayType {
	return ArrayTypeOrdered
}

func (x DateList) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *DateList) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	// be resilient against broken writers (PDF)
	if len(node.Nodes) == 0 && len(node.Value) > 0 {
		if d, err := ParseDate(node.Value); err != nil {
			return err
		} else {
			*x = append(*x, d)
		}
		return nil
	}
	return UnmarshalArray(d, node, x.Typ(), x)
}

// Non-Standard Time
type Time time.Time

func NewTime(t time.Time) Time {
	return Time(t)
}

func (x Time) Value() time.Time {
	return time.Time(x)
}

func (x Time) String() string {
	return time.Time(x).Format("15:04:05-07:00")
}

func (x Time) IsZero() bool {
	return time.Time(x).IsZero()
}

func (x Time) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(x.String()), nil
}

var timeFormats []string = []string{
	"15:04:05-07:00", // time with timezone (IPTC)
	"15:04:05",       // time without timezone (IPTC)
	"150405-0700",    // time with timezone (Getty)
}

func ParseTime(value string) (Time, error) {
	if value != "" {
		for _, f := range timeFormats {
			if t, err := time.Parse(f, value); err == nil {
				return Time(t), nil
			}
		}
	}
	return Time{}, fmt.Errorf("xmp: invalid time value '%s'", value)
}

func (x *Time) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if d, err := ParseTime(string(data)); err != nil {
		return err
	} else {
		*x = d
	}
	return nil
}

// 8.2.2.1 AgentName
//
type AgentName string

func (x AgentName) IsZero() bool {
	return len(x) == 0
}

func (x AgentName) String() string {
	return string(x)
}

// 8.2.2.3 GUID (a simple non-Uri identifier)
//
type GUID string

// func newGUID(issuer string, u uuid.UUID) GUID {
// 	return GUID(strings.Join([]string{issuer, u.String()}, ":"))
// }

// func NewGUID() GUID {
// 	return newGUID("uuid", uuid.NewV4())
// }

// func NewGUIDFrom(issuer string, u uuid.UUID) GUID {
// 	if issuer == "" {
// 		issuer = "uuid"
// 	}
// 	return newGUID(issuer, u)
// }

func (x GUID) IsZero() bool {
	return x == ""
}

// func (x GUID) UUID() uuid.UUID {
// 	if idx := strings.LastIndex(string(x), ":"); idx > -1 {
// 		return uuid.FromStringOrNil(string(x[idx+1:]))
// 	}
// 	return uuid.FromStringOrNil(string(x))
// }

func (x GUID) Issuer() string {
	if idx := strings.LastIndex(string(x), ":"); idx > -1 {
		return string(x[:idx])
	}
	return ""
}

func (x GUID) Value() string {
	return string(x)
}

func (x GUID) String() string {
	return string(x)
}

func (x GUID) MarshalText() ([]byte, error) {
	return []byte(x), nil
}

func (x *GUID) UnmarshalText(data []byte) error {
	*x = GUID(data)
	return nil
}

func (x GUID) MarshalXMP(e *Encoder, node *Node, m Model) error {
	if x.IsZero() {
		return nil
	}
	b, _ := x.MarshalText()
	return e.EncodeElement(b, node)
}

type GUIDList []GUID

func (x GUIDList) Typ() ArrayType {
	return ArrayTypeOrdered
}

func (x GUIDList) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *GUIDList) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

type Url string

func (x Url) Value() string {
	return string(x)
}

func (x Url) IsZero() bool {
	return x == ""
}

type UrlArray []Url

func (x UrlArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func (x UrlArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *UrlArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

type UrlList []Url

func (x UrlList) Typ() ArrayType {
	return ArrayTypeOrdered
}

func (x UrlList) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *UrlList) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

// 8.2.2.10 Uri (Note: special handling as rdf:resource attribute without node content)
//
type Uri string

func (x Uri) Value() string {
	return string(x)
}

func (x Uri) IsZero() bool {
	return x == ""
}

func NewUri(u string) Uri {
	return Uri(u)
}

// - supported
// <xmp:BaseUrl rdf:resource="http://www.adobe.com/"/>
//
// - supported
// <xmp:BaseUrl rdf:parseType="Resource">
//   <rdf:value rdf:resource="http://www.adobe.com/"/>
//   <xe:qualifier>artificial example</xe:qualifier>
// </xmp:BaseUrl>
//
// func (x Uri) MarshalXMP(e *Encoder, node *Node, m Model) error {
// 	node.Attr = append(node.Attr, xml.Attr{
// 		Name:  xml.Name{Local: "rdf:resource"},
// 		Value: string(x),
// 	})
// 	return nil
// }

func (x *Uri) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	if attr := node.GetAttr("rdf", "resource"); len(attr) > 0 {
		*x = Uri(attr[0].Value)
		return nil
	}
	var u string
	if err := d.DecodeElement(&u, node); err != nil {
		return err
	}
	*x = Uri(u)
	return nil
}

type UriArray []Uri

func (x UriArray) IsZero() bool {
	return len(x) == 0
}

func (x UriArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func (x UriArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *UriArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

type UriList []Uri

func (x UriList) IsZero() bool {
	return len(x) == 0
}

func (x UriList) Typ() ArrayType {
	return ArrayTypeOrdered
}

func (x UriList) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *UriList) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

// Rational "n/m"
//
type Rational struct {
	Num int64
	Den int64
}

func (x *Rational) Addr() *Rational {
	if x.IsZero() {
		return nil
	}
	return x
}

func (x Rational) IsZero() bool {
	return x.Den == 0 && x.Num == 0
}

func (x Rational) Value() float64 {
	if x.Den == 0 {
		return 1
	}
	return float64(x.Num) / float64(x.Den)
}

func (x Rational) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(strconv.FormatInt(x.Num, 10))
	buf.WriteByte('/')
	buf.WriteString(strconv.FormatInt(x.Den, 10))
	return buf.String()
}

func (x Rational) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func gcd(x, y int64) int64 {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// Beware: primitive conversion algorithm
func FloatToRational(f float32) Rational {
	var (
		den int64   = 1000000
		rnd float32 = 0.5
	)
	switch {
	case f > 2147:
		den = 10000
	case f > 214748:
		den = 100
	case f > 21474836:
		den = 1
	}
	if f < 0 {
		rnd = -0.5
	}
	nom := int64(f*float32(den) + rnd)
	g := gcd(nom, den)
	return Rational{nom / g, den / g}
}

func (x *Rational) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var err error
	v := string(data)
	r := Rational{}
	switch {
	case strings.Contains(v, "."):
		var f float64
		f, err = strconv.ParseFloat(v, 32)
		if err == nil {
			if f < 1 {
				r = Rational{1, int64(1000/f+50) / 1000}
			} else {
				r = FloatToRational(float32(f))
			}
			// fmt.Printf("Float %f in Rational %d/%d\n", f, r.Num, r.Den)
		}
	case strings.Contains(v, "/"):
		_, err = fmt.Sscanf(v, "%d/%d", &r.Num, &r.Den)
	case strings.Contains(v, " "):
		_, err = fmt.Sscanf(v, "%d %d", &r.Num, &r.Den)
	default:
		r.Num, err = strconv.ParseInt(v, 10, 64)
		r.Den = 1
	}
	if err != nil {
		return fmt.Errorf("xmp: invalid rational '%s': %v", v, err)
	}
	*x = r
	return nil
}

type RationalArray []Rational

func (a RationalArray) Typ() ArrayType {
	return ArrayTypeOrdered
}

func NewRationalArray(items ...Rational) RationalArray {
	a := make(RationalArray, 0, len(items))
	return append(a, items...)
}

func (a RationalArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	if len(a) == 0 {
		return nil
	}
	return MarshalArray(e, node, a.Typ(), a)
}

func (a *RationalArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, a.Typ(), a)
}

// GPS Coordinates
//
// “DDD,MM,SSk” or “DDD,MM.mmk”
//  DDD is a number of degrees
//  MM is a number of minutes
//  SS is a number of seconds
//  mm is a fraction of minutes
//  k is a single character N, S, E, or W indicating a direction (north, south, east, west)
type GPSCoord string

func (x GPSCoord) Value() string {
	return string(x)
}

func (x GPSCoord) IsZero() bool {
	return x == "" || x == "0,0.0000000"
}
