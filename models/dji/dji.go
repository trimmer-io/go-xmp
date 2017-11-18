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

// Package dji implements metadata models written by video cameras on DJI drones.
package dji

import (
	"fmt"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsDJI = xmp.NewNamespace("drone-dji", "http://www.dji.com/drone-dji/1.0/", NewModel)
)

func init() {
	xmp.Register(NsDJI, xmp.CameraMetadata)
}

func NewModel(name string) xmp.Model {
	return &DJI{}
}

func MakeModel(d *xmp.Document) (*DJI, error) {
	m, err := d.MakeModel(NsDJI)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*DJI)
	return x, nil
}

func FindModel(d *xmp.Document) *DJI {
	if m := d.FindModel(NsDJI); m != nil {
		return m.(*DJI)
	}
	return nil
}

func (m *DJI) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsDJI}
}

func (m *DJI) Can(nsName string) bool {
	return nsName == NsDJI.GetName()
}

func (x *DJI) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x DJI) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *DJI) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *DJI) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsDJI.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *DJI) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsDJI.GetName(), err)
	}
	return nil
}

type DJI struct {
	AbsoluteAltitude  float32 `drone-dji:"AbsoluteAltitude"  qt:"-"     xmp:"drone-dji:AbsoluteAltitude"`  //  "+543.44",
	FlightPitchDegree float32 `drone-dji:"FlightPitchDegree" qt:"©fpt"  xmp:"drone-dji:FlightPitchDegree"` //  "+4.80",
	FlightRollDegree  float32 `drone-dji:"FlightRollDegree"  qt:"©frl"  xmp:"drone-dji:FlightRollDegree"`  //  "+1.30",
	FlightYawDegree   float32 `drone-dji:"FlightYawDegree"   qt:"©fyw"  xmp:"drone-dji:FlightYawDegree"`   //  "-1.90",
	GimbalPitchDegree float32 `drone-dji:"GimbalPitchDegree" qt:"©gpt"  xmp:"drone-dji:GimbalPitchDegree"` //  "-90.00",
	GimbalRollDegree  float32 `drone-dji:"GimbalRollDegree"  qt:"©grl"  xmp:"drone-dji:GimbalRollDegree"`  //  "+0.00",
	GimbalYawDegree   float32 `drone-dji:"GimbalYawDegree"   qt:"©gyw"  xmp:"drone-dji:GimbalYawDegree"`   //  "-2.00",
	RelativeAltitude  float32 `drone-dji:"RelativeAltitude"  qt:"-"     xmp:"drone-dji:RelativeAltitude"`  //  "+46.60"
	SpeedX            float32 `drone-dji:"SpeedX"            qt:"©xsp"  xmp:"drone-dji:SpeedX"`
	SpeedY            float32 `drone-dji:"SpeedY"            qt:"©ysp"  xmp:"drone-dji:SpeedY"` // - +0.00
	SpeedZ            float32 `drone-dji:"SpeedZ"            qt:"©zsp"  xmp:"drone-dji:SpeedZ"` //  - +0.00,+0.40
	Model             string  `drone-dji:"Model"             qt:"©mdl"  xmp:"tiff:Model"`
}

// unknown
// "©dji"  UserData_dji
// "©res"  UserData_res
// "©uid"  UserData_uid
