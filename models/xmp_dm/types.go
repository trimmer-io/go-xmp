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

// Package xmpdm implements the XMP Dynamic Media namespace as defined by XMP Specification Part 2.
package xmpdm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"trimmer.io/go-xmp/xmp"
)

// 1.2.3.1 Part
// [/]time:##     - duration from node=0
// [/]time:##d##  - duration from given node
// [/]time:##r##  - range from given node to end
//
// ## is FrameCount
type Part struct {
	Path     string
	Start    FrameCount
	Duration FrameCount
}

func prefixIfNot(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s
	}
	return p + s
}

func NewPart(v string) Part {
	return Part{
		Path: prefixIfNot(v, "/"),
	}
}

func NewPartList(s ...string) PartList {
	l := make(PartList, len(s))
	for i, v := range s {
		l[i].Path = prefixIfNot(v, "/")
	}
	return l
}

func (x Part) IsZero() bool {
	return x.Path == "" && x.Start.IsZero() && x.Duration.IsZero()
}

func (x Part) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t Part
	return e.EncodeElement(_t(x), node)
}

func (x *Part) UnmarshalText(data []byte) error {
	p := Part{}
	v := string(data)
	tokens := strings.Split(v, "time:")
	p.Path = tokens[0]
	if len(tokens) > 1 {
		t := tokens[1]
		if i := strings.Index(t, "d"); i > -1 {
			// node + duration
			if err := p.Start.UnmarshalText([]byte(t[:i])); err != nil {
				return fmt.Errorf("xmp: invalid part node: %v", err)
			}
			if err := p.Duration.UnmarshalText([]byte(t[i+1:])); err != nil {
				return fmt.Errorf("xmp: invalid part duration: %v", err)
			}
		} else if i = strings.Index(t, "r"); i > -1 {
			// node + end range
			if err := p.Start.UnmarshalText([]byte(t[:i])); err != nil {
				return fmt.Errorf("xmp: invalid part range node: %v", err)
			}
			var end FrameCount
			if err := end.UnmarshalText([]byte(t[i+1:])); err != nil {
				return fmt.Errorf("xmp: invalid part range end: %v", err)
			}
			p.Duration = end.Sub(p.Start)
		} else {
			// duration with zero node
			if err := p.Duration.UnmarshalText([]byte(t)); err != nil {
				return fmt.Errorf("xmp: invalid part duration: %v", err)
			}
		}
	}
	*x = p
	return nil
}

func (x Part) MarshalText() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString(x.Path)
	if len(x.Path) > 0 && !strings.HasSuffix(x.Path, "/") {
		buf.WriteByte('/')
	}
	if x.Duration.IsZero() {
		return buf.Bytes(), nil
	}
	buf.WriteString("time:")
	if !x.Start.IsZero() {
		buf.WriteString(x.Start.String())
		buf.WriteByte('d')
	}
	buf.WriteString(x.Duration.String())
	return buf.Bytes(), nil
}

type PartList []Part

func (x PartList) String() string {
	if b, err := x.MarshalText(); err == nil {
		return string(b)
	}
	return ""
}

func (x *PartList) Add(v string) {
	*x = append(*x, Part{Path: v})
}

func (x *PartList) AddPart(p Part) {
	*x = append(*x, p)
}

func (x *PartList) UnmarshalText(data []byte) error {
	l := make(PartList, 0)
	for _, v := range bytes.Split(data, []byte(";")) {
		p := Part{}
		if err := p.UnmarshalText(v); err != nil {
			return err
		}
		l = append(l, p)
	}
	*x = l
	return nil
}

func (x PartList) MarshalText() ([]byte, error) {
	s := make([][]byte, len(x))
	for i, v := range x {
		b, _ := v.MarshalText()
		if !bytes.HasPrefix(b, []byte("/")) {
			b = append([]byte("/"), b...)
		}
		s[i] = b
	}
	return bytes.Join(s, []byte(";")), nil
}

// 1.2.4.1 ResourceRef - stRef:maskMarkers closed choice
type MaskType string

