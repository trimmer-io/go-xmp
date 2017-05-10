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

package riff

import (
	"github.com/echa/go-xmp/xmp"
	"strings"
)

type StringArray xmp.StringArray

func (a StringArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x *StringArray) UnmarshalText(data []byte) error {
	*x = strings.Split(string(data), "; ")
	return nil
}

func (x StringArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *StringArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type AltString xmp.AltString

func (x *AltString) UnmarshalText(data []byte) error {
	as := xmp.AltString{}
	for _, v := range strings.Split(string(data), "; ") {
		as.Add("eng", v)
	}
	*x = AltString(as)
	return nil
}

func (x AltString) Typ() xmp.ArrayType {
	return xmp.ArrayTypeAlternative
}

func (x AltString) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, xmp.AltString(x).Typ(), xmp.AltString(x))
}

func (x *AltString) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, xmp.AltString(*x).Typ(), (*xmp.AltString)(x))
}
