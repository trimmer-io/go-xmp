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
	"fmt"
	"reflect"
)

type Document struct {
	// name and version of the toolkit that generated an XMP document instance
	toolkit string

	nodes NodeList

	// flag indicating the document content has changed
	dirty bool

	// rdf:about
	about string

	// local namespace map for tracking registered namespaces
	intNsMap map[string]*Namespace

	// local namespace map for tracking unknown namespaces
	extNsMap map[string]*Namespace
}

// high-level XMP document interface
func NewDocument() *Document {
	d := &Document{
		toolkit:  XMP_TOOLKIT_VERSION,
		nodes:    make(NodeList, 0),
		intNsMap: make(map[string]*Namespace),
		extNsMap: make(map[string]*Namespace),
	}
	return d
}

func (d *Document) SetDirty() {
	d.dirty = true
}

func (d Document) IsDirty() bool {
	return d.dirty
}

func (d *Document) Close() {
	if d == nil {
		return
	}
	for _, v := range d.nodes {
		v.Close()
	}
	d.nodes = nil
}

// cross-model sync, must be explicitly called to merge across models
func (d *Document) SyncModels() error {
	for _, n := range d.nodes {
		if n.Model != nil {
			if err := n.Model.SyncModel(d); err != nil {
				return err
			}
		}
	}
	return nil
}

// implicit sync to align standard properties across models (e.g. xmp <-> photoshop)
// called each time a model is loaded from XMP/XML or XMP/JSON
func (d *Document) syncFromXMP() error {
	for _, n := range d.nodes {
		if n.Model != nil {
			if err := n.Model.SyncFromXMP(d); err != nil {
				return err
			}
		}
	}
	return nil
}

// called each time a model is stored as XMP/XML or XMP/JSON
func (d *Document) syncToXMP() error {
	if !d.dirty {
		return nil
	}
	for _, n := range d.nodes {
		if n.Model != nil {
			if err := n.Model.SyncToXMP(d); err != nil {
				return err
			}
		}
	}
	d.dirty = false
	return nil
}

func (d *Document) Namespaces() NamespaceList {
	var l NamespaceList
	for _, v := range d.intNsMap {
		l = append(l, v)
	}
	for _, v := range d.extNsMap {
		l = append(l, v)
	}
	return l
}

func (d *Document) Nodes() NodeList {
	return d.nodes
}

func (d *Document) FindNode(ns *Namespace) *Node {
	return d.nodes.FindNode(ns)
}

func (d *Document) FindModel(ns *Namespace) Model {
	prefix := ns.GetName()
	for _, v := range d.nodes {
		if v.Model != nil && v.Model.Can(prefix) {
			return v.Model
		}
	}
	return nil
}

func (d Document) FindNs(name, uri string) *Namespace {
	var ns *Namespace
	if uri != "" {
		ns = d.findNsByURI(uri)
	}
	if ns == nil {
		ns = d.findNsByPrefix(getPrefix(name))
	}
	return ns
}

func (d *Document) MakeModel(ns *Namespace) (Model, error) {
	m := d.FindModel(ns)
	if m == nil {
		if m = ns.NewModel(); m == nil {
			return nil, fmt.Errorf("xmp: cannot create '%s' model", ns.GetName())
		}
		n := NewNode(ns.XMLName(""))
		n.Model = m
		d.nodes = append(d.nodes, n)
		for _, v := range m.Namespaces() {
			d.intNsMap[v.GetURI()] = v
		}
		d.SetDirty()
	}
	return m, nil
}

func (d *Document) AddModel(v Model) (*Node, error) {
	if v == nil {
		return nil, fmt.Errorf("xmp: AddModel called with nil model")
	}
	ns := v.Namespaces()
	if len(ns) == 0 {
		return nil, fmt.Errorf("xmp: model '%s' must implement at least one namespace", reflect.TypeOf(v).String())
	}
	for _, v := range ns {
		d.intNsMap[v.GetURI()] = v
	}
	n := d.FindNode(ns[0])
	if n == nil {
		n = NewNode(ns[0].XMLName(""))
		d.nodes = append(d.nodes, n)
	}
	n.Model = v
	d.SetDirty()
	return n, nil
}

func (d *Document) RemoveNamespace(ns *Namespace) bool {
	if ns == nil {
		return false
	}
	name := ns.GetName()
	removed := false
	for i, v := range d.nodes {
		if v.Name() == name {
			d.nodes = append(d.nodes[:i], d.nodes[i+1:]...)
			removed = true
			delete(d.intNsMap, ns.GetURI())
			delete(d.extNsMap, ns.GetURI())
			v.Close()
			d.SetDirty()
			break
		}
	}
	return removed
}

func (d *Document) RemoveNamespaceByName(n string) bool {
	ns := d.findNsByPrefix(n)
	return d.RemoveNamespace(ns)
}

// will delete all document nodes matching the list
func (d *Document) RemoveNamespaces(remove NamespaceList) bool {
	removed := false
	for _, ns := range remove {
		r := d.RemoveNamespace(ns)
		removed = removed || r
	}
	return removed
}

// will keep document nodes matching the list and delete all others.
// returns true if any namespace was removed.
func (d *Document) FilterNamespaces(keep NamespaceList) bool {
	removed := false
	allNs := make(NamespaceList, 0)
	// stage one: collect all top-level namespaces
	for _, v := range d.nodes {
		ns := d.findNsByPrefix(v.Namespace())
		if ns != nil {
			allNs = append(allNs, ns)
		}
	}
	// stage two: remove
	for _, v := range allNs {
		if !keep.Contains(v) {
			r := d.RemoveNamespace(v)
			removed = removed || r
		}
	}
	return removed
}

type DocumentArray []*Document

func (x DocumentArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func (x DocumentArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	return MarshalArray(e, node, x.Typ(), x)
}

func (x *DocumentArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

func (d Document) findNsByURI(uri string) *Namespace {
	if v, ok := d.intNsMap[uri]; ok {
		return v
	}
	if v, ok := d.extNsMap[uri]; ok {
		return v
	}
	return nil
}

func (d Document) findNsByPrefix(pre string) *Namespace {
	for _, v := range d.intNsMap {
		if v.GetName() == pre {
			return v
		}
	}
	for _, v := range d.extNsMap {
		if v.GetName() == pre {
			return v
		}
	}
	if ns, err := NsRegistry.GetNamespace(pre); err == nil {
		return ns
	}
	return nil
}