const (
	MaskAll  MaskType = "All"
	MaskNone MaskType = "None"
)

// 1.2.6.3 FrameCount
// ##<FrameRate>
type FrameCount struct {
	Count int64
	Rate  FrameRate
}

func (x FrameCount) IsZero() bool {
	return x.Count == 0 && x.Rate.IsZero()
}

func (x FrameCount) Sub(y FrameCount) FrameCount {
	return FrameCount{
		Count: x.Count - y.Count,
		Rate:  x.Rate,
	}
}

func (x FrameCount) IsSmaller(y FrameCount) bool {
	return float64(x.Count)*x.Rate.Value() < float64(y.Count)*y.Rate.Value()
}

func (x *FrameCount) UnmarshalText(data []byte) error {
	var err error
	c := FrameCount{
		Count: 0,
		Rate:  FrameRate{1, 1},
	}
	v := string(data)
	if fpos := strings.Index(v, "f"); fpos > -1 {
		// ##f##s## format
		if err := c.Rate.UnmarshalText([]byte(v[fpos:])); err != nil {
			return fmt.Errorf("xmp: invalid frame count '%s': %v", v, err)
		}
		v = v[:fpos]
	}
	// special "maximum"
	if v == "maximum" {
		c.Count = -1
	} else {
		// ## format
		if c.Count, err = strconv.ParseInt(v, 10, 64); err != nil {
			return fmt.Errorf("xmp: invalid frame counter value '%s': %v", v, err)
		}
	}
	*x = c
	return nil
}

func (x FrameCount) String() string {
	switch x.Count {
	case 0:
		return "0"
	case -1:
		return "maximum"
	default:
		buf := bytes.Buffer{}
		buf.WriteString(strconv.FormatInt(x.Count, 10))
		buf.WriteString(x.Rate.String())
		return buf.String()
	}
}

func (x FrameCount) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// 1.2.6.1 beatSpliceStretch
type BeatSpliceStretch struct {
	RiseInDecibel      float64   `xmp:"xmpDM:riseInDecibel,attr"`
	RiseInTimeDuration MediaTime `xmp:"xmpDM:riseInTimeDuration"`
	UseFileBeatsMarker xmp.Bool  `xmp:"xmpDM:useFileBeatsMarker,attr"`
}

func (x BeatSpliceStretch) IsZero() bool {
	return x.RiseInDecibel == 0 && x.RiseInTimeDuration.IsZero() && !x.UseFileBeatsMarker.Value()
}

// 1.2.6.2 CuePointParam
type CuePointParam struct {
	Key   string `xmp:"xmpDM:key,attr"`
	Value string `xmp:"xmpDM:value,attr"`
}

func (x CuePointParam) IsZero() bool {
	return x.Value == "" && x.Key == ""
}

type CuePointParamArray []CuePointParam

func (x CuePointParamArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x CuePointParamArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *CuePointParamArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type VideoFrameRate float64

func (x *VideoFrameRate) UnmarshalText(data []byte) error {
	// try detecting supported strings
	switch strings.ToUpper(string(data)) {
	case "NTSC":
		*x = VideoFrameRate(29.97)
		return nil
	case "PAL":
		*x = VideoFrameRate(25)
		return nil
	}

	// try parsing as rational
	rat := &xmp.Rational{}
	if err := rat.UnmarshalText(data); err == nil {
		*x = VideoFrameRate(rat.Value())
		return nil
	}

	// parse as float value, strip trailing "fps"
	v := strings.TrimSpace(strings.TrimSuffix(string(data), "fps"))
	v = strings.TrimSuffix(v, "i")
	v = strings.TrimSuffix(v, "p")
	if f, err := strconv.ParseFloat(v, 64); err != nil {
		return fmt.Errorf("xmp: invalid video frame rate value '%s': %v", v, err)
	} else {
		*x = VideoFrameRate(f)
	}
	return nil
}

// 1.2.6.4 FrameRate
// "f###"
// "f###s###"
type FrameRate struct {
	Rate int64
	Base int64
}

func (x FrameRate) Value() float64 {
	if x.Base == 0 {
		return 0
	}
	return float64(x.Rate) / float64(x.Base)
}

func (x FrameRate) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('f')
	buf.WriteString(strconv.FormatInt(x.Rate, 10))
	if x.Base > 1 {
		buf.WriteByte('s')
		buf.WriteString(strconv.FormatInt(x.Base, 10))
	}
	return buf.String()
}

