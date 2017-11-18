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

// Package xmpbj implements the XMP Basic Job Ticket namespace as defined by XMP Specification Part 2.
package xmpbj

import (
	"trimmer.io/go-xmp/xmp"
)

// Part 2: 1.2.5.1 Job
type Job struct {
	ID   string  `xmp:"stJob:id,attr"`
	Name string  `xmp:"stJob:name,attr"`
	Url  xmp.Uri `xmp:"stJob:url,attr"`
}

func (x Job) IsZero() bool {
	return x.ID == "" && x.Name == "" && x.Url == ""
}

func (x Job) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t Job
	return e.EncodeElement(_t(x), node)
}

type JobArray []*Job

func (x JobArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x JobArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *JobArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}
