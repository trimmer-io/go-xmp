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

package crs

type CropUnits int

const (
	CropUnitsPixels = iota // 0
	CropUnitsInches        // 1
	CropUnitsCm            // 2
)

type ToneCurve string

const (
	ToneCurveLinear         ToneCurve = "Linear"
	ToneCurveMediumContrast ToneCurve = "Medium Contrast"
	ToneCurveStrongContrast ToneCurve = "Strong Contrast"
	ToneCurveCustom         ToneCurve = "Custom"
)

type WhiteBalance string

const (
	WhiteBalanceAsShot      WhiteBalance = "As Shot"
	WhiteBalanceAuto        WhiteBalance = "Auto"
	WhiteBalanceDaylight    WhiteBalance = "Daylight"
	WhiteBalanceCloudy      WhiteBalance = "Cloudy"
	WhiteBalanceShade       WhiteBalance = "Shade"
	WhiteBalanceTungsten    WhiteBalance = "Tungsten"
	WhiteBalanceFluorescent WhiteBalance = "Fluorescent"
	WhiteBalanceFlash       WhiteBalance = "Flash"
	WhiteBalanceCustom      WhiteBalance = "Custom"
)