func (x FrameRate) IsZero() bool {
	return x.Base == 0 || x.Rate == 0
}

func (x *FrameRate) UnmarshalText(data []byte) error {
	if !bytes.HasPrefix(data, []byte("f")) {
		return fmt.Errorf("xmp: invalid frame rate value '%s'", string(data))
	}
	v := data[1:]
	r := FrameRate{}
	tokens := bytes.Split(v, []byte("s"))
	var err error
	r.Rate, err = strconv.ParseInt(string(tokens[0]), 10, 64)
	if err != nil {
		return fmt.Errorf("xmp: invalid frame rate value '%s': %v", string(data), err)
	}
	if len(tokens) > 1 {
		r.Base, err = strconv.ParseInt(string(tokens[1]), 10, 64)
		if err != nil {
			return fmt.Errorf("xmp: invalid frame rate value '%s': %v", string(data), err)
		}
	}
	if r.Base == 0 {
		r.Base = 1
	}
	*x = r
	return nil
}

func (x FrameRate) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// 1.2.6.5 Marker
type Marker struct {
	Comment        string             `xmp:"xmpDM:comment,attr"`
	CuePointParams CuePointParamArray `xmp:"xmpDM:cuePointParams"`
	CuePointType   string             `xmp:"xmpDM:cuePointType,attr"`
	Duration       FrameCount         `xmp:"xmpDM:duration,attr"`
	Location       xmp.Uri            `xmp:"xmpDM:location"`
	Name           string             `xmp:"xmpDM:name,attr"`
	Probability    float64            `xmp:"xmpDM:probability,attr"`
	Speaker        string             `xmp:"xmpDM:speaker,attr"`
	StartTime      FrameCount         `xmp:"xmpDM:startTime,attr"`
	Target         string             `xmp:"xmpDM:target,attr"`
	Types          MarkerTypeList     `xmp:"xmpDM:type"`
}

func (x Marker) IsZero() bool {
	return x.Comment == "" &&
		len(x.CuePointParams) == 0 &&
		x.CuePointType == "" &&
		x.Duration.IsZero() &&
		x.Location.IsZero() &&
		x.Name == "" &&
		x.Probability == 0 &&
		x.Speaker == "" &&
		x.StartTime.IsZero() &&
		x.Target == "" &&
		len(x.Types) == 0
}

type MarkerList []Marker

func (x MarkerList) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x MarkerList) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *MarkerList) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

func (x MarkerList) Len() int           { return len(x) }
func (x MarkerList) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x MarkerList) Less(i, j int) bool { return x[i].StartTime.IsSmaller(x[j].StartTime) }

func (x *MarkerList) Sort() {
	sort.Sort(x)
}

func (x MarkerList) Filter(types MarkerTypeList) MarkerList {
	l := make(MarkerList, 0)
	for _, v := range x {
		if len(types.Intersect(v.Types)) > 0 {
			l = append(l, v)
		}
	}
	return l
}

type MarkerTypeList []MarkerType

func (x MarkerTypeList) Intersect(y MarkerTypeList) MarkerTypeList {
	m := make(map[MarkerType]bool)
	for _, v := range x {
		m[v] = true
	}
	r := make(MarkerTypeList, 0)
	for _, v := range y {
		if _, ok := m[v]; ok {
			r = append(r, v)
		}
	}
	return r
}

func (x MarkerTypeList) Contains(t MarkerType) bool {
	for _, v := range x {
		if v == t {
			return true
		}
	}
	return false
}

