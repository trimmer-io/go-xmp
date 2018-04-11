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

package ixml

import (
	"fmt"
	"strings"
)

type TakeType string

const (
	TakeTypeDefault      TakeType = "DEFAULT"
	TakeTypeNoGood       TakeType = "NO_GOOD"
	TakeTypeFalseStart   TakeType = "FALSE_START"
	TakeTypeWildTrack    TakeType = "WILD_TRACK"
	TakeTypePickup       TakeType = "PICKUP"
	TakeTypeRehearsal    TakeType = "REHEARSAL"
	TakeTypeAnnouncement TakeType = "ANNOUNCEMENT"
	TakeTypeSoundGuide   TakeType = "SOUND_GUIDE"
)

type TakeTypeList []TakeType

func (x TakeTypeList) Contains(t TakeType) bool {
	for _, v := range x {
		if v == t {
			return true
		}
	}
	return false
}

func (x *TakeTypeList) Add(t TakeType) {
	if !x.Contains(t) {
		*x = append(*x, t)
	}
}

func (x *TakeTypeList) Del(t TakeType) {
	idx := -1
	for i, v := range *x {
		if v == t {
			idx = i
		}
	}
	if idx > -1 {
		*x = append((*x)[:idx], (*x)[idx+1:]...)
	}
}

func (x TakeTypeList) IsDefault() bool {
	if len(x) == 0 {
		return true
	}
	return x.Contains(TakeTypeDefault)
}

// allow future extensions
func ParseTakeType(s string) TakeType {
	switch s {
	case "DEFAULT":
		return TakeTypeDefault
	case "NO_GOOD":
		return TakeTypeNoGood
	case "FALSE_START":
		return TakeTypeFalseStart
	case "WILD_TRACK":
		return TakeTypeWildTrack
	case "PICKUP":
		return TakeTypePickup
	case "REHEARSAL":
		return TakeTypeRehearsal
	case "ANNOUNCEMENT":
		return TakeTypeAnnouncement
	case "SOUND_GUIDE":
		return TakeTypeSoundGuide
	default:
		return TakeType(s)
	}
}

func (x TakeTypeList) MarshalText() ([]byte, error) {
	if len(x) == 0 {
		return nil, nil
	}
	l := make([]string, len(x))
	for i, v := range x {
		l[i] = string(v)
	}
	return []byte(strings.Join(l, ",")), nil
}

func (x *TakeTypeList) UnmarshalText(data []byte) error {
	fields := strings.Split(string(data), ",")
	if len(fields) == 0 {
		return nil
	}
	ttl := make(TakeTypeList, len(fields))
	for i, v := range fields {
		ttl[i] = ParseTakeType(v)
	}
	*x = ttl
	return nil
}

// Location Type
//
type LocationType string

const (
	LocationTypeInterior LocationType = "INT"
	LocationTypeExterior LocationType = "EXT"
)

type LocationTypeList []LocationType

// allow future extensions
func ParseLocationType(s string) LocationType {
	switch s {
	case "INT":
		return LocationTypeInterior
	case "EXT":
		return LocationTypeExterior
	default:
		return LocationType(s)
	}
}

func (x *LocationTypeList) Add(t LocationType) {
	if !x.Contains(t) {
		*x = append(*x, t)
	}
}

func (x *LocationTypeList) Del(t LocationType) {
	idx := -1
	for i, v := range *x {
		if v == t {
			idx = i
		}
	}
	if idx > -1 {
		*x = append((*x)[:idx], (*x)[idx+1:]...)
	}
}

func (x LocationTypeList) MarshalText() ([]byte, error) {
	if len(x) == 0 {
		return nil, nil
	}
	l := make([]string, len(x))
	for i, v := range x {
		l[i] = string(v)
	}
	return []byte(strings.Join(l, ",")), nil
}

func (x *LocationTypeList) UnmarshalText(data []byte) error {
	fields := strings.Split(string(data), ",")
	if len(fields) == 0 {
		return nil
	}
	ttl := make(LocationTypeList, len(fields))
	for i, v := range fields {
		ttl[i] = ParseLocationType(v)
	}
	*x = ttl
	return nil
}

