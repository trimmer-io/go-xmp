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

// Package crs implements Adobe Camera Raw metadata as defined in chapter 3.3 of XMP Specification Part 2.
package crs

import (
	"fmt"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsCrs = xmp.NewNamespace("crs", "http://ns.adobe.com/camera-raw-settings/1.0/", NewModel)
)

func init() {
	xmp.Register(NsCrs, xmp.ImageMetadata)
}

func NewModel(name string) xmp.Model {
	return &CameraRawInfo{}
}

func MakeModel(d *xmp.Document) (*CameraRawInfo, error) {
	m, err := d.MakeModel(NsCrs)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*CameraRawInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *CameraRawInfo {
	if m := d.FindModel(NsCrs); m != nil {
		return m.(*CameraRawInfo)
	}
	return nil
}

// Part 2: 3.3 Camera Raw namespace
type CameraRawInfo struct {
	AutoBrightness       xmp.Bool        `xmp:"crs:AutoBrightness"`
	AutoContrast         xmp.Bool        `xmp:"crs:AutoContrast"`
	AutoExposure         xmp.Bool        `xmp:"crs:AutoExposure"`
	AutoShadows          xmp.Bool        `xmp:"crs:AutoShadows"`
	BlueHue              int64           `xmp:"crs:BlueHue"`
	BlueSaturation       int64           `xmp:"crs:BlueSaturation"`
	Brightness           int64           `xmp:"crs:Brightness"`
	CameraProfile        string          `xmp:"crs:CameraProfile"`
	ChromaticAberrationB int64           `xmp:"crs:ChromaticAberrationB"`
	ChromaticAberrationR int64           `xmp:"crs:ChromaticAberrationR"`
	ColorNoiseReduction  int64           `xmp:"crs:ColorNoiseReduction"`
	Contrast             int64           `xmp:"crs:Contrast"`
	CropTop              float32         `xmp:"crs:CropTop"`
	CropLeft             float32         `xmp:"crs:CropLeft"`
	CropBottom           float32         `xmp:"crs:CropBottom"`
	CropRight            float32         `xmp:"crs:CropRight"`
	CropAngle            float32         `xmp:"crs:CropAngle"`
	CropWidth            float32         `xmp:"crs:CropWidth"`
	CropHeight           float32         `xmp:"crs:CropHeight"`
	CropUnits            CropUnits       `xmp:"crs:CropUnits"`
	Exposure             float32         `xmp:"crs:Exposure"`
	GreenHue             int64           `xmp:"crs:GreenHue"`
	GreenSaturation      int64           `xmp:"crs:GreenSaturation"`
	HasCrop              xmp.Bool        `xmp:"crs:HasCrop"`
	HasSettings          xmp.Bool        `xmp:"crs:HasSettings"`
	LuminanceSmoothing   int64           `xmp:"crs:LuminanceSmoothing"`
	RawFileName          string          `xmp:"crs:RawFileName"`
	RedHue               int64           `xmp:"crs:RedHue"`
	RedSaturation        int64           `xmp:"crs:RedSaturation"`
	Saturation           int64           `xmp:"crs:Saturation"`
	Shadows              int64           `xmp:"crs:Shadows"`
	ShadowTint           int64           `xmp:"crs:ShadowTint"`
	Sharpness            int64           `xmp:"crs:Sharpness"`
	Temperature          int64           `xmp:"crs:Temperature"`
	Tint                 int64           `xmp:"crs:Tint"`
	ToneCurve            xmp.StringArray `xmp:"crs:ToneCurve"`
	ToneCurveName        ToneCurve       `xmp:"crs:ToneCurveName"`
	Version              string          `xmp:"crs:Version"`
	VignetteAmount       int64           `xmp:"crs:VignetteAmount"`
	VignetteMidpoint     int64           `xmp:"crs:VignetteMidpoint"`
	WhiteBalance         WhiteBalance    `xmp:"crs:WhiteBalance"`

	// non-standard fields found in samples in the wild
	AlreadyApplied                        xmp.Bool        `xmp:"crs:AlreadyApplied"`
	Dehaze                                int             `xmp:"crs:Dehaze"`
	Converter                             string          `xmp:"crs:Converter"`
	MoireFilter                           string          `xmp:"crs:MoireFilter"`
	Smoothness                            int64           `xmp:"crs:Smoothness"`
	CameraProfileDigest                   string          `xmp:"crs:CameraProfileDigest"`
	Clarity                               int64           `xmp:"crs:Clarity"`
	ConvertToGrayscale                    xmp.Bool        `xmp:"crs:ConvertToGrayscale"`
	Defringe                              int64           `xmp:"crs:Defringe"`
	FillLight                             int64           `xmp:"crs:FillLight"`
	HighlightRecovery                     int64           `xmp:"crs:HighlightRecovery"`
	HueAdjustmentAqua                     int64           `xmp:"crs:HueAdjustmentAqua"`
	HueAdjustmentBlue                     int64           `xmp:"crs:HueAdjustmentBlue"`
	HueAdjustmentGreen                    int64           `xmp:"crs:HueAdjustmentGreen"`
	HueAdjustmentMagenta                  int64           `xmp:"crs:HueAdjustmentMagenta"`
	HueAdjustmentOrange                   int64           `xmp:"crs:HueAdjustmentOrange"`
	HueAdjustmentPurple                   int64           `xmp:"crs:HueAdjustmentPurple"`
	HueAdjustmentRed                      int64           `xmp:"crs:HueAdjustmentRed"`
	HueAdjustmentYellow                   int64           `xmp:"crs:HueAdjustmentYellow"`
	IncrementalTemperature                int64           `xmp:"crs:IncrementalTemperature"`
	IncrementalTint                       int64           `xmp:"crs:IncrementalTint"`
	LuminanceAdjustmentAqua               int64           `xmp:"crs:LuminanceAdjustmentAqua"`
	LuminanceAdjustmentBlue               int64           `xmp:"crs:LuminanceAdjustmentBlue"`
	LuminanceAdjustmentGreen              int64           `xmp:"crs:LuminanceAdjustmentGreen"`
	LuminanceAdjustmentMagenta            int64           `xmp:"crs:LuminanceAdjustmentMagenta"`
	LuminanceAdjustmentOrange             int64           `xmp:"crs:LuminanceAdjustmentOrange"`
	LuminanceAdjustmentPurple             int64           `xmp:"crs:LuminanceAdjustmentPurple"`
	LuminanceAdjustmentRed                int64           `xmp:"crs:LuminanceAdjustmentRed"`
	LuminanceAdjustmentYellow             int64           `xmp:"crs:LuminanceAdjustmentYellow"`
	ParametricDarks                       int64           `xmp:"crs:ParametricDarks"`
	ParametricHighlights                  int64           `xmp:"crs:ParametricHighlights"`
	ParametricHighlightSplit              int64           `xmp:"crs:ParametricHighlightSplit"`
	ParametricLights                      int64           `xmp:"crs:ParametricLights"`
	ParametricMidtoneSplit                int64           `xmp:"crs:ParametricMidtoneSplit"`
	ParametricShadows                     int64           `xmp:"crs:ParametricShadows"`
	ParametricShadowSplit                 int64           `xmp:"crs:ParametricShadowSplit"`
	SaturationAdjustmentAqua              int64           `xmp:"crs:SaturationAdjustmentAqua"`
	SaturationAdjustmentBlue              int64           `xmp:"crs:SaturationAdjustmentBlue"`
	SaturationAdjustmentGreen             int64           `xmp:"crs:SaturationAdjustmentGreen"`
	SaturationAdjustmentMagenta           int64           `xmp:"crs:SaturationAdjustmentMagenta"`
	SaturationAdjustmentOrange            int64           `xmp:"crs:SaturationAdjustmentOrange"`
	SaturationAdjustmentPurple            int64           `xmp:"crs:SaturationAdjustmentPurple"`
	SaturationAdjustmentRed               int64           `xmp:"crs:SaturationAdjustmentRed"`
	SaturationAdjustmentYellow            int64           `xmp:"crs:SaturationAdjustmentYellow"`
	SharpenDetail                         int64           `xmp:"crs:SharpenDetail"`
	SharpenEdgeMasking                    int64           `xmp:"crs:SharpenEdgeMasking"`
	SharpenRadius                         float32         `xmp:"crs:SharpenRadius"`
	SplitToningBalance                    int64           `xmp:"crs:SplitToningBalance"`
	SplitToningHighlightHue               int64           `xmp:"crs:SplitToningHighlightHue"`
	SplitToningHighlightSaturation        int64           `xmp:"crs:SplitToningHighlightSaturation"`
	SplitToningShadowHue                  int64           `xmp:"crs:SplitToningShadowHue"`
	SplitToningShadowSaturation           int64           `xmp:"crs:SplitToningShadowSaturation"`
	Vibrance                              int64           `xmp:"crs:Vibrance"`
	GrayMixerRed                          int64           `xmp:"crs:GrayMixerRed"`
	GrayMixerOrange                       int64           `xmp:"crs:GrayMixerOrange"`
	GrayMixerYellow                       int64           `xmp:"crs:GrayMixerYellow"`
	GrayMixerGreen                        int64           `xmp:"crs:GrayMixerGreen"`
	GrayMixerAqua                         int64           `xmp:"crs:GrayMixerAqua"`
	GrayMixerBlue                         int64           `xmp:"crs:GrayMixerBlue"`
	GrayMixerPurple                       int64           `xmp:"crs:GrayMixerPurple"`
	GrayMixerMagenta                      int64           `xmp:"crs:GrayMixerMagenta"`
	RetouchInfo                           xmp.StringArray `xmp:"crs:RetouchInfo"`
	RedEyeInfo                            xmp.StringArray `xmp:"crs:RedEyeInfo"`
	CropUnit                              int64           `xmp:"crs:CropUnit"`
	PostCropVignetteAmount                int64           `xmp:"crs:PostCropVignetteAmount"`
	PostCropVignetteMidpoint              int64           `xmp:"crs:PostCropVignetteMidpoint"`
	PostCropVignetteFeather               int64           `xmp:"crs:PostCropVignetteFeather"`
	PostCropVignetteRoundness             int64           `xmp:"crs:PostCropVignetteRoundness"`
	PostCropVignetteStyle                 int64           `xmp:"crs:PostCropVignetteStyle"`
	ProcessVersion                        string          `xmp:"crs:ProcessVersion"`
	LensProfileEnable                     int64           `xmp:"crs:LensProfileEnable"`
	LensProfileSetup                      string          `xmp:"crs:LensProfileSetup"`
	LensProfileName                       string          `xmp:"crs:LensProfileName"`
	LensProfileFilename                   string          `xmp:"crs:LensProfileFilename"`
	LensProfileDigest                     string          `xmp:"crs:LensProfileDigest"`
	LensProfileDistortionScale            int64           `xmp:"crs:LensProfileDistortionScale"`
	LensProfileChromaticAberrationScale   int64           `xmp:"crs:LensProfileChromaticAberrationScale"`
	LensProfileVignettingScale            int64           `xmp:"crs:LensProfileVignettingScale"`
	LensManualDistortionAmount            int64           `xmp:"crs:LensManualDistortionAmount"`
	PerspectiveVertical                   int64           `xmp:"crs:PerspectiveVertical"`
	PerspectiveHorizontal                 int64           `xmp:"crs:PerspectiveHorizontal"`
	PerspectiveRotate                     float32         `xmp:"crs:PerspectiveRotate"`
	PerspectiveScale                      int64           `xmp:"crs:PerspectiveScale"`
	CropConstrainToWarp                   int64           `xmp:"crs:CropConstrainToWarp"`
	LuminanceNoiseReductionDetail         int64           `xmp:"crs:LuminanceNoiseReductionDetail"`
	LuminanceNoiseReductionContrast       int64           `xmp:"crs:LuminanceNoiseReductionContrast"`
	ColorNoiseReductionDetail             int64           `xmp:"crs:ColorNoiseReductionDetail"`
	GrainAmount                           int64           `xmp:"crs:GrainAmount"`
	GrainSize                             int64           `xmp:"crs:GrainSize"`
	GrainFrequency                        int64           `xmp:"crs:GrainFrequency"`
	AutoLateralCA                         int64           `xmp:"crs:AutoLateralCA"`
	Exposure2012                          float32         `xmp:"crs:Exposure2012"`
	Contrast2012                          int64           `xmp:"crs:Contrast2012"`
	Highlights2012                        int64           `xmp:"crs:Highlights2012"`
	Shadows2012                           int64           `xmp:"crs:Shadows2012"`
	Whites2012                            int64           `xmp:"crs:Whites2012"`
	Blacks2012                            int64           `xmp:"crs:Blacks2012"`
	Clarity2012                           int64           `xmp:"crs:Clarity2012"`
	PostCropVignetteHighlightContrast     int64           `xmp:"crs:PostCropVignetteHighlightContrast"`
	ToneCurveName2012                     ToneCurve       `xmp:"crs:ToneCurveName2012"`
	ToneCurveRed                          PointList       `xmp:"crs:ToneCurveRed"`
	ToneCurveGreen                        PointList       `xmp:"crs:ToneCurveGreen"`
	ToneCurveBlue                         PointList       `xmp:"crs:ToneCurveBlue"`
	ToneCurvePV2012                       PointList       `xmp:"crs:ToneCurvePV2012"`
	ToneCurvePV2012Red                    PointList       `xmp:"crs:ToneCurvePV2012Red"`
	ToneCurvePV2012Green                  PointList       `xmp:"crs:ToneCurvePV2012Green"`
	ToneCurvePV2012Blue                   PointList       `xmp:"crs:ToneCurvePV2012Blue"`
	ToneMapStrength                       int             `xmp:"crs:ToneMapStrength"`
	DefringePurpleAmount                  int64           `xmp:"crs:DefringePurpleAmount"`
	DefringePurpleHueLo                   int64           `xmp:"crs:DefringePurpleHueLo"`
	DefringePurpleHueHi                   int64           `xmp:"crs:DefringePurpleHueHi"`
	DefringeGreenAmount                   int64           `xmp:"crs:DefringeGreenAmount"`
	DefringeGreenHueLo                    int64           `xmp:"crs:DefringeGreenHueLo"`
	DefringeGreenHueHi                    int64           `xmp:"crs:DefringeGreenHueHi"`
	AutoWhiteVersion                      int64           `xmp:"crs:AutoWhiteVersion"`
	ColorNoiseReductionSmoothness         int64           `xmp:"crs:ColorNoiseReductionSmoothness"`
	PerspectiveAspect                     int64           `xmp:"crs:PerspectiveAspect"`
	PerspectiveUpright                    int64           `xmp:"crs:PerspectiveUpright"`
	UprightVersion                        int64           `xmp:"crs:UprightVersion"`
	UprightCenterMode                     int64           `xmp:"crs:UprightCenterMode"`
	UprightCenterNormX                    float32         `xmp:"crs:UprightCenterNormX"`
	UprightCenterNormY                    float32         `xmp:"crs:UprightCenterNormY"`
	UprightFocalMode                      int64           `xmp:"crs:UprightFocalMode"`
	UprightFocalLength35mm                float32         `xmp:"crs:UprightFocalLength35mm"`
	UprightPreview                        xmp.Bool        `xmp:"crs:UprightPreview"`
	UprightTransformCount                 int64           `xmp:"crs:UprightTransformCount"`
	UprightDependentDigest                string          `xmp:"crs:UprightDependentDigest"`
	UprightTransform_0                    string          `xmp:"crs:UprightTransform_0"`
	UprightTransform_1                    string          `xmp:"crs:UprightTransform_1"`
	UprightTransform_2                    string          `xmp:"crs:UprightTransform_2"`
	UprightTransform_3                    string          `xmp:"crs:UprightTransform_3"`
	UprightTransform_4                    string          `xmp:"crs:UprightTransform_4"`
	LensProfileMatchKeyExifMake           string          `xmp:"crs:LensProfileMatchKeyExifMake"`
	LensProfileMatchKeyExifModel          string          `xmp:"crs:LensProfileMatchKeyExifModel"`
	LensProfileMatchKeyCameraModelName    string          `xmp:"crs:LensProfileMatchKeyCameraModelName"`
	LensProfileMatchKeyLensInfo           string          `xmp:"crs:LensProfileMatchKeyLensInfo"`
	LensProfileMatchKeyLensID             string          `xmp:"crs:LensProfileMatchKeyLensID"`
	LensProfileMatchKeyLensName           string          `xmp:"crs:LensProfileMatchKeyLensName"`
	LensProfileMatchKeyIsRaw              xmp.Bool        `xmp:"crs:LensProfileMatchKeyIsRaw"`
	LensProfileMatchKeySensorFormatFactor float32         `xmp:"crs:LensProfileMatchKeySensorFormatFactor"`
	DefaultAutoTone                       xmp.Bool        `xmp:"crs:DefaultAutoTone"`
	DefaultAutoGray                       xmp.Bool        `xmp:"crs:DefaultAutoGray"`
	DefaultsSpecificToSerial              xmp.Bool        `xmp:"crs:DefaultsSpecificToSerial"`
	DefaultsSpecificToISO                 xmp.Bool        `xmp:"crs:DefaultsSpecificToISO"`
	DNGIgnoreSidecars                     xmp.Bool        `xmp:"crs:DNGIgnoreSidecars"`
	NegativeCachePath                     string          `xmp:"crs:NegativeCachePath"`
	NegativeCacheMaximumSize              float32         `xmp:"crs:NegativeCacheMaximumSize"`
	NegativeCacheLargePreviewSize         int64           `xmp:"crs:NegativeCacheLargePreviewSize"`
	JPEGHandling                          string          `xmp:"crs:JPEGHandling"`
	TIFFHandling                          string          `xmp:"crs:TIFFHandling"`
	GradientBasedCorrections              CorrectionList  `xmp:"crs:GradientBasedCorrections"`
	CircularGradientBasedCorrections      CorrectionList  `xmp:"crs:CircularGradientBasedCorrections"`
	PaintBasedCorrections                 CorrectionList  `xmp:"crs:PaintBasedCorrections"`
	RetouchAreas                          RetouchAreaList `xmp:"crs:RetouchAreas"`
}

func (x CameraRawInfo) Can(nsName string) bool {
	return NsCrs.GetName() == nsName
}

func (x CameraRawInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsCrs}
}

func (x *CameraRawInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *CameraRawInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x CameraRawInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *CameraRawInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *CameraRawInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsCrs.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *CameraRawInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsCrs.GetName(), err)
	}
	return nil
}
