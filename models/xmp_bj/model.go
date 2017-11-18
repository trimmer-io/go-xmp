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
	"fmt"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsXmpBJ = xmp.NewNamespace("xmpBJ", "http://ns.adobe.com/xap/1.0/bj/", NewModel)
	nsStJob = xmp.NewNamespace("stJob", "http://ns.adobe.com/xap/1.0/sType/Job#", nil)
)

func init() {
	xmp.Register(NsXmpBJ, xmp.XmpMetadata)
	xmp.Register(nsStJob)
}

func NewModel(name string) xmp.Model {
	return &JobTicket{}
}

func MakeModel(d *xmp.Document) (*JobTicket, error) {
	m, err := d.MakeModel(NsXmpBJ)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*JobTicket)
	return x, nil
}

func FindModel(d *xmp.Document) *JobTicket {
	if m := d.FindModel(NsXmpBJ); m != nil {
		return m.(*JobTicket)
	}
	return nil
}

type JobTicket struct {
	JobRef JobArray `xmp:"xmpBJ:JobRef"`
}

func (x JobTicket) Can(nsName string) bool {
	return NsXmpBJ.GetName() == nsName
}

func (x JobTicket) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsXmpBJ}
}

func (x *JobTicket) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x JobTicket) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *JobTicket) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *JobTicket) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsXmpBJ.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *JobTicket) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsXmpBJ.GetName(), err)
	}
	return nil
}
