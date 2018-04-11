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

// http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf

package exif

type ColorSpace int

const (
	ColorSpaceSRGB         ColorSpace = 1
	ColorSpaceUncalibrated ColorSpace = 0xffff
)

type Component int

const (
	ComponentNone Component = 0
	ComponentY    Component = 1
	ComponentCb   Component = 2
	ComponentCr   Component = 3
	ComponentR    Component = 4
	ComponentG    Component = 5
	ComponentB    Component = 6
)

type ExposureProgram int

const (
	ExposureProgramUndefined        ExposureProgram = 0
	ExposureProgramManual           ExposureProgram = 1
	ExposureProgramNormalProgram    ExposureProgram = 2
	ExposureProgramAperturePriority ExposureProgram = 3
	ExposureProgramShutterPriority  ExposureProgram = 4
	ExposureProgramCreative         ExposureProgram = 5
	ExposureProgramAction           ExposureProgram = 6
	ExposureProgramPortrait         ExposureProgram = 7
	ExposureProgramLandscape        ExposureProgram = 8
	ExposureProgramBulb             ExposureProgram = 9
)

type SensitivityType int

const (
	SensitivityTypeUnknown     SensitivityType = 0 // Unknown
	SensitivityTypeSOS         SensitivityType = 1 // Standard output sensitivity (SOS)
	SensitivityTypeREI         SensitivityType = 2 // Recommended exposure index (REI)
	SensitivityTypeISO         SensitivityType = 3 // ISOspeed
	SensitivityTypeSOS_REI     SensitivityType = 4 // Standard output sensitivity (SOS) and recommended exposure index (REI)
	SensitivityTypeSOS_ISO     SensitivityType = 5 // Standard output sensitivity (SOS) and ISOspeed
	SensitivityTypeREI_ISO     SensitivityType = 6 // Recommended exposure index (REI) and ISO speed
	SensitivityTypeSOS_REI_ISO SensitivityType = 7 // Standard output sensitivity (SOS) and recommended exposure index (REI) and ISO speed
)

type MeteringMode int

const (
	MeteringModeUnknown               MeteringMode = 0
	MeteringModeAverage               MeteringMode = 1
	MeteringModeCenterWeightedAverage MeteringMode = 2
	MeteringModeSpot                  MeteringMode = 3
	MeteringModeMultiSpot             MeteringMode = 4
	MeteringModePattern               MeteringMode = 5
	MeteringModePartial               MeteringMode = 6
	MeteringModeother                 MeteringMode = 255
)

type LightSource int

const (
	LightSourceUnknown              LightSource = 0
	LightSourceDaylight             LightSource = 1
	LightSourceFluorescent          LightSource = 2
	LightSourceTungsten             LightSource = 3
	LightSourceFlash                LightSource = 4
	LightSourceFineWeather          LightSource = 9
	LightSourceCloudyWeather        LightSource = 10
	LightSourceShade                LightSource = 11
	LightSourceDaylightFluorescent  LightSource = 12 // D 5700 - 7100K
	LightSourceDayWhiteFluorescent  LightSource = 13 // N 4600 - 5500K
	LightSourceCoolWhiteFluorescent LightSource = 14 // W 3800 - 4500K
	LightSourceWhiteFluorescent     LightSource = 15 // WW 3250 - 3800K
	LightSourceWarmWhiteFluorescent LightSource = 16 // L 2600 - 3250K
	LightSourceStandardA            LightSource = 17
	LightSourceStandardB            LightSource = 18
	LightSourceStandardC            LightSource = 19
	LightSourceD55                  LightSource = 20
	LightSourceD65                  LightSource = 21
	LightSourceD75                  LightSource = 22
	LightSourceD50                  LightSource = 23
	LightSourceISOStudioTungsten    LightSource = 24
	LightSourceOther                LightSource = 255
)

type SensingMode int

const (
	SensingModeUndefined SensingMode = 1 // Not defined
	SensingModeOneChip   SensingMode = 2 // One-chip color area sensor
	SensingModeTwoChip   SensingMode = 3 // Two-chip color area sensor
	SensingModeThreeChip SensingMode = 4 // Three-chip color area sensor
	SensingModeArea      SensingMode = 5 // Color sequential area sensor
	SensingModeTrilinear SensingMode = 7 // Trilinear sensor
	SensingModeLinear    SensingMode = 8 // Color sequential linear sensor
)

type FileSourceType int

const (
	FileSourceOther              FileSourceType = 0
	FileSourceScannerTransparent FileSourceType = 1
	FileSourceScannerReflex      FileSourceType = 2
	FileSourceDSC                FileSourceType = 3
)

type RenderMode int

const (
	RenderModeNormal RenderMode = 0
	RenderModeCustom RenderMode = 1
)

type ExposureMode int

const (
	ExposureModeAuto    ExposureMode = 0
	ExposureModeManual  ExposureMode = 1
	ExposureModeBracket ExposureMode = 2
)

type WhiteBalanceMode int

const (
	WhiteBalanceAuto   WhiteBalanceMode = 0
	WhiteBalanceManual WhiteBalanceMode = 1
)

type SceneCaptureType int

const (
	SceneCaptureTypeStandard  SceneCaptureType = 0
	SceneCaptureTypeLandscape SceneCaptureType = 1
	SceneCaptureTypePortrait  SceneCaptureType = 2
	SceneCaptureTypeNight     SceneCaptureType = 3
)

type GainMode int

const (
	GainModeNone     GainMode = 0
	GainModeLowUp    GainMode = 1
	GainModeHighUp   GainMode = 2
	GainModeLowDown  GainMode = 3
	GainModeHighDown GainMode = 4
)

type ContrastMode int

const (
	ContrastModeNormal ContrastMode = 0
	ContrastModeSoft   ContrastMode = 1
	ContrastModeHard   ContrastMode = 2
)

type SaturationMode int

const (
	SaturationModeNormal SaturationMode = 0
	SaturationModeLow    SaturationMode = 1
	SaturationModeHigh   SaturationMode = 2
)

type SharpnessMode int

const (
	SharpnessModeNormal SharpnessMode = 0
	SharpnessModeSoft   SharpnessMode = 1
	SharpnessModeHard   SharpnessMode = 2
)

type SubjectDistanceMode int

const (
	SubjectDistanceUnknown SubjectDistanceMode = 0
	SubjectDistanceMacro   SubjectDistanceMode = 1 // macro mode
	SubjectDistanceClose   SubjectDistanceMode = 2 // 1-3 m
	SubjectDistanceDistant SubjectDistanceMode = 3 // >3m
)

type InteropMode string

const (
	InteropModeR98 InteropMode = "R98"
	InteropModeTHM InteropMode = "THM"
	InteropModeR03 InteropMode = "R03"
)

type FlashMode int

const (
	FlashModeUnknown  FlashMode = 0
	FlashModeFire     FlashMode = 1
	FlashModeSuppress FlashMode = 2
	FlashModeAudo     FlashMode = 3
)

type FlashReturnMode int

const (
	FlashReturnModeDisabled FlashReturnMode = 0 // no strobe return detection
	FlashReturnModeNo       FlashReturnMode = 2 // strobe return light not detected
	FlashReturnModeYes      FlashReturnMode = 3 // strobe return light detected
)
