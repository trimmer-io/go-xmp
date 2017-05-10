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

package xmpdm

// Part 1: 8.2.2.8 RenditionClass
type RenditionClass string

const (
	RenditionClassInvalid   RenditionClass = ""
	RenditionClassDefault   RenditionClass = "default"
	RenditionClassDraft     RenditionClass = "draft"
	RenditionClassLowRes    RenditionClass = "low-res"
	RenditionClassProof     RenditionClass = "proof"
	RenditionClassScreen    RenditionClass = "screen"
	RenditionClassThumbnail RenditionClass = "thumbnail"
)

type AudioChannelType string

const (
	AudioChannelTypeInvalid AudioChannelType = ""
	AudioChannelTypeMono    AudioChannelType = "Mono"
	AudioChannelTypeStereo  AudioChannelType = "Stereo"
	AudioChannelType5_1     AudioChannelType = "5.1"
	AudioChannelType7_1     AudioChannelType = "7.1"
	AudioChannelType16Ch    AudioChannelType = "16 Channel"
	AudioChannelTypeOther   AudioChannelType = "Other"
)

type AudioSampleType string

const (
	AudioSampleTypeInvalid    AudioSampleType = ""
	AudioSampleType8i         AudioSampleType = "8Int"
	AudioSampleType16i        AudioSampleType = "16Int"
	AudioSampleType24i        AudioSampleType = "24Int"
	AudioSampleType32i        AudioSampleType = "32Int"
	AudioSampleType32f        AudioSampleType = "32Float"
	AudioSampleTypeCompressed AudioSampleType = "Compressed"
	AudioSampleTypePacked     AudioSampleType = "Packed"
	AudioSampleTypeOther      AudioSampleType = "Other"
)

type CameraAngle string

const (
	CameraAngleInvalid         CameraAngle = ""
	CameraAngleLowAngle        CameraAngle = "Low Angle"
	CameraAngleEyeLevel        CameraAngle = "Eye Level"
	CameraAngleHighAngle       CameraAngle = "High Angle"
	CameraAngleOverheadShot    CameraAngle = "Overhead Shot"
	CameraAngleBirdsEyeShot    CameraAngle = "Birds Eye Shot"
	CameraAngleDutchAngle      CameraAngle = "Dutch Angle"
	CameraAnglePOV             CameraAngle = "POV"
	CameraAngleOverTheShoulder CameraAngle = "Over the Shoulder"
	CameraAngleReactionShot    CameraAngle = "Reaction Shot"
)

type CameraMove string

const (
	CameraMoveInvalid      CameraMove = ""
	CameraMoveAerial       CameraMove = "Aerial"
	CameraMoveBoomUp       CameraMove = "Boom Up"
	CameraMoveBoomDown     CameraMove = "Boom Down"
	CameraMoveCraneUp      CameraMove = "Crane Up"
	CameraMoveCraneDown    CameraMove = "Crane Down"
	CameraMoveDollyIn      CameraMove = "Dolly In"
	CameraMoveDollyOut     CameraMove = "Dolly Out"
	CameraMovePanLeft      CameraMove = "Pan Left"
	CameraMovePanRight     CameraMove = "Pan Right"
	CameraMovePedestalUp   CameraMove = "Pedestal Up"
	CameraMovePedestalDown CameraMove = "Pedestal Down"
	CameraMoveTiltUp       CameraMove = "Tilt Up"
	CameraMoveTiltDown     CameraMove = "Tilt Down"
	CameraMoveTracking     CameraMove = "Tracking"
	CameraMoveTruckLeft    CameraMove = "Truck Left"
	CameraMoveTruckRight   CameraMove = "Truck Right"
	CameraMoveZoomIn       CameraMove = "Zoom In"
	CameraMoveZoomOut      CameraMove = "Zoom Out"
)

type MusicalKey string

const (
	MusicalKeyInvalid MusicalKey = ""
	MusicalKeyC       MusicalKey = "C"
	MusicalKeyC_      MusicalKey = "C#"
	MusicalKeyD       MusicalKey = "D"
	MusicalKeyD_      MusicalKey = "D#"
	MusicalKeyE       MusicalKey = "E"
	MusicalKeyF       MusicalKey = "F"
	MusicalKeyF_      MusicalKey = "F#"
	MusicalKeyG       MusicalKey = "G"
	MusicalKeyG_      MusicalKey = "G#"
	MusicalKeyA       MusicalKey = "A"
	MusicalKeyA_      MusicalKey = "A#"
	MusicalKeyB       MusicalKey = "B"
)

type MusicalScale string

const (
	MusicalScaleInvalid MusicalScale = ""
	MusicalScaleMajor   MusicalScale = "Major"
	MusicalScaleMinor   MusicalScale = "Minor"
	MusicalScaleBoth    MusicalScale = "Both"
	MusicalScaleNeither MusicalScale = "Neither"
)

type PullDown string

const (
	PullDownInvalid PullDown = ""
	PullDownWSSWW   PullDown = "WSSWW"
	PullDownSSWWW   PullDown = "SSWWW"
	PullDownSWWWS   PullDown = "SWWWS"
	PullDownWWWSS   PullDown = "WWWSS"
	PullDownWWSSW   PullDown = "WWSSW"
	PullDownWWWSW   PullDown = "WWWSW"
	PullDownWWSWW   PullDown = "WWSWW"
	PullDownWSWWW   PullDown = "WSWWW"
	PullDownSWWWW   PullDown = "SWWWW"
	PullDownWWWWS   PullDown = "WWWWS"
)