func (x LocationTypeList) Contains(t LocationType) bool {
	for _, v := range x {
		if v == t {
			return true
		}
	}
	return false
}

// Location Time
//
type LocationTime string

const (
	LocationTimeSunrise   LocationTime = "SUNRISE"
	LocationTimeMorning   LocationTime = "MORNING"
	LocationTimeMidday    LocationTime = "MIDDAY"
	LocationTimeDay       LocationTime = "DAY"
	LocationTimeAfternoon LocationTime = "AFTERNOON"
	LocationTimeEvening   LocationTime = "EVENING"
	LocationTimeSunset    LocationTime = "SUNSET"
	LocationTimeNight     LocationTime = "NIGHT"
)

type LocationTimeList []LocationTime

// allow future extensions
func ParseLocationTime(s string) LocationTime {
	switch s {
	case "SUNRISE":
		return LocationTimeSunrise
	case "MORNING":
		return LocationTimeMorning
	case "MIDDAY":
		return LocationTimeMidday
	case "DAY":
		return LocationTimeDay
	case "AFTERNOON":
		return LocationTimeAfternoon
	case "EVENING":
		return LocationTimeEvening
	case "SUNSET":
		return LocationTimeSunset
	case "NIGHT":
		return LocationTimeNight
	default:
		return LocationTime(s)
	}
}

func (x *LocationTimeList) Add(t LocationTime) {
	if !x.Contains(t) {
		*x = append(*x, t)
	}
}

func (x *LocationTimeList) Del(t LocationTime) {
	idx := -1
	for i, v := range *x {
		if v == t {
			idx = i
		}
	}
	if idx > -1 {
		*x = append((*x)[:idx], (*x)[idx+1:]...)
	}
}

func (x LocationTimeList) MarshalText() ([]byte, error) {
	if len(x) == 0 {
		return nil, nil
	}
	l := make([]string, len(x))
	for i, v := range x {
		l[i] = string(v)
	}
	return []byte(strings.Join(l, ",")), nil
}

func (x *LocationTimeList) UnmarshalText(data []byte) error {
	fields := strings.Split(string(data), ",")
	if len(fields) == 0 {
		return nil
	}
	ttl := make(LocationTimeList, len(fields))
	for i, v := range fields {
		ttl[i] = ParseLocationTime(v)
	}
	*x = ttl
	return nil
}

func (x LocationTimeList) Contains(t LocationTime) bool {
	for _, v := range x {
		if v == t {
			return true
		}
	}
	return false
}

type FunctionType string