// 1.2.6.6 Media
type Media struct {
	Duration     MediaTime `xmp:"xmpDM:duration"`
	Managed      xmp.Bool  `xmp:"xmpDM:managed,attr"`
	Path         xmp.Uri   `xmp:"xmpDM:path,attr"`
	StartTime    MediaTime `xmp:"xmpDM:startTime"`
	Track        string    `xmp:"xmpDM:track,attr"`
	WebStatement xmp.Uri   `xmp:"xmpDM:webStatement"`
}

type MediaArray []Media

func (x MediaArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x MediaArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *MediaArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// 1.2.6.7 ProjectLink
type ProjectLink struct {
	Path xmp.Uri  `xmp:"xmpDM:path"`
	Type FileType `xmp:"xmpDM:type,attr"`
}

func (x ProjectLink) IsZero() bool {
	return x.Path == "" && x.Type == ""
}

// 1.2.6.8 resampleStretch
type ResampleStretch struct {
	Quality Quality `xmp:"xmpDM:quality,attr"`
}

func (x ResampleStretch) IsZero() bool {
	return x.Quality == ""
}

// 1.2.6.9 Time
type MediaTime struct {
	Scale xmp.Rational `xmp:"xmpDM:scale,attr"`
	Value int64        `xmp:"xmpDM:value,attr"`
}

func (x MediaTime) IsZero() bool {
	return x.Value == 0 && x.Scale.Den == 0
}

// Parse when used as attribute (not standard-conform)
// xmpDM:duration="0:07.680"
func (x *MediaTime) UnmarshalXMPAttr(d *xmp.Decoder, a xmp.Attr) error {
	v := string(a.Value)
	var dur time.Duration
	fields := strings.Split(v, ".")
	if len(fields) > 1 {
		if i, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
			dur = time.Duration(i * int64(math.Pow10(9-len(fields[1]))))
		}
	}
	fields = strings.Split(fields[0], ":")
	var h, m, s int
	switch len(fields) {
	case 3:
		h, _ = strconv.Atoi(fields[0])
		m, _ = strconv.Atoi(fields[1])
		s, _ = strconv.Atoi(fields[2])
	case 2:
		m, _ = strconv.Atoi(fields[0])
		s, _ = strconv.Atoi(fields[1])
	case 1:
		s, _ = strconv.Atoi(fields[0])
	}
	dur += time.Hour*time.Duration(h) + time.Minute*time.Duration(m) + time.Second*time.Duration(s)
	*x = MediaTime{
		Value: int64(dur / time.Millisecond),
		Scale: xmp.Rational{
			Num: 1,
			Den: 1000,
		},
	}
	return nil
}

// Part 2: 1.2.6.10 Timecode

func ParseTimecodeFormat(f float64, isdrop bool) TimecodeFormat {
	switch {
	case 23.975 <= f && f < 23.997:
		return TimecodeFormat23976
	case f == 24:
		return TimecodeFormat24
	case f == 25:
		return TimecodeFormat25
	case 29.96 < f && f < 29.98:
		if isdrop {
			return TimecodeFormat2997
		} else {
			return TimecodeFormat2997ND
		}
	case f == 30:
		return TimecodeFormat30
	case f == 50:
		return TimecodeFormat50
	case 59.93 < f && f < 59.95:
		if isdrop {
			return TimecodeFormat5994
		} else {
			return TimecodeFormat5994
		}
	case f == 60:
		return TimecodeFormat60
	default:
		return TimecodeFormatInvalid
	}
}

type Timecode struct {
	Format      TimecodeFormat `xmp:"xmpDM:timeFormat"`
	Value       string         `xmp:"xmpDM:timeValue"` // hh:mm:ss:ff or hh;mm;ss;ff
	H           int            `xmp:"-"`
	M           int            `xmp:"-"`
	S           int            `xmp:"-"`
	F           int            `xmp:"-"`
	IsDropFrame bool           `xmp:"-"`
}

