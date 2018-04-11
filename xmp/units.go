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

// Duration
// NullInt
// SplitInt (1:2:3:4)
// NullString
// NullBool
// NullFloat
// NullFloat64

package xmp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var EEmptyValue = errors.New("empty value")

// Int
type NullInt int

func (i NullInt) Int() int {
	return int(i)
}

func (i NullInt) Value() int {
	return int(i)
}

func ParseNullInt(d string) (NullInt, error) {
	switch d {
	case "", "-", "--", "---", "NaN", "unknown":
		return 0, EEmptyValue
	default:
		// parse integer
		if i, err := strconv.ParseInt(d, 10, 64); err == nil {
			return NullInt(i), nil
		}
		return 0, fmt.Errorf("xmp: parsing NullInt '%s': invalid syntax", d)
	}
}

func (i NullInt) MarshalText() ([]byte, error) {
	return []byte(strconv.Itoa(int(i))), nil
}

func (i *NullInt) UnmarshalText(data []byte) error {
	v, err := ParseNullInt(string(data))
	if err == EEmptyValue {
		return nil
	}
	if err != nil {
		return err
	}
	*i = v
	return nil
}

func (i NullInt) MarshalJSON() ([]byte, error) {
	return i.MarshalText()
}

func (i *NullInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return i.UnmarshalText(bytes.Trim(data, "\""))
	}
	return i.UnmarshalText(data)
}

// 32bit Integer split byte-wise (06:13:14:61)
type SplitInt uint32

func (x SplitInt) Uint32() uint32 {
	return uint32(x)
}

func (x SplitInt) Value() uint32 {
	return uint32(x)
}

func (x *SplitInt) UnmarshalText(data []byte) error {
	d := string(data)
	switch d {
	case "", "-", "--", "---", "NaN", "unknown":
		return nil
	default:
		// parse integer
		if i, err := strconv.ParseInt(d, 10, 32); err == nil {
			*x = SplitInt(i)
		} else if f := strings.Split(d, ":"); len(f) == 4 {
			var u uint32
			for i := 0; i < 4; i++ {
				if v, err := strconv.ParseInt(f[i], 10, 8); err != nil {
					return fmt.Errorf("xmp: parsing SplitInt '%s': invalid syntax", d)
				} else {
					u += uint32(v) << uint32((3-i)*8)
				}
			}
			*x = SplitInt(u)
		} else {
			return fmt.Errorf("xmp: parsing SplitInt '%s': invalid syntax", d)
		}
	}
	return nil
}

// String
type NullString string

func (s NullString) String() string {
	return string(s)
}

func (s NullString) Value() string {
	return string(s)
}

func ParseNullString(s string) (NullString, error) {
	s = strings.TrimSpace(s)
	switch s {
	case "", "-", "--", "unknown":
		return "", EEmptyValue
	default:
		return NullString(s), nil
	}
}

func (s NullString) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *NullString) UnmarshalText(data []byte) error {
	// only overwrite when not set
	if len(*s) > 0 {
		return nil
	}
	// only overwrite when not empty
	p, err := ParseNullString(string(data))
	if err == EEmptyValue {
		return nil
	}
	if err != nil {
		return err
	}
	*s = p
	return nil
}

// Bool
type NullBool bool

func (x NullBool) Value() bool {
	return bool(x)
}

func ParseNullBool(d string) (NullBool, error) {
	switch strings.ToLower(d) {
	case "", "-", "--", "---":
		return NullBool(false), EEmptyValue
	case "true", "on", "yes", "1", "enabled":
		return NullBool(true), nil
	default:
		return NullBool(false), nil
	}
}

func (x NullBool) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatBool(bool(x))), nil
}

func (x *NullBool) UnmarshalText(data []byte) error {
	v, err := ParseNullBool(string(data))
	if err == EEmptyValue {
		return nil
	}
	if err != nil {
		return err
	}
	*x = v
	return nil
}

func (x NullBool) MarshalJSON() ([]byte, error) {
	return x.MarshalText()
}

func (x *NullBool) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return x.UnmarshalText(bytes.Trim(data, "\""))
	}
	return x.UnmarshalText(data)
}

// Duration
type Duration time.Duration

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) String() string {
	return time.Duration(d).String()
}

func ParseDuration(d string) (Duration, error) {
	// parse integer values as seconds
	if i, err := strconv.ParseInt(d, 10, 64); err == nil {
		return Duration(time.Duration(i) * time.Second), nil
	}
	// parse as duration string (note: no whitespace allowed)
	if i, err := time.ParseDuration(d); err == nil {
		return Duration(i), nil
	}
	// parse as duration string with whitespace removed
	d = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, d)
	if i, err := time.ParseDuration(d); err == nil {
		return Duration(i), nil
	}
	return 0, fmt.Errorf("xmp: parsing duration '%s': invalid syntax", d)
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

func (d *Duration) UnmarshalText(data []byte) error {
	i, err := ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = i
	return nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return d.UnmarshalText(bytes.Trim(data, "\""))
	}
	if i, err := strconv.ParseInt(string(data), 10, 64); err == nil {
		*d = Duration(time.Duration(i) * time.Second)
		return nil
	}
	return fmt.Errorf("xmp: parsing duration '%s': invalid syntax", string(data))
}

