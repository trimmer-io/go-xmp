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

// http://forums.camerabits.com/index.php?topic=903.0

// Package pm implements metadata written by Photomechanic software found in many professional images.
package pm

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsPm = xmp.NewNamespace("photomechanic", "http://ns.camerabits.com/photomechanic/1.0/", NewModel)
)

func init() {
	xmp.Register(NsPm, xmp.ImageMetadata)
}

func NewModel(name string) xmp.Model {
	return &Photomechanic{}
}

func MakeModel(d *xmp.Document) (*Photomechanic, error) {
	m, err := d.MakeModel(NsPm)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*Photomechanic)
	return x, nil
}

func FindModel(d *xmp.Document) *Photomechanic {
	if m := d.FindModel(NsPm); m != nil {
		return m.(*Photomechanic)
	}
	return nil
}

func (x Photomechanic) Can(nsName string) bool {
	return nsName == "photomechanic"
}

func (x *Photomechanic) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsPm}
}

func (x *Photomechanic) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *Photomechanic) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x Photomechanic) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *Photomechanic) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *Photomechanic) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsPm.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *Photomechanic) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsPm.GetName(), err)
	}
	return nil
}

type Photomechanic struct {
	EditStatus        string          `xmp:"photomechanic:EditStatus"`        // edit status field from IPTC record
	ColorClass        int             `xmp:"photomechanic:ColorClass"`        // USA
	CountryCode       string          `xmp:"photomechanic:CountryCode"`       // USA
	TimeCreated       Time            `xmp:"photomechanic:TimeCreated"`       // HHMMSS+dHdS
	Prefs             Preferences     `xmp:"photomechanic:Prefs"`             // T:C:R:F
	Tagged            xmp.Bool        `xmp:"photomechanic:Tagged"`            // "True"
	Version           string          `xmp:"photomechanic:PMVersion"`         // PM5
	ColorClassEval    int             `xmp:"photomechanic:ColorClassEval"`    // ="2"
	ColorClassApply   xmp.Bool        `xmp:"photomechanic:ColorClassApply"`   // ="True"
	RatingEval        int             `xmp:"photomechanic:RatingEval"`        // ="4"
	RatingApply       xmp.Bool        `xmp:"photomechanic:RatingApply"`       // ="True"
	TagEval           int             `xmp:"photomechanic:TagEval"`           // ="0"
	TagApply          xmp.Bool        `xmp:"photomechanic:TagApply"`          // ="False"
	CaptionMergeStyle int             `xmp:"photomechanic:CaptionMergeStyle"` // ="1"
	ApplyAPCustom     int             `xmp:"photomechanic:ApplyAPCustom"`     // ="0"
	MergeAPCustom     int             `xmp:"photomechanic:MergeAPCustom"`     // ="0"
	ApplyDateType     int             `xmp:"photomechanic:ApplyDateType"`     // ="2"
	FieldsToApply     xmp.StringArray `xmp:"photomechanic:FieldsToApply"`
}

type Time time.Time

const pmTimeFormat = "150405-0700"

// HH = 24hr hour value (0-23)
// MM = minute value (0-59)
// SS = second value (0-59)
// dH = GMT hour delta (+/- 12)
// dS = GMT minute delta (0-59)

func (x Time) IsZero() bool {
	return time.Time(x).IsZero()
}

func (x Time) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(time.Time(x).Format(pmTimeFormat)), nil
}

func (x *Time) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if t, err := time.Parse(pmTimeFormat, string(data)); err != nil {
		return fmt.Errorf("photomechanic: invalid time value '%s'", string(data))
	} else {
		*x = Time(t)
	}
	return nil
}

// T:C:R:F
type Preferences struct {
	TagStatus  int // T = Tag status (0 or 1)
	ColorClass int // C = Color Class value (0-8)
	Rating     int // R = Rating (0-5)
	Frame      int // F = Frame number of image or -1 if undetermined
}

func (x Preferences) IsZero() bool {
	return x.TagStatus == 0 && x.ColorClass == 0 && x.Rating == 0 && x.Frame == 0
}

func (x Preferences) MarshalText() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString(strconv.FormatInt(int64(x.TagStatus), 10))
	buf.WriteByte(':')
	buf.WriteString(strconv.FormatInt(int64(x.ColorClass), 10))
	buf.WriteByte(':')
	buf.WriteString(strconv.FormatInt(int64(x.Rating), 10))
	buf.WriteByte(':')
	buf.WriteString(strconv.FormatInt(int64(x.Frame), 10))
	return buf.Bytes(), nil
}

func (x *Preferences) UnmarshalText(data []byte) error {
	p := Preferences{}
	var err error
	fields := strings.Split(string(data), ":")
	if len(fields) != 4 {
		return fmt.Errorf("photomechanic: invalid prefs value '%s'", string(data))
	}
	if p.TagStatus, err = strconv.Atoi(fields[0]); err != nil {
		return fmt.Errorf("photomechanic: invalid tag status value '%s': %v", string(data), err)
	}
	if p.ColorClass, err = strconv.Atoi(fields[1]); err != nil {
		return fmt.Errorf("photomechanic: invalid color class value '%s': %v", string(data), err)
	}
	if p.Rating, err = strconv.Atoi(fields[2]); err != nil {
		return fmt.Errorf("photomechanic: invalid rating value '%s': %v", string(data), err)
	}
	if p.Frame, err = strconv.Atoi(fields[3]); err != nil {
		return fmt.Errorf("photomechanic: invalid frame number value '%s': %v", string(data), err)
	}
	*x = p
	return nil
}