func (x Timecode) String() string {
	sep := ':'
	if x.IsDropFrame {
		sep = ';'
	}
	buf := bytes.Buffer{}
	if x.H < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.FormatInt(int64(x.H), 10))
	buf.WriteRune(sep)
	if x.M < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.FormatInt(int64(x.M), 10))
	buf.WriteRune(sep)
	if x.S < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.FormatInt(int64(x.S), 10))
	buf.WriteRune(sep)
	if x.F < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.FormatInt(int64(x.F), 10))
	return buf.String()
}

func (x *Timecode) Unpack() error {
	sep := ":"
	if x.IsDropFrame {
		sep = ";"
	}
	fields := strings.Split(x.Value, sep)
	if len(fields) < 4 {
		return fmt.Errorf("found only %d of 4 fields", len(fields))
	}
	var err error
	if x.H, err = strconv.Atoi(fields[0]); err != nil {
		return err
	}
	if x.M, err = strconv.Atoi(fields[1]); err != nil {
		return err
	}
	if x.S, err = strconv.Atoi(fields[2]); err != nil {
		return err
	}
	if x.F, err = strconv.Atoi(fields[3]); err != nil {
		return err
	}
	return nil
}

func (x Timecode) IsZero() bool {
	return x.Value == "" && x.Format == ""
}

func (x Timecode) MarshalJSON() ([]byte, error) {
	x.Value = x.String()
	type _t Timecode
	return json.Marshal(_t(x))
}

func (x *Timecode) UnmarshalJSON(data []byte) error {
	xt := struct {
		Format TimecodeFormat
		Value  string
	}{}
	if err := json.Unmarshal(data, &xt); err != nil {
		return fmt.Errorf("xmp: invalid timecode: %v", err)
	}
	if len(xt.Value) == 0 {
		xt.Value = "00:00:00:00"
	}
	tc := Timecode{
		Value:       xt.Value,
		Format:      xt.Format,
		IsDropFrame: strings.Contains(xt.Value, ";"),
	}
	if err := tc.Unpack(); err != nil {
		return fmt.Errorf("xmp: invalid timecode '%s': %v", tc.Value, err)
	}
	*x = tc
	return nil
}

func (x Timecode) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	x.Value = x.String()
	type _t Timecode
	return e.EncodeElement(_t(x), node)
}

func (x *Timecode) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	// Note: to avoid recursion we unmarshal into a temporary struct type
	//       and copy over the result

	// support both the attribute form and the element form
	var tc Timecode
	xt := struct {
		Format TimecodeFormat `xmp:"xmpDM:timeFormat"`
		Value  string         `xmp:"xmpDM:timeValue"`
	}{}
	if err := d.DecodeElement(&xt, node); err != nil {
		return fmt.Errorf("xmp: invalid timecode '%s': %v", node.Value, err)
	}
	if len(xt.Value) == 0 {
		xt.Value = "00:00:00:00"
	}
	tc = Timecode{
		Value:       xt.Value,
		Format:      xt.Format,
		IsDropFrame: strings.Contains(xt.Value, ";"),
	}
	if err := tc.Unpack(); err != nil {
		return fmt.Errorf("xmp: invalid timecode '%s': %v", tc.Value, err)
	}
	*x = tc
	return nil
}

// 1.2.6.11 timeScaleStretch
type TimeScaleStretch struct {
	Overlap   float64 `xmp:"xmpDM:frameOverlappingPercentage,attr"`
	FrameSize float64 `xmp:"xmpDM:frameSize,attr"`
	Quality   Quality `xmp:"xmpDM:quality,attr"`
}

func (x TimeScaleStretch) IsZero() bool {
	return x.Overlap == 0 && x.FrameSize == 0 && x.Quality == ""
}

// 1.2.6.12 Track
type Track struct {
	FrameRate FrameRate      `xmp:"xmpDM:frameRate,attr"`
	Markers   MarkerList     `xmp:"xmpDM:markers"`
	Name      string         `xmp:"xmpDM:trackName,attr"`
	Type      MarkerTypeList `xmp:"xmpDM:trackType,attr"`
}

type TrackArray []Track

func (x TrackArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x TrackArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *TrackArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}
