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

package tiff

// TIFF enumeration values
// http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf

type CompressionType int

const (
	CompressionUncompressed CompressionType = 1
	CompressionJPEG         CompressionType = 6
)

type ColorModel int

const (
	ColorModelWhiteIsZero      ColorModel = 0     // WhiteIsZero
	ColorModelBlackIsZero      ColorModel = 1     // BlackIsZero
	ColorModelRGB              ColorModel = 2     // RGB
	ColorModelRGBPalette       ColorModel = 3     // RGB Palette
	ColorModelTransparencyMask ColorModel = 4     // Transparency Mask
	ColorModelCMYK             ColorModel = 5     // CMYK
	ColorModelYCbCr            ColorModel = 6     // YCbCr
	ColorModelCIELab           ColorModel = 8     // CIELab
	ColorModelICCLab           ColorModel = 9     // ICCLab
	ColorModelITULab           ColorModel = 10    // ITULab
	ColorModelColorFilterArray ColorModel = 32803 // Color Filter Array
	ColorModelPixarLogL        ColorModel = 32844 // Pixar LogL
	ColorModelPixarLogLuv      ColorModel = 32845 // Pixar LogLuv
	ColorModelLinearRaw        ColorModel = 34892 // Linear Raw
)

type OrientationType int

const (
	OrientationTopLeft     OrientationType = 1 //  1 = Horizontal (normal)
	OrientationTopRight    OrientationType = 2 //  2 = Mirror horizontal
	OrientationBottomRight OrientationType = 3 //  3 = Rotate 180
	OrientationBottomLeft  OrientationType = 4 //  4 = Mirror vertical
	OrientationLeftTop     OrientationType = 5 //  5 = Mirror horizontal and rotate 270 CW
	OrientationRightTop    OrientationType = 6 //  6 = Rotate 90 CW
	OrientationRightBottom OrientationType = 7 //  7 = Mirror horizontal and rotate 90 CW
	OrientationLeftBottom  OrientationType = 8 //  8 = Rotate 270 CW
)

type PlanarType int

const (
	PlanarChunky PlanarType = 1 // 1 = chunky
	PlanarPlanar PlanarType = 2 // 2 = planar
)

type YCbCrPosition int

const (
	YCbCrPositionCentered YCbCrPosition = 1
	YCbCrPositionCoSited  YCbCrPosition = 2
)

type ResolutionUnit int

const (
	ResolutionUnitInches     ResolutionUnit = 2
	ResolutionUnitCentimeter ResolutionUnit = 3
)
