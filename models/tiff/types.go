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

package tiff

import (
	"github.com/trimmer-io/go-xmp/xmp"
)

type YCbCrSubSampling xmp.IntList

// [2, 1] = YCbCr4:2:2 [2, 2] = YCbCr4:2:0
func (x YCbCrSubSampling) String() string {
	if len(x) < 2 {
		return "undefined"
	}
	switch {
	case x[0] == 1 && x[1] == 1:
		return "4:4:4"
	case x[0] == 1 && x[1] == 2:
		return "4:4:0"
	case x[0] == 1 && x[1] == 4:
		return "4:4:1"
	case x[0] == 2 && x[1] == 1:
		return "4:2:2"
	case x[0] == 2 && x[1] == 2:
		return "4:2:0"
	case x[0] == 2 && x[1] == 4:
		return "4:2:1"
	case x[0] == 4 && x[1] == 1:
		return "4:1:1"
	case x[0] == 4 && x[1] == 2:
		return "4:1:0"
	default:
		return "undefined"
	}
}

func (x YCbCrSubSampling) IsZero() bool {
	return len(x) < 2 || (x[0] == 0 && x[1] == 0)
}

func (x YCbCrSubSampling) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x YCbCrSubSampling) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *YCbCrSubSampling) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}
