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

// Package dc implements Dublin Core metadata as defined in XMP Specification Part 1.
package dc

import (
	"trimmer.io/go-xmp/xmp"
)

type DataType string

const (
	DataTypeCollection  DataType = "Collection"           // Collection
	DataTypeDataset     DataType = "Dataset"              // Dataset
	DataTypeEvent       DataType = "Event"                // Event
	DataTypeImage       DataType = "Image"                // Image
	DataTypeInteractive DataType = "Interactive Resource" // Interactive Resource
	DataTypeMoving      DataType = "Moving"               // Moving Image
	DataTypePhysical    DataType = "Physical Object"      // Physical Object
	DataTypeService     DataType = "Service"              // Service
	DataTypeSoftware    DataType = "Software"             // Software
	DataTypeSound       DataType = "Sound"                // Sound
	DataTypeStillImage  DataType = "Still Image"          // Still Image
	DataTypeText        DataType = "Text"                 // Text
)

// Locale
type LocaleArray []string

func (x LocaleArray) Default() string {
	if len(x) == 0 {
		return ""
	}
	return x[0]
}

func (x LocaleArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x LocaleArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *LocaleArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	if len(node.Nodes) == 0 && len(node.Value) > 0 {
		*x = append(*x, node.Value)
		return nil
	}
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}
