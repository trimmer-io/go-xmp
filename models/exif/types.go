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

// http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf

package exif

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"trimmer.io/go-xmp/xmp"
)

const EXIF_DATE_FORMAT = "2006:01:02 15:04:05"

type ByteArray []byte

func (x ByteArray) String() string {
	s := make([]string, len(x))
	for i, v := range x {
		s[i] = strconv.Itoa(int(v))
	}
	return strings.Join(s, " ")
}

func (x ByteArray) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *ByteArray) UnmarshalText(data []byte) error {
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return nil
	}
	a := make(ByteArray, len(fields))
	for i, v := range fields {
		if j, err := strconv.ParseUint(v, 10, 8); err != nil {
			return fmt.Errorf("exif: invalid byte '%s': %v", v, err)
		} else {
			a[i] = byte(j)
		}
	}
	*x = a
	return nil
}

type Date time.Time

func Now() Date {
	return Date(time.Now())
}

func (x Date) IsZero() bool {
	return time.Time(x).IsZero()
}

func (x Date) Value() time.Time {
	return time.Time(x)
}

func (x Date) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(time.Time(x).Format(EXIF_DATE_FORMAT)), nil
}

func (x Date) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return e.EncodeElement(time.Time(x).Format(EXIF_DATE_FORMAT), node)
}

var timeFormats []string = []string{
	"2006-01-02T15:04:05.999999999",       // XMP
	"2006-01-02T15:04:05.999999999Z07:00", // XMP
	"2006-01-02T15:04:05.999999999Z",      // XMP
	"2006-01-02T15:04:05Z",                // XMP
	"2006-01-02T15:04:05",                 // XMP
	"2006-01-02T15:04",                    // XMP
	"2006:01:02 15:04:05.999999999-07:00", // Exif
	"2006:01:02 15:04:05-07:00",           // Exif
	"2006:01:02 15:04:05",                 // Exif
	"2006:01:02 15:04",                    // Exif
	"2006:01:02",                          // Exif GPS datestamp
	"2006:01",                             // Exif
	"2006",                                // Exif
}

func (x *Date) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	value := string(data)
	for _, f := range timeFormats {
		if t, err := time.Parse(f, value); err == nil {
			*x = Date(t)
			return nil
		}
	}
	return fmt.Errorf("exif: invalid datetime value '%s'", value)
}

func convertDateToXMP(d Date, s string) (xmp.Date, error) {
	if d.IsZero() {
		return xmp.Date{}, nil
	}
	xd := xmp.Date(d)
	s = strings.TrimSpace(s)
	if s != "" {
		return xd, nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err != nil {
		return xd, fmt.Errorf("exif: invalid subsecond format '%s': %v", s, err)
	} else {
		sec := d.Value().Unix()
		nsec := i * int64(math.Pow10(9-len(s)))
		xd = xmp.Date(time.Unix(sec, nsec))
	}
	return xd, nil
}

type ComponentArray []Component

func (x ComponentArray) String() string {
	s := make([]string, len(x))
	for i, v := range x {
		s[i] = strconv.Itoa(int(v))
	}
	return strings.Join(s, " ")
}

func (x ComponentArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x ComponentArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *ComponentArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type OECF struct {
	Columns int               `xmp:"exif:Columns"`
	Rows    int               `xmp:"exif:Rows"`
	Names   xmp.StringArray   `xmp:"exif:Names"`
	Values  xmp.RationalArray `xmp:"exif:Values"`
}

func (x OECF) IsZero() bool {
	return x.Columns == 0 || x.Rows == 0
}

func (x *OECF) Addr() *OECF {
	if x.IsZero() {
		return nil
	}
	return x
}

func (x *OECF) UnmarshalText(data []byte) error {
	var err error
	o := OECF{}
	for i, v := range strings.Fields(string(data)) {
		switch {
		case i == 0:
			if o.Columns, err = strconv.Atoi(v); err != nil {
				return fmt.Errorf("exif: invalid OECF columns value '%s': %v", v, err)
			}
			o.Names = make(xmp.StringArray, 0, o.Columns)
		case i == 1:
			if o.Rows, err = strconv.Atoi(v); err != nil {
				return fmt.Errorf("exif: invalid OECF rows value '%s': %v", v, err)
			}
			o.Values = make(xmp.RationalArray, 0, o.Columns*o.Rows)
		case i > 1 && i < o.Columns+1:
			o.Names = append(o.Names, v)
		default:
			r := xmp.Rational{}
			if err := r.UnmarshalText([]byte(v)); err != nil {
				return fmt.Errorf("exif: invalid OECF value '%s': %v", v, err)
			}
			o.Values = append(o.Values, r)
		}
	}
	*x = o
	return nil
}

func (x OECF) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	type _t OECF
	return e.EncodeElement(_t(x), node)
}

type CFAPattern struct {
	Columns int       `xmp:"exif:Columns"`
	Rows    int       `xmp:"exif:Rows"`
	Values  ByteArray `xmp:"exif:Values"`
}

func (x CFAPattern) IsZero() bool {
	return x.Columns == 0 || x.Rows == 0
}

func (x *CFAPattern) Addr() *CFAPattern {
	if x.IsZero() {
		return nil
	}
	return x
}

func (x *CFAPattern) UnmarshalText(data []byte) error {
	var err error
	c := CFAPattern{}
	for i, v := range strings.Fields(string(data)) {
		switch {
		case i == 0:
			if c.Columns, err = strconv.Atoi(v); err != nil {
				return fmt.Errorf("exif: invalid CFA columns value '%s': %v", v, err)
			}
		case i == 1:
			if c.Rows, err = strconv.Atoi(v); err != nil {
				return fmt.Errorf("exif: invalid CFA rows value '%s': %v", v, err)
			}
			c.Values = make(ByteArray, 0, c.Columns*c.Rows)
		default:
			if j, err := strconv.ParseUint(v, 10, 8); err != nil {
				return fmt.Errorf("exif: invalid CFA value '%s': %v", v, err)
			} else {
				c.Values = append(c.Values, byte(j))
			}
		}
	}
	*x = c
	return nil
}

func (x CFAPattern) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t CFAPattern
	return e.EncodeElement(_t(x), node)
}

