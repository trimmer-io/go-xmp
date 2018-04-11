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

package xmp

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type UnknownNamespaceError struct {
	Name xml.Name
}

func (e *UnknownNamespaceError) Error() string {
	if e.Name.Space == "" {
		return fmt.Sprintf("xmp: undefined XML namespace for node %s", e.Name.Local)
	} else {
		return fmt.Sprintf("xmp: undefined XML namespace for node %s:%s", e.Name.Space, e.Name.Local)
	}
}

type Namespace struct {
	Name    string
	URI     string
	Factory ModelFactory
}

// Namespace Groups
type NamespaceGroup string
type NamespaceGroupList []NamespaceGroup

const (
	NoMetadata     NamespaceGroup = ""
	XmpMetadata    NamespaceGroup = "xmp"
	ImageMetadata  NamespaceGroup = "image"
	MusicMetadata  NamespaceGroup = "music"
	MovieMetadata  NamespaceGroup = "movie"
	SoundMetadata  NamespaceGroup = "sound"
	CameraMetadata NamespaceGroup = "camera"
	VfxMetadata    NamespaceGroup = "vfx"
	RightsMetadata NamespaceGroup = "rights"
)

func ParseNamespaceGroup(s string) NamespaceGroup {
	switch s {
	case "xmp":
		return XmpMetadata
	case "image":
		return ImageMetadata
	case "music":
		return MusicMetadata
	case "movie":
		return MovieMetadata
	case "sound":
		return SoundMetadata
	case "camera":
		return CameraMetadata
	case "vfx":
		return VfxMetadata
	case "rights":
		return RightsMetadata
	default:
		return NoMetadata
	}
}

func (g NamespaceGroup) Namespaces() NamespaceList {
	l, _ := GetGroupNamespaces(g)
	return l
}

func (g NamespaceGroup) Contains(ns *Namespace) bool {
	for _, v := range g.Namespaces() {
		if v.GetName() == ns.GetName() {
			return true
		}
	}
	return false
}

func (l NamespaceGroupList) Contains(ns *Namespace) bool {
	for _, v := range l {
		if v.Contains(ns) {
			return true
		}
	}
	return false
}

func NewNamespace(name, uri string, factory ModelFactory) *Namespace {
	return &Namespace{
		Name:    name,
		URI:     uri,
		Factory: factory,
	}
}

func (n Namespace) NewModel() Model {
	if n.Factory != nil {
		return n.Factory(n.GetName())
	}
	return nil
}

func (n Namespace) MatchAttrName(name string) bool {
	return getPrefix(name) == n.Name
}

func (n Namespace) MatchXMLName(name xml.Name) bool {
	if name.Space != "" {
		return name.Space == n.URI
	}
	return getPrefix(name.Local) == n.Name
}

func (n Namespace) Expand(local string) string {
	if local == "" {
		return n.GetName()
	}
	return strings.Join([]string{n.GetName(), local}, ":")
}

func (n Namespace) Split(name string) string {
	return stripPrefix(name)
}

func (n Namespace) GetName() string {
	return n.Name
}

func (n Namespace) GetURI() string {
	return n.URI
}

func (n Namespace) GetAttr() Attr {
	return Attr{Name: xml.Name{Local: "xmlns:" + n.Name}, Value: n.URI}
}

func (n Namespace) XMLName(local string) xml.Name {
	return xml.Name{
		Space: "",
		Local: n.Expand(local),
	}
}

func (n Namespace) RootName() xml.Name {
	return xml.Name{
		Space: "",
		Local: n.Name,
	}
}

type NamespaceList []*Namespace

func (l NamespaceList) Contains(ns *Namespace) bool {
	if ns == nil {
		return false
	}
	for _, v := range l {
		if v.GetName() == ns.GetName() {
			return true
		}
	}
	return false
}

func (l NamespaceList) ContainsName(n string) bool {
	if n == "" {
		return false
	}
	for _, v := range l {
		if v.GetName() == n {
			return true
		}
	}
	return false
}

func (l *NamespaceList) RemoveDups() NamespaceList {
	found := make(map[string]bool)
	j := 0
	for i, v := range *l {
		u := v.GetURI()
		if !found[u] {
			found[u] = true
			(*l)[j] = (*l)[i]
			j++
		}
	}
	*l = (*l)[:j]
	return *l
}

func stripPrefix(n string) string {
	if i := strings.Index(n, ":"); i > -1 {
		return n[i+1:]
	}
	return n
}

func getPrefix(n string) string {
	if i := strings.Index(n, ":"); i > -1 {
		return n[:i]
	}
	return n
}

// true when n has a non-empty namespace prefix
func hasPrefix(n string) bool {
	return strings.Index(n, ":") > 0
}