type ShotSize string

const (
	ShotSizeInvalid ShotSize = ""    // implementation internal
	ShotSizeECU     ShotSize = "ECU" // extreme close-up
	ShotSizeMCU     ShotSize = "MCU" //  medium close-up
	ShotSizeCU      ShotSize = "CU"  //  close-up
	ShotSizeMS      ShotSize = "MS"  //  medium shot
	ShotSizeWS      ShotSize = "WS"  //  wide shot
	ShotSizeMWS     ShotSize = "MWS" //  medium wide shot
	ShotSizeEWS     ShotSize = "EWS" //  extreme wide shot
)

type StretchMode string

const (
	StretchModeInvalid     StretchMode = ""
	StretchModeFixedLength StretchMode = "Fixed length"
	StretchModeTimeScale   StretchMode = "Time-Scale"
	StretchModeResample    StretchMode = "Resample"
	StretchModeBeatSplice  StretchMode = "Beat Splice"
	StretchModeHybrid      StretchMode = "Hybrid"
)

type TimeSignature string

const (
	TimeSignatureInvalid TimeSignature = ""
	TimeSignature2_4     TimeSignature = "2/4"
	TimeSignature3_4     TimeSignature = "3/4"
	TimeSignature4_4     TimeSignature = "4/4"
	TimeSignature5_4     TimeSignature = "5/4"
	TimeSignature7_4     TimeSignature = "7/4"
	TimeSignature6_8     TimeSignature = "6/8"
	TimeSignature9_8     TimeSignature = "9/8"
	TimeSignature12_8    TimeSignature = "12/8"
)

type AlphaMode string

const (
	AlphaModeInvalid       AlphaMode = ""
	AlphaModeStraight      AlphaMode = "straight"
	AlphaModePreMultiplied AlphaMode = "pre-multiplied"
	AlphaModeNone          AlphaMode = "none"
)

type ColorSpace string

const (
	ColorSpaceInvalid ColorSpace = ""
	ColorSpaceSRGB    ColorSpace = "sRGB"
	ColorSpaceCCIR601 ColorSpace = "CCIR-601"
	ColorSpaceCCIR709 ColorSpace = "CCIR-709"
)

type FieldOrder string

const (
	FieldOrderInvalid     FieldOrder = ""
	FieldOrderUpper       FieldOrder = "Upper"
	FieldOrderLower       FieldOrder = "Lower"
	FieldOrderProgressive FieldOrder = "Progressive"
)

type PixelDepth string

const (
	PixelDepthInvalid PixelDepth = ""
	PixelDepth8       PixelDepth = "8Int"
	PixelDepth16      PixelDepth = "16Int"
	PixelDepth24      PixelDepth = "24Int"
	PixelDepth32      PixelDepth = "32Int"
	PixelDepth32f     PixelDepth = "32Float"
	PixelDepthOther   PixelDepth = "Other"
)

type TimecodeFormat string

const (
	TimecodeFormatInvalid TimecodeFormat = ""
	TimecodeFormat24      TimecodeFormat = "24Timecode"
	TimecodeFormat25      TimecodeFormat = "25Timecode"
	TimecodeFormat2997    TimecodeFormat = "2997DropTimecode" // (semicolon delimiter)
	TimecodeFormat2997ND  TimecodeFormat = "2997NonDropTimecode"
	TimecodeFormat30      TimecodeFormat = "30Timecode"
	TimecodeFormat50      TimecodeFormat = "50Timecode"
	TimecodeFormat5994    TimecodeFormat = "5994DropTimecode" // (semicolon delimiter)
	TimecodeFormat5994ND  TimecodeFormat = "5994NonDropTimecode"
	TimecodeFormat60      TimecodeFormat = "60Timecode"
	TimecodeFormat23976   TimecodeFormat = "23976Timecode"
)

type Quality string

const (
	QualityInvalid Quality = ""
	QualityHigh    Quality = "High"
	QualityMedium  Quality = "Medium"
	QualityLow     Quality = "Low"
)

type FileType string

const (
	FileTypeInvalid FileType = ""
	FileTypeMovie   FileType = "movie"
	FileTypeStill   FileType = "still"
	FileTypeAudio   FileType = "audio"
	FileTypeCustom  FileType = "custom"
)

type MarkerType string

const (
	MarkerTypeInvalid      MarkerType = ""                            // implementation internal
	MarkerTypeChapter      MarkerType = "Chapter"                     // XMP DM
	MarkerTypeCue          MarkerType = "Cue"                         // XMP DM
	MarkerTypeIndex        MarkerType = "Index"                       // XMP DM
	MarkerTypeSpeech       MarkerType = "Speech"                      // XMP DM
	MarkerTypeTrack        MarkerType = "Track"                       // XMP DM
	MarkerTypeIptcDuration MarkerType = "ivqu:editorialDuration"      // IPTC Video Metadata Nov 2016
	MarkerTypeIptcStart    MarkerType = "ivqu:editorialDurationStart" // IPTC Video Metadata Nov 2016
	MarkerTypeIptcEnd      MarkerType = "ivqu:editorialDurationEnd"   // IPTC Video Metadata Nov 2016
)
