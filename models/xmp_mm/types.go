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
	"trimmer.io/go-xmp/models/xmp_dm"
	"trimmer.io/go-xmp/xmp"
)

// Part 1: 8.2.2.9 ResourceRef + Part 2: 1.2.4.1 ResourceRef
type ResourceRef struct {
	AlternatePaths     xmp.UriArray         `xmp:"stRef:alternatePaths"`
	OriginalDocumentID xmp.GUID             `xmp:"stRef:originalDocumentID,attr"`
	DocumentID         xmp.GUID             `xmp:"stRef:documentID,attr"`
	FilePath           xmp.Uri              `xmp:"stRef:filePath,attr"`
	FromPart           *xmpdm.Part          `xmp:"stRef:fromPart,attr"`
	InstanceID         xmp.GUID             `xmp:"stRef:instanceID,attr"`
	LastModifyDate     xmp.Date             `xmp:"stRef:lastModifyDate,attr"`
	Manager            xmp.AgentName        `xmp:"stRef:manager,attr"`
	ManagerVariant     string               `xmp:"stRef:managerVariant,attr"`
	ManageTo           xmp.Uri              `xmp:"stRef:manageTo,attr"`
	ManageUI           xmp.Uri              `xmp:"stRef:manageUI,attr"`
	MaskMarkers        xmpdm.MaskType       `xmp:"stRef:maskMarkers,attr"`
	PartMapping        string               `xmp:"stRef:partMapping,attr"`
	RenditionClass     xmpdm.RenditionClass `xmp:"stRef:renditionClass,attr"`
	RenditionParams    string               `xmp:"stRef:renditionParams,attr"`
	ToPart             *xmpdm.Part          `xmp:"stRef:toPart,attr"`
	VersionID          string               `xmp:"stRef:versionID,attr"`
}

func (x ResourceRef) IsZero() bool {
	return len(x.AlternatePaths) == 0 &&
		x.OriginalDocumentID.IsZero() &&
		x.DocumentID.IsZero() &&
		x.FilePath == "" &&
		(x.FromPart == nil || x.FromPart.IsZero()) &&
		x.InstanceID.IsZero() &&
		x.LastModifyDate.IsZero() &&
		x.Manager.IsZero() &&
		x.ManagerVariant == "" &&
		x.ManageTo == "" &&
		x.MaskMarkers == "" &&
		x.PartMapping == "" &&
		x.RenditionClass == "" &&
		x.RenditionParams == "" &&
		(x.ToPart == nil || x.ToPart.IsZero()) &&
		x.VersionID == ""
}

func (x ResourceRef) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	if x.IsZero() {
		return nil
	}
	type _t ResourceRef
	return e.EncodeElement(_t(x), node)
}

type ResourceRefArray []*ResourceRef

func (x ResourceRefArray) Index(did xmp.GUID) int {
	for i, v := range x {
		if v.DocumentID == did {
			return i
		}
	}
	return -1
}

func (x *ResourceRefArray) DeleteIndex(idx int) {
	if idx < 0 || idx > len(*x) {
		return
	}
	*x = append((*x)[:idx], ((*x)[idx+1:])...)
}

func (x ResourceRefArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x ResourceRefArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *ResourceRefArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// Part 2: 1.2.4.2 Version
type StVersion struct {
	Comments   string        `xmp:"stVer:comments,attr"`
	Event      ResourceEvent `xmp:"stVer:event"`
	Modifier   string        `xmp:"stVer:modifier,attr"`
	ModifyDate xmp.Date      `xmp:"stVer:modifyDate,attr"`
	Version    string        `xmp:"stVer:version,attr"`
}

type StVersionArray []*StVersion

func (x StVersionArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeUnordered
}

func (x StVersionArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *StVersionArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

// 1.2.4 ResourceEvent
type ResourceEvent struct {
	Action        ActionType     `xmp:"stEvt:action,attr"`
	Changed       xmpdm.PartList `xmp:"stEvt:changed,attr"`
	InstanceID    xmp.GUID       `xmp:"stEvt:instanceID,attr"`
	Parameters    string         `xmp:"stEvt:parameters,attr"`
	SoftwareAgent xmp.AgentName  `xmp:"stEvt:softwareAgent,attr"`
	When          xmp.Date       `xmp:"stEvt:when,attr"`
}

type ResourceEventArray []*ResourceEvent

func (x ResourceEventArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x ResourceEventArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *ResourceEventArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type ActionType string

const (
	ActionConverted      ActionType = "converted"
	ActionCopied         ActionType = "copied"
	ActionCreated        ActionType = "created"
	ActionCropped        ActionType = "cropped"
	ActionEdited         ActionType = "edited"
	ActionFiltered       ActionType = "filtered"
	ActionFormatted      ActionType = "formatted"
	ActionVersionUpdated ActionType = "version_updated"
	ActionPrinted        ActionType = "printed"
	ActionPublished      ActionType = "published"
	ActionManaged        ActionType = "managed"
	ActionProduced       ActionType = "produced"
	ActionResized        ActionType = "resized"
	ActionSaved          ActionType = "saved"
	// more types not in XMP standard
	ActionAdded   ActionType = "media_added"
	ActionRemoved ActionType = "media_removed"
	ActionForked  ActionType = "forked"
	ActionMerged  ActionType = "merged"
)
