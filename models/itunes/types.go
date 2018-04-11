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

// iTunes-style metadata as found in .mp4, .m4a, .m4p, .m4v, .m4b files
package itunes

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// https://sourceforge.net/p/mediainfo/feature-requests/398/
// https://forums.mp3tag.de/index.php?showtopic=12640
type SMPB struct {
	EncoderDelay        int64 `xmp:"iTunes:EncoderDelay,attr"`
	EndPadding          int64 `xmp:"iTunes:EndPadding,attr"`
	OriginalSampleCount int64 `xmp:"iTunes:OriginalSampleCount,attr"`
	EndOffset           int64 `xmp:"iTunes:EndOffset,attr"`
}

// Gapless Playback info
// iTunSMPB 00000000 00000840 000001CA 00000000003F31F6 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
//                   Delay    Padding  SampleCount               EndOffset
func (x *SMPB) UnmarshalText(data []byte) error {
	var err error
	sm := SMPB{}
	for i, v := range bytes.Fields(data) {
		switch i {
		case 1:
			if sm.EncoderDelay, err = strconv.ParseInt(string(v), 16, 32); err != nil {
				return fmt.Errorf("iTunes: cannot parse SMPB '%s': %v", string(data), err)
			}
		case 2:
			if sm.EndPadding, err = strconv.ParseInt(string(v), 16, 32); err != nil {
				return fmt.Errorf("iTunes: cannot parse SMPB '%s': %v", string(data), err)
			}
		case 3:
			if sm.OriginalSampleCount, err = strconv.ParseInt(string(v), 16, 64); err != nil {
				return fmt.Errorf("iTunes: cannot parse SMPB '%s': %v", string(data), err)
			}
		case 5:
			if sm.EndOffset, err = strconv.ParseInt(string(v), 16, 64); err != nil {
				return fmt.Errorf("iTunes: cannot parse SMPB '%s': %v", string(data), err)
			}
		}
	}
	*x = sm
	return nil
}

type Bool int // 0 no, 1 yes

func (b Bool) Value() bool {
	return b == 1
}

type ContentRating struct {
	Standard string `xmp:"iTunes:Standard,attr"`
	Rating   string `xmp:"iTunes:Rating,attr"`
	Score    string `xmp:"iTunes:Score,attr"`
	Reasons  string `xmp:"iTunes:Reasons,attr"`
}

func (x ContentRating) IsZero() bool {
	return x.Standard == "" && x.Rating == ""
}

func (x ContentRating) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('(')
	buf.WriteString(x.Standard)
	buf.WriteByte('|')
	buf.WriteString(x.Rating)
	buf.WriteByte('|')
	buf.WriteString(x.Score)
	buf.WriteByte('|')
	buf.WriteString(x.Reasons)
	return buf.String()
}

func (x *ContentRating) UnmarshalText(data []byte) error {
	cr := ContentRating{}
	v := string(data)
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "(")
	v = strings.TrimSuffix(v, ")")
	for i, s := range strings.Split(v, "|") {
		switch i {
		case 1:
			cr.Standard = s
		case 2:
			cr.Rating = s
		case 3:
			cr.Score = s
		case 4:
			cr.Reasons = s
		}
	}
	*x = cr
	return nil
}

func (x ContentRating) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(x.String()), nil
}