func (d Duration) Truncate(r time.Duration) Duration {
	if d > 0 {
		return Duration(math.Ceil(float64(d)/float64(r))) * Duration(r)
	} else {
		return Duration(math.Floor(float64(d)/float64(r))) * Duration(r)
	}
}

func (d Duration) RoundToDays() int {
	return int(d.Truncate(time.Hour*24) / Duration(time.Hour*24))
}

func (d Duration) RoundToHours() int64 {
	return int64(d.Truncate(time.Hour) / Duration(time.Hour))
}

func (d Duration) RoundToMinutes() int64 {
	return int64(d.Truncate(time.Minute) / Duration(time.Minute))
}

func (d Duration) RoundToSeconds() int64 {
	return int64(d.Truncate(time.Second) / Duration(time.Second))
}

func (d Duration) RoundToMillisecond() int64 {
	return int64(d.Truncate(time.Millisecond) / Duration(time.Millisecond))
}

// Float32
type NullFloat float32

func (f NullFloat) Float32() float32 {
	return float32(f)
}

func (f NullFloat) Value() float64 {
	return float64(f)
}

var cleanFloat = regexp.MustCompile("[^0-9e.+-]")

func ParseNullFloat(d string) (NullFloat, error) {
	switch strings.ToLower(d) {
	case "", "-", "--", "---", "unknown":
		return 0, EEmptyValue
	case "nan":
		return NullFloat(math.NaN()), nil
	case "-inf":
		return NullFloat(math.Inf(-1)), nil
	case "+inf", "inf":
		return NullFloat(math.Inf(1)), nil
	default:
		// parse float
		c := cleanFloat.ReplaceAllString(d, "")
		if i, err := strconv.ParseFloat(c, 32); err == nil {
			return NullFloat(i), nil
		}
		return 0, fmt.Errorf("xmp: parsing NullFloat '%s': invalid syntax", d)
	}
}

func (f NullFloat) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatFloat(f.Value(), 'f', -1, 32)), nil
}

func (f *NullFloat) UnmarshalText(data []byte) error {
	v, err := ParseNullFloat(string(data))
	if err == EEmptyValue {
		return nil
	}
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// JSON does not define NaN, +Inf and -Inf
func (f NullFloat) MarshalJSON() ([]byte, error) {
	switch {
	case math.IsNaN(f.Value()):
		return []byte("\"NaN\""), nil
	case math.IsInf(f.Value(), -1):
		return []byte("\"+Inf\""), nil
	case math.IsInf(f.Value(), 1):
		return []byte("\"-Inf\""), nil
	default:
		return json.Marshal(f.Value())
	}
}

func (f *NullFloat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return f.UnmarshalText(bytes.Trim(data, "\""))
	}
	return f.UnmarshalText(data)
}

// Float64
type NullFloat64 float64

func (f NullFloat64) Float64() float64 {
	return float64(f)
}

func (f NullFloat64) Value() float64 {
	return float64(f)
}

func ParseNullFloat64(d string) (NullFloat64, error) {
	switch strings.ToLower(d) {
	case "", "-", "--", "---", "unknown":
		return 0, EEmptyValue
	case "nan":
		return NullFloat64(math.NaN()), nil
	case "-inf":
		return NullFloat64(math.Inf(-1)), nil
	case "+inf", "inf":
		return NullFloat64(math.Inf(1)), nil
	default:
		// parse float
		c := cleanFloat.ReplaceAllString(d, "")
		if i, err := strconv.ParseFloat(c, 64); err == nil {
			return NullFloat64(i), nil
		}
		return 0, fmt.Errorf("xmp: parsing NullFloat64 '%s': invalid syntax", d)
	}
}

func (f NullFloat64) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatFloat(f.Value(), 'f', -1, 64)), nil
}

func (f *NullFloat64) UnmarshalText(data []byte) error {
	v, err := ParseNullFloat64(string(data))
	if err == EEmptyValue {
		return nil
	}
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// JSON does not define NaN, +Inf and -Inf
func (f NullFloat64) MarshalJSON() ([]byte, error) {
	switch {
	case math.IsNaN(f.Value()):
		return []byte("\"NaN\""), nil
	case math.IsInf(f.Value(), -1):
		return []byte("\"+Inf\""), nil
	case math.IsInf(f.Value(), 1):
		return []byte("\"-Inf\""), nil
	default:
		return json.Marshal(f.Value())
	}
}

func (f *NullFloat64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return f.UnmarshalText(bytes.Trim(data, "\""))
	}
	return f.UnmarshalText(data)
}

func Max(x, y int) int {
	if x < y {
		return y
	} else {
		return x
	}
}

func Min(x, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}

func Max64(x, y int64) int64 {
	if x < y {
		return y
	} else {
		return x
	}
}

func Min64(x, y int64) int64 {
	if x > y {
		return y
	} else {
		return x
	}
}
