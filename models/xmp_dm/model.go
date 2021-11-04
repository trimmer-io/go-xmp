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
	"fmt"
	"github.com/trimmer-io/go-xmp/models/xmp_tpg"
	"github.com/trimmer-io/go-xmp/xmp"
)

var (
	NsXmpDM = xmp.NewNamespace("xmpDM", "http://ns.adobe.com/xmp/1.0/DynamicMedia/", NewModel)
)

func init() {
	xmp.Register(NsXmpDM, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &XmpDM{}
}

func MakeModel(d *xmp.Document) (*XmpDM, error) {
	m, err := d.MakeModel(NsXmpDM)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*XmpDM)
	return x, nil
}

func FindModel(d *xmp.Document) *XmpDM {
	if m := d.FindModel(NsXmpDM); m != nil {
		return m.(*XmpDM)
	}
	return nil
}

type XmpDM struct {
	AbsPeakAudioFilePath         string            `xmp:"xmpDM:absPeakAudioFilePath"`
	Album                        string            `xmp:"xmpDM:album"`
	AltTapeName                  string            `xmp:"xmpDM:altTapeName"`
	AltTimecode                  Timecode          `xmp:"xmpDM:altTimecode"`
	Artist                       string            `xmp:"xmpDM:artist"`
	AudioChannelType             AudioChannelType  `xmp:"xmpDM:audioChannelType"`
	AudioCompressor              string            `xmp:"xmpDM:audioCompressor"`
	AudioSampleRate              float64           `xmp:"xmpDM:audioSampleRate"`
	AudioSampleType              AudioSampleType   `xmp:"xmpDM:audioSampleType"`
	BeatSpliceParams             BeatSpliceStretch `xmp:"xmpDM:beatSpliceParams"`
	CameraAngle                  CameraAngle       `xmp:"xmpDM:cameraAngle"`
	CameraLabel                  string            `xmp:"xmpDM:cameraLabel"`
	CameraModel                  string            `xmp:"xmpDM:cameraModel"`
	CameraMove                   CameraMove        `xmp:"xmpDM:cameraMove"`
	Client                       string            `xmp:"xmpDM:client"`
	Comment                      string            `xmp:"xmpDM:comment"`
	Composer                     string            `xmp:"xmpDM:composer"`
	ContributedMedia             MediaArray        `xmp:"xmpDM:contributedMedia"`
	Director                     string            `xmp:"xmpDM:director"`
	DirectorPhotography          string            `xmp:"xmpDM:directorPhotography"`
	Duration                     MediaTime         `xmp:"xmpDM:duration"`
	Engineer                     string            `xmp:"xmpDM:engineer"`
	FileDataRate                 xmp.Rational      `xmp:"xmpDM:fileDataRate"`
	Genre                        string            `xmp:"xmpDM:genre"`
	Good                         xmp.Bool          `xmp:"xmpDM:good"`
	Instrument                   string            `xmp:"xmpDM:instrument"`
	IntroTime                    MediaTime         `xmp:"xmpDM:introTime"`
	Key                          MusicalKey        `xmp:"xmpDM:key"`
	LogComment                   string            `xmp:"xmpDM:logComment"`
	Loop                         xmp.Bool          `xmp:"xmpDM:loop"`
	NumberOfBeats                float64           `xmp:"xmpDM:numberOfBeats"`
	Markers                      MarkerList        `xmp:"xmpDM:markers"`
	OutCue                       MediaTime         `xmp:"xmpDM:outCue"`
	ProjectName                  string            `xmp:"xmpDM:projectName"`
	ProjectRef                   ProjectLink       `xmp:"xmpDM:projectRef"`
	PullDown                     PullDown          `xmp:"xmpDM:pullDown"`
	RelativePeakAudioFilePath    xmp.Uri           `xmp:"xmpDM:relativePeakAudioFilePath"`
	RelativeTimestamp            MediaTime         `xmp:"xmpDM:relativeTimestamp"`
	ReleaseDate                  xmp.Date          `xmp:"xmpDM:releaseDate"`
	ResampleParams               ResampleStretch   `xmp:"xmpDM:resampleParams"`
	ScaleType                    MusicalScale      `xmp:"xmpDM:scaleType"`
	Scene                        string            `xmp:"xmpDM:scene"`
	ShotDate                     xmp.Date          `xmp:"xmpDM:shotDate"`
	ShotDay                      string            `xmp:"xmpDM:shotDay"`
	ShotLocation                 string            `xmp:"xmpDM:shotLocation"`
	ShotName                     string            `xmp:"xmpDM:shotName"`
	ShotNumber                   string            `xmp:"xmpDM:shotNumber"`
	ShotSize                     ShotSize          `xmp:"xmpDM:shotSize"`
	SpeakerPlacement             string            `xmp:"xmpDM:speakerPlacement"`
	StartTimecode                Timecode          `xmp:"xmpDM:startTimecode"`
	StretchMode                  StretchMode       `xmp:"xmpDM:stretchMode"`
	TakeNumber                   int               `xmp:"xmpDM:takeNumber"`
	TapeName                     string            `xmp:"xmpDM:tapeName"`
	Tempo                        float64           `xmp:"xmpDM:tempo"`
	TimeScaleParams              TimeScaleStretch  `xmp:"xmpDM:timeScaleParams"`
	TimeSignature                TimeSignature     `xmp:"xmpDM:timeSignature"`
	TrackNumber                  int               `xmp:"xmpDM:trackNumber"`
	Tracks                       TrackArray        `xmp:"xmpDM:Tracks"`
	VideoAlphaMode               AlphaMode         `xmp:"xmpDM:videoAlphaMode"`
	VideoAlphaPremultipleColor   xmptpg.Colorant   `xmp:"xmpDM:videoAlphaPremultipleColor"`
	VideoAlphaUnityIsTransparent xmp.Bool          `xmp:"xmpDM:videoAlphaUnityIsTransparent"`
	VideoColorSpace              ColorSpace        `xmp:"xmpDM:videoColorSpace"`
	VideoCompressor              string            `xmp:"xmpDM:videoCompressor"`
	VideoFieldOrder              FieldOrder        `xmp:"xmpDM:videoFieldOrder"`
	VideoFrameRate               VideoFrameRate    `xmp:"xmpDM:videoFrameRate"`
	VideoFrameSize               xmptpg.Dimensions `xmp:"xmpDM:videoFrameSize"`
	VideoPixelDepth              PixelDepth        `xmp:"xmpDM:videoPixelDepth"`
	VideoPixelAspectRatio        xmp.Rational      `xmp:"xmpDM:videoPixelAspectRatio"`
	PartOfCompilation            xmp.Bool          `xmp:"xmpDM:partOfCompilation"`
	Lyrics                       string            `xmp:"xmpDM:lyrics"`
	DiscNumber                   int               `xmp:"xmpDM:discNumber"`

	// not found in XMP spec, but in files
	StartTimeScale      int `xmp:"xmpDM:startTimeScale"`      // "24"
	StartTimeSampleSize int `xmp:"xmpDM:startTimeSampleSize"` // "1"
}

func (x XmpDM) Can(nsName string) bool {
	return NsXmpDM.GetName() == nsName
}

func (x XmpDM) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsXmpDM}
}

func (x *XmpDM) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *XmpDM) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x XmpDM) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *XmpDM) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *XmpDM) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsXmpDM.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *XmpDM) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsXmpDM.GetName(), err)
	}
	return nil
}