type DeviceSettings struct {
	Columns int             `xmp:"exif:Columns"`
	Rows    int             `xmp:"exif:Rows"`
	Values  xmp.StringArray `xmp:"exif:Values"`
}

func (x DeviceSettings) IsZero() bool {
	return x.Columns == 0 || x.Rows == 0
}

func (x *DeviceSettings) Addr() *DeviceSettings {
	if x.IsZero() {
		return nil
	}
	return x
}

func (x *DeviceSettings) UnmarshalText(data []byte) error {
	var err error
	s := DeviceSettings{}
	for i, v := range strings.Fields(string(data)) {
		switch {
		case i == 0:
			if s.Columns, err = strconv.Atoi(v); err != nil {
				return fmt.Errorf("exif: invalid device settings columns value '%s': %v", v, err)
			}
		case i == 1:
			if s.Rows, err = strconv.Atoi(v); err != nil {
				return fmt.Errorf("exif: invalid device settings rows value '%s': %v", v, err)
			}
			s.Values = make(xmp.StringArray, 0, s.Columns*s.Rows)
		default:
			s.Values = append(s.Values, v)
		}
	}
	*x = s
	return nil
}

func (x DeviceSettings) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t DeviceSettings
	return e.EncodeElement(_t(x), node)
}

type Flash struct {
	Fired      xmp.Bool        `xmp:"exif:Fired,attr,empty"`
	Function   xmp.Bool        `xmp:"exif:Function,attr,empty"`
	Mode       FlashMode       `xmp:"exif:Mode,attr,empty"`
	RedEyeMode xmp.Bool        `xmp:"exif:RedEyeMode,attr,empty"`
	Return     FlashReturnMode `xmp:"exif:Return,attr,empty"`
}

func (x Flash) IsZero() bool {
	return x.Mode == 0 && x.Return == 0 && !x.Fired.Value() && !x.Function.Value() && !x.RedEyeMode.Value()
}

func (x *Flash) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	v, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		return fmt.Errorf("exif: invalid flash value '%d': %v", v, err)
	}
	x.Fired = v&0x01 > 0
	x.Return = FlashReturnMode(v >> 1 & 0x3)
	x.Mode = FlashMode(v >> 3 & 0x03)
	x.Function = v&0x20 > 0
	x.RedEyeMode = v&0x40 > 0
	return nil
}

func (x Flash) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t Flash
	return e.EncodeElement(_t(x), node)
}

func (x *Flash) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	// strip TextUnmarhal method
	type _t Flash
	f := _t{}
	if err := d.DecodeElement(&f, node); err != nil {
		return err
	}
	*x = Flash(f)
	return nil
}

type GPSCoord [3]xmp.Rational

// unmarshal from decimal
func (x *GPSCoord) UnmarshalText(data []byte) error {
	s := string(data)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || s[0] != '-' && len(strings.Split(s, ".")[0]) > 3 {
		return fmt.Errorf("exif: invalid decimal GPS coordinate '%s'", s)
	}
	val := math.Abs(f)
	degrees := int(math.Floor(val))
	minutes := int(math.Floor(60 * (val - float64(degrees))))
	seconds := 3600 * (val - float64(degrees) - (float64(minutes) / 60))
	(*x)[0] = xmp.FloatToRational(float32(degrees))
	(*x)[1] = xmp.FloatToRational(float32(minutes))
	(*x)[2] = xmp.FloatToRational(float32(seconds))
	return nil
}

func convertGPStoXMP(r GPSCoord, ref string) (xmp.GPSCoord, error) {
	if ref == "" {
		return "", nil
	}
	var deg [3]float64
	for i, v := range r {
		if v.Den == 0 {
			return "", fmt.Errorf("exif: invalid GPS coordinate '%s'", v.String())
		}
		deg[i] = float64(v.Num) / float64(v.Den)
	}
	min := deg[0]*60 + deg[1] + deg[2]/60
	ideg := int(min / 60)
	min -= float64(ideg) * 60

	buf := bytes.Buffer{}
	buf.Write(strconv.AppendInt(make([]byte, 0, 3), int64(ideg), 10))
	buf.WriteByte(',')
	buf.Write(strconv.AppendFloat(make([]byte, 0, 24), min, 'f', 7, 64))
	buf.WriteString(ref)
	return xmp.GPSCoord(buf.String()), nil
}

func convertGPSTimestamp(date Date, x xmp.RationalArray) xmp.Date {
	if date.IsZero() {
		return xmp.Date{}
	}
	var d time.Duration
	for i, v := range x {
		switch i {
		case 0:
			d += time.Duration(v.Value() * float64(time.Hour))
		case 1:
			d += time.Duration(v.Value() * float64(time.Minute))
		case 2:
			d += time.Duration(v.Value() * float64(time.Second))
		}
	}
	return xmp.Date(date.Value().Add(d))
}
