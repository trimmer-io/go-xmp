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

// Creative Commons Metadata
//
// - URL: https://wiki.creativecommons.org/wiki/XMP
// - Namespace: http://creativecommons.org/ns#
// - prefix: cc

// Package cc implements Creative Commons licensing metadata.
package cc

import (
	"fmt"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsCC = xmp.NewNamespace("cc", "http://creativecommons.org/ns#", NewModel)
)

func init() {
	xmp.Register(NsCC, xmp.RightsMetadata)
}

func NewModel(name string) xmp.Model {
	return &CC{}
}

func MakeModel(d *xmp.Document) (*CC, error) {
	m, err := d.MakeModel(NsCC)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*CC)
	return x, nil
}

func FindModel(d *xmp.Document) *CC {
	if m := d.FindModel(NsCC); m != nil {
		return m.(*CC)
	}
	return nil
}

type CC struct {
	License         xmp.Uri `xmp:"cc:license"`
	MorePermissions xmp.Uri `xmp:"cc:morePermissions"`
	AttributionURL  xmp.Uri `xmp:"cc:attributionURL"`
	AttributionName string  `xmp:"cc:attributionName"`
}

func (x CC) Can(nsName string) bool {
	return NsCC.GetName() == nsName
}

func (x CC) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsCC}
}

func (x *CC) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x CC) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *CC) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *CC) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsCC.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *CC) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsCC.GetName(), err)
	}
	return nil
}