const (
	FunctionMS_M         FunctionType = "M-MID_SIDE"
	FunctionMS_S         FunctionType = "S-MID_SIDE"
	FunctionXY_X         FunctionType = "X-X_Y "
	FunctionXY_Y         FunctionType = "Y-X_Y"
	FunctionMix_L        FunctionType = "L-MIX"
	FunctionMix_R        FunctionType = "R-MIX"
	FunctionMix          FunctionType = "MIX"
	FunctionLeft         FunctionType = "LEFT"
	FunctionRight        FunctionType = "RIGHT"
	FunctionDMS_F        FunctionType = "F-DMS"
	FunctionDMS_8        FunctionType = "8-DMS"
	FunctionDMS_R        FunctionType = "R-DMS"
	FunctionLCR_L        FunctionType = "L-LCR"
	FunctionLCR_C        FunctionType = "C-LCR"
	FunctionLCR_R        FunctionType = "R-LCR"
	FunctionLCRS_L       FunctionType = "L-LCRS"
	FunctionLCRS_C       FunctionType = "C-LCRS"
	FunctionLCRS_R       FunctionType = "R-LCRS"
	FunctionLCRS_S       FunctionType = "S-LCRS"
	Function51_L         FunctionType = "L-5.1"
	Function51_C         FunctionType = "C-5.1"
	Function51_R         FunctionType = "R-5.1"
	Function51_Ls        FunctionType = "Ls-5.1"
	Function51_Rs        FunctionType = "Rs-5.1"
	Function51_LFE       FunctionType = "LFE-5.1"
	Function71_L         FunctionType = "L-7.1"
	Function71_Lc        FunctionType = "Lc-7.1"
	Function71_C         FunctionType = "C-7.1"
	Function71_Rc        FunctionType = "Rc-7.1"
	Function71_R         FunctionType = "R-7.1"
	Function71_Ls        FunctionType = "Ls-7.1"
	Function71_Rs        FunctionType = "Rs-7.1"
	Function71_LFE       FunctionType = "LFE-7.1"
	FunctionSurroundL    FunctionType = "L-GENERIC"    // Main Layer Front Left
	FunctionSurroundLc   FunctionType = "Lc-GENERIC"   // Main Layer Front Left Center
	FunctionSurroundC    FunctionType = "C-GENERIC"    // Main Layer Front Center
	FunctionSurroundRc   FunctionType = "Rc-GENERIC"   // Main Layer Front Right Center
	FunctionSurroundR    FunctionType = "R-GENERIC"    // Main Layer Front Right
	FunctionSurroundLs   FunctionType = "Ls-GENERIC"   // Main Layer Rear Left
	FunctionSurroundCs   FunctionType = "Cs-GENERIC"   // Main Layer Rear Center
	FunctionSurroundRs   FunctionType = "Rs-GENERIC"   // Main Layer Rear Right
	FunctionSurroundLFE  FunctionType = "LFE-GENERIC"  // LFE
	FunctionSurroundSl   FunctionType = "Sl-GENERIC"   // Main Layer Side Left
	FunctionSurroundSr   FunctionType = "Sr-GENERIC"   // Main Layer Side Right
	FunctionSurroundLcs  FunctionType = "Lcs-GENERIC"  // Main Layer Rear Left Center
	FunctionSurroundRcs  FunctionType = "Rcs-GENERIC"  // Main Layer Rear Right Center
	FunctionSurroundLFE2 FunctionType = "LFE2-GENERIC" // LFE 2
	FunctionSurroundVoG  FunctionType = "VoG-GENERIC"  // Top Layer Voice of God
	FunctionSurroundTl   FunctionType = "Tl-GENERIC"   // Top Layer Front Left
	FunctionSurroundTc   FunctionType = "Tc-GENERIC"   // Top Layer Front Center
	FunctionSurroundTr   FunctionType = "Tr-GENERIC"   // Top Layer Front Right
	FunctionSurroundTrl  FunctionType = "Trl-GENERIC"  // Top Layer Rear Left
	FunctionSurroundTrc  FunctionType = "Trc-GENERIC"  // Top Layer Rear Center
	FunctionSurroundTrr  FunctionType = "Trr-GENERIC"  // Top Layer Rear Right
	FunctionSurroundTsl  FunctionType = "Tsl-GENERIC"  // Top Layer Side Left
	FunctionSurroundTsr  FunctionType = "Tsr-GENERIC"  // Top Layer Side Right
	FunctionSurroundVoD  FunctionType = "VoD-GENERIC"  // Bottom Layer Voice of Devil
	FunctionSurroundBl   FunctionType = "Bl-GENERIC"   // Bottom Layer Front Left
	FunctionSurroundBc   FunctionType = "Bc-GENERIC"   // Bottom Layer Front Center
	FunctionSurroundBr   FunctionType = "Br-GENERIC"   // Bottom Layer Front Right
	FunctionSurroundBrl  FunctionType = "Brl-GENERIC"  // Bottom Layer Rear Left
	FunctionSurroundBrc  FunctionType = "Brc-GENERIC"  // Bottom Layer Rear Center
	FunctionSurroundBrr  FunctionType = "Brr-GENERIC"  // Bottom Layer Rear Right
	FunctionSurroundBsl  FunctionType = "Bsl-GENERIC"  // Bottom Layer Side Left
	FunctionSurroundBsr  FunctionType = "Bsr-GENERIC"  // Bottom Layer Side Right
	FunctionSoundfield_W FunctionType = "W-SOUNDFIELD"
	FunctionSoundfield_X FunctionType = "X-SOUNDFIELD"
	FunctionSoundfield_Y FunctionType = "Y-SOUNDFIELD"
	FunctionSoundfield_Z FunctionType = "Z-SOUNDFIELD"
	FunctionVideo        FunctionType = "VIDEO"
)

type FunctionSet []FunctionType

