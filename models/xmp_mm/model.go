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

// Package xmpmm implements the XMP Media Management namespace as defined by XMP Specification Part 2.
package xmpmm

import (
	"fmt"
	"github.com/trimmer-io/go-xmp/models/xmp_dm"
	"github.com/trimmer-io/go-xmp/xmp"
)

var (
	NsXmpMM = xmp.NewNamespace("xmpMM", "http://ns.adobe.com/xap/1.0/mm/", NewModel)
	nsStEvt = xmp.NewNamespace("stEvt", "http://ns.adobe.com/xap/1.0/sType/ResourceEvent#", nil)
	nsStRef = xmp.NewNamespace("stRef", "http://ns.adobe.com/xap/1.0/sType/ResourceRef#", nil)
	nsStVer = xmp.NewNamespace("stVer", "http://ns.adobe.com/xap/1.0/sType/Version#", nil)
)

func init() {
	xmp.Register(NsXmpMM, xmp.XmpMetadata)
	xmp.Register(nsStEvt)
	xmp.Register(nsStRef)
	xmp.Register(nsStVer)
}

func NewModel(name string) xmp.Model {
	return &XmpMM{}
}

func MakeModel(d *xmp.Document) (*XmpMM, error) {
	m, err := d.MakeModel(NsXmpMM)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*XmpMM)
	return x, nil
}

func FindModel(d *xmp.Document) *XmpMM {
	if m := d.FindModel(NsXmpMM); m != nil {
		return m.(*XmpMM)
	}
	return nil
}

type XmpMM struct {
	DerivedFrom        *ResourceRef         `xmp:"xmpMM:DerivedFrom"`
	DocumentID         xmp.GUID             `xmp:"xmpMM:DocumentID"`
	History            ResourceEventArray   `xmp:"xmpMM:History"`
	Ingredients        ResourceRefArray     `xmp:"xmpMM:Ingredients"`
	InstanceID         xmp.GUID             `xmp:"xmpMM:InstanceID"`
	ManagedFrom        *ResourceRef         `xmp:"xmpMM:ManagedFrom"`
	Manager            xmp.AgentName        `xmp:"xmpMM:Manager"`
	ManageTo           xmp.Uri              `xmp:"xmpMM:ManageTo"`
	ManageUI           xmp.Uri              `xmp:"xmpMM:ManageUI"`
	ManagerVariant     string               `xmp:"xmpMM:ManagerVariant"`
	OriginalDocumentID xmp.GUID             `xmp:"xmpMM:OriginalDocumentID"`
	Pantry             xmp.ExtensionArray   `xmp:"xmpMM:Pantry,omitempty"`
	RenditionClass     xmpdm.RenditionClass `xmp:"xmpMM:RenditionClass"`
	RenditionParams    string               `xmp:"xmpMM:RenditionParams"`
	VersionID          string               `xmp:"xmpMM:VersionID"`
	Versions           StVersionArray       `xmp:"xmpMM:Versions"`
}

func (x XmpMM) Can(nsName string) bool {
	return NsXmpMM.GetName() == nsName
}

func (x XmpMM) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsXmpMM}
}

func (x *XmpMM) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *XmpMM) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x XmpMM) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *XmpMM) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *XmpMM) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsXmpMM.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *XmpMM) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsXmpMM.GetName(), err)
	}
	return nil
}

func (x *XmpMM) AddPantry(d *xmp.Document) {
	node := xmp.NewNode(xmp.NewName("rdf:Description"))
	cpy := xmp.NewDocument()
	cpy.Merge(d, xmp.MERGE)
	node.Nodes = cpy.Nodes()
	x.Pantry = append(x.Pantry, (*xmp.Extension)(node))
}

func (x *XmpMM) AddHistory(e *ResourceEvent) {
	x.History = append(x.History, e)
}

func (x *XmpMM) AddVersion(v *StVersion) {
	x.Versions = append(x.Versions, v)
}

func (x *XmpMM) GetLastEvent() *ResourceEvent {
	if l := len(x.History); l > 0 {
		return x.History[l-1]
	}
	return nil
}

func (x *XmpMM) GetPreviousVersion() *StVersion {
	if l := len(x.Versions); l > 1 {
		return x.Versions[l-2]
	}
	return nil
}

func (x *XmpMM) GetLastVersion() *StVersion {
	if l := len(x.Versions); l > 0 {
		return x.Versions[l-1]
	}
	return nil
}

func (x *XmpMM) GetPreviousVersionId() string {
	if v := x.GetPreviousVersion(); v != nil {
		return v.Version
	}
	return ""
}

func (x *XmpMM) GetLastVersionId() string {
	if v := x.GetLastVersion(); v != nil {
		return v.Version
	}
	return ""
}

func (x *XmpMM) SetPreviousVersionId(version string) {
	if l := len(x.Versions); l > 0 {
		x.Versions[l-1].Version = version
	}
}

func (x *XmpMM) SelfResourceRef(version string) *ResourceRef {
	ref := &ResourceRef{
		OriginalDocumentID: x.OriginalDocumentID,
		DocumentID:         x.DocumentID,
		InstanceID:         x.InstanceID,
		Manager:            x.Manager,
		ManagerVariant:     x.ManagerVariant,
		ManageTo:           x.ManageTo,
		ManageUI:           x.ManageUI,
		RenditionClass:     x.RenditionClass,
		RenditionParams:    x.RenditionParams,
		MaskMarkers:        xmpdm.MaskNone,
		VersionID:          version,
	}
	if e := x.GetLastEvent(); e != nil {
		ref.LastModifyDate = e.When
	}
	return ref
}

// assumes InstanceID and VersionID are changed outside
func (x *XmpMM) AppendVersionHistory(action ActionType, modifier, changed string, date xmp.Date) {
	// append change to last version
	v := x.GetLastVersion()

	// make new version if none exists or if list does not contain
	// entry for the current version
	if v == nil || v.Version != x.VersionID {
		v = &StVersion{
			Event: ResourceEvent{
				Action:        action,
				Changed:       xmpdm.NewPartList(changed),
				InstanceID:    x.InstanceID,
				SoftwareAgent: xmp.Agent,
				When:          date,
			},
			Modifier:   modifier,
			ModifyDate: date,
			Version:    x.VersionID,
		}
		x.AddVersion(v)
		return
	}
	v.Event.Changed.Add(changed)
}
