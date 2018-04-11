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

package itunes

import (
	"trimmer.io/go-xmp/xmp"
	// "howett.net/plist"
)

type MovieInfo struct {
	AssetInfo     *AssetInfo  `plist:"asset-info"    xmp:"iTunes:AssetInfo,attr"`
	Studio        string      `plist:"studio"        xmp:"iTunes:Studio,attr"`
	Cast          PersonArray `plist:"cast"          xmp:"iTunes:Cast"`
	Directors     PersonArray `plist:"directors"     xmp:"iTunes:Directors"`
	CoDirectors   PersonArray `plist:"codirectors"   xmp:"iTunes:CoDirectors"`
	Screenwriters PersonArray `plist:"screenwriters" xmp:"iTunes:Screenwriters"`
	Producers     PersonArray `plist:"producers"     xmp:"iTunes:Producers"`
	CopyWarning   string      `plist:"copy-warning"  xmp:"iTunes:CopyWarning"`
}

// unmarshal Apple plist style XML file requires to import external
// dependency howett.net/plist
// func (x *MovieInfo) UnmarshalText(data []byte) error {
// 	_, err := plist.Unmarshal(data, x)
// 	return err
// }

type AssetInfo struct {
	FileSize     int64  `plist:"file-size"     xmp:"iTunes:FileSize,attr"`
	Flavor       string `plist:"flavor"        xmp:"iTunes:Flavor,attr"`
	ScreenFormat string `plist:"screen-format" xmp:"iTunes:ScreenFormat,attr"`
	Soundtrack   string `plist:"soundtrack"    xmp:"iTunes:Soundtrack,attr"`
}

type Person struct {
	ID   string `plist:"adamId" xmp:"iTunes:AdamID,attr"`
	Name string `plist:"name"   xmp:"iTunes:Name,attr"`
}

type PersonArray []Person

func (a PersonArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x PersonArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *PersonArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}