var (
	FunctionSetMS            FunctionSet = FunctionSet{FunctionMS_M, FunctionMS_S}
	FunctionSetXY            FunctionSet = FunctionSet{FunctionXY_X, FunctionXY_Y}
	FunctionSetStereoDownmix FunctionSet = FunctionSet{FunctionMix_L, FunctionMix_R}
	FunctionSetMonoDownmix   FunctionSet = FunctionSet{FunctionMix}
	FunctionSetStereo        FunctionSet = FunctionSet{FunctionLeft, FunctionRight}
	FunctionSetDMS           FunctionSet = FunctionSet{FunctionDMS_F, FunctionDMS_8, FunctionDMS_R}
	FunctionSetLCR           FunctionSet = FunctionSet{FunctionLCR_L, FunctionLCR_C, FunctionLCR_R}
	FunctionSetLCRS          FunctionSet = FunctionSet{FunctionLCRS_L, FunctionLCRS_C, FunctionLCRS_R, FunctionLCRS_S}
	FunctionSet51            FunctionSet = FunctionSet{Function51_L, Function51_C, Function51_R, Function51_Ls, Function51_Rs, Function51_LFE}
	FunctionSet71            FunctionSet = FunctionSet{Function71_L, Function71_Lc, Function71_C, Function71_Rc, Function71_R, Function71_Ls, Function71_Rs, Function71_LFE}
	FunctionSetSoundfiled    FunctionSet = FunctionSet{FunctionSoundfield_W, FunctionSoundfield_X, FunctionSoundfield_Y, FunctionSoundfield_Z}
)

func (x FunctionSet) Contains(t FunctionType) bool {
	for _, v := range x {
		if v == t {
			return true
		}
	}
	return false
}

func (x *FunctionSet) Add(t FunctionType) {
	if !x.Contains(t) {
		*x = append(*x, t)
	}
}

func (x *FunctionSet) Del(t FunctionType) {
	idx := -1
	for i, v := range *x {
		if v == t {
			idx = i
		}
	}
	if idx > -1 {
		*x = append((*x)[:idx], (*x)[idx+1:]...)
	}
}

func (x *FunctionSet) Merge(s FunctionSet) {
	for _, v := range s {
		x.Add(v)
	}
}

func (x *FunctionSet) Sub(s FunctionSet) {
	for _, v := range s {
		x.Del(v)
	}
}

type TimecodeFlag string

const (
	TimecodeFlagInvalid TimecodeFlag = ""
	TimecodeFlagDF      TimecodeFlag = "DF"
	TimecodeFlagNDF     TimecodeFlag = "NDF"
)

func ParseTimecodeFlag(s string) TimecodeFlag {
	switch s {
	case "DF":
		return TimecodeFlagDF
	case "NDF":
		return TimecodeFlagNDF
	default:
		return TimecodeFlagInvalid
	}
}

func (x *TimecodeFlag) UnmarshalText(data []byte) error {
	s := string(data)
	if f := ParseTimecodeFlag(s); f == TimecodeFlagInvalid {
		return fmt.Errorf("ixml: invalid timecode flag '%s'", s)
	} else {
		*x = f
	}
	return nil
}

type SyncPointType string

const (
	SyncPointRelative SyncPointType = "RELATIVE"
	SyncPointAbsolute SyncPointType = "ABSOLUTE"
)

type SyncPointFunctionType string

const (
	SyncPointPreRecordSamplecount SyncPointFunctionType = "PRE_RECORD_SAMPLECOUNT"
	SyncPointSlateGeneric         SyncPointFunctionType = "SLATE_GENERIC"
	SyncPointHeadSlate            SyncPointFunctionType = "HEAD_SLATE"
	SyncPointTailSlate            SyncPointFunctionType = "TAIL_SLATE"
	SyncPointMarkerGeneric        SyncPointFunctionType = "MARKER_GENERIC"
	SyncPointMarkerAutoplay       SyncPointFunctionType = "MARKER_AUTOPLAY"
	SyncPointMarkerAutoplayStop   SyncPointFunctionType = "MARKER_AUTOPLAYSTOP"
	SyncPointMarkerAutoplayLoop   SyncPointFunctionType = "MARKER_AUTOPLAYLOOP"
	SyncPointGroupOffset          SyncPointFunctionType = "GROUP_OFFSET"
)
