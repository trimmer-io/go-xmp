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

package xmp

import (
	"encoding/xml"
	"fmt"
	"strings"
)

var (
	nodePool                             = make(chan *Node, 5000)
	npAllocs, npFrees, npHits, npReturns int64
)

type Node struct {
	XMLName xml.Name // node name and namespace
	Attr    AttrList // captures all unbound attributes and XMP qualifiers
	Model   Model    // XmpCore, DublinCore, etc
	Value   string
	Nodes   NodeList // child nodes
}

type Attr struct {
	Name  xml.Name
	Value string
}

type AttrList []Attr

func (x AttrList) XML() []xml.Attr {
	l := make([]xml.Attr, len(x))
	for i, v := range x {
		l[i] = xml.Attr(v)
	}
	return l
}

func (x *AttrList) From(l []xml.Attr) {
	*x = make(AttrList, len(l))
	for i, v := range l {
		(*x)[i] = Attr(v)
	}
}

func NewNode(name xml.Name) *Node {
	var n *Node
	select {
	case n = <-nodePool: // Try to get one from the nodePool
		npHits++
		n.XMLName = name
		n.Attr = nil
		n.Nodes = nil
		n.Model = nil
		n.Value = ""
	default: // All in use, create a new, temporary:
		npAllocs++
		n = &Node{
			XMLName: name,
			Attr:    nil,
			Nodes:   nil,
			Model:   nil,
			Value:   "",
		}
	}

	return n
}

func (n *Node) Close() {
	for _, v := range n.Nodes {
		v.Close()
	}
	n.XMLName = xml.Name{}
	n.Attr = nil
	n.Nodes = nil
	n.Model = nil
	n.Value = ""
	select {
	case nodePool <- n: // try to put back into the nodePool
		npReturns++
	default: // pool is full, will be garbage collected
		npFrees++
	}
}

func copyNode(x *Node) *Node {
	n := NewNode(x.XMLName)
	n.Value = x.Value
	n.Model = x.Model
	n.Attr = make([]Attr, 0, len(x.Attr))
	copy(n.Attr, x.Attr)
	n.Nodes = copyNodes(x.Nodes)
	return n
}

func copyNodes(x NodeList) NodeList {
	l := make(NodeList, 0, len(x))
	for _, v := range x {
		l = append(l, copyNode(v))
	}
	return l
}

type NodeList []*Node

func (x *NodeList) AddNode(n *Node) {
	for i, l := 0, len(*x); i < l; i++ {
		if (*x)[i].XMLName.Local == n.XMLName.Local {
			(*x)[i] = n
			return
		}
	}
	*x = append(*x, n)
}

func (x *NodeList) FindNode(ns *Namespace) *Node {
	prefix := ns.GetName()
	for _, v := range *x {
		if v.Name() == prefix {
			return v
		}
		if v.Model != nil && v.Model.Can(prefix) {
			return v
		}
	}
	return nil
}

// keep list of nodes unique, overwrite contents when names equal
func (n *Node) AddNode(x *Node) {
	if x == n {
		panic(fmt.Errorf("xmp: node loop detected"))
	}
	n.Nodes.AddNode(x)
}

func (n *Node) AddAttr(attr Attr) {
	for i, l := 0, len(n.Attr); i < l; i++ {
		if n.Attr[i].Name.Local == attr.Name.Local {
			n.Attr[i].Value = attr.Value
			return
		}
	}
	n.Attr = append(n.Attr, attr)
}

// keep list of attributes unique, overwrite value when names equal
func (n *Node) AddStringAttr(name, value string) {
	n.AddAttr(Attr{Name: xml.Name{Local: name}, Value: value})
}

func (n *Node) GetAttr(ns, name string) []Attr {
	l := make([]Attr, 0)
	for _, v := range n.Attr {
		if ns != "" && v.Name.Space != ns {
			continue
		}
		if name != "" && dropNsName(v.Name.Local) != name {
			continue
		}
		l = append(l, v)
	}
	return l
}

func (n *Node) Name() string {
	return dropNsName(n.XMLName.Local)
}

func (n *Node) FullName() string {
	if n.XMLName.Space != "" {
		return NsRegistry.Short(n.XMLName.Space, n.XMLName.Local)
	}
	return n.XMLName.Local
}

func (n *Node) Namespace() string {
	ns := n.XMLName.Space
	if ns == "" {
		ns = getPrefix(n.XMLName.Local)
	}
	return ns
}

func (n *Node) IsZero() bool {
	empty := n.Model == nil && n.Value == "" && len(n.Attr) == 0
	if !empty {
		return false
	}
	for _, v := range n.Nodes {
		empty = empty && v.IsZero()
	}
	return empty
}

func (n *Node) translate(d *Decoder) {
	d.translate(&n.XMLName)
	// Note: don't use `for .. range` here because it copies
	// structures, but we intend to alter the node tree
	for i, l := 0, len(n.Attr); i < l; i++ {
		d.translate(&n.Attr[i].Name)
	}
	for i, l := 0, len(n.Nodes); i < l; i++ {
		n.Nodes[i].translate(d)
	}
}

func (n *Node) Namespaces(d *Document) NamespaceList {
	m := make(map[string]bool)

	// keep node namespace
	if name := n.Namespace(); name != "" {
		m[name] = true
	}

	// add model namespaces
	if n.Model != nil {
		for _, v := range n.Model.Namespaces() {
			m[v.GetName()] = true
		}
	}

	// walk attributes
	for _, v := range n.Attr {
		m[getPrefix(v.Name.Local)] = true
	}

	// walk subnodes and capture used namespaces
	var l NamespaceList
	for _, v := range n.Nodes {
		l = append(l, v.Namespaces(d)...)
	}

	// l := make(NamespaceList, 0)
	for name, _ := range m {
		ns := d.findNsByPrefix(name)
		if ns != nil && ns != nsRDF && ns != nsXML {
			l = append(l, ns)
		}
	}

	// keep unique namespaces only
	return l.RemoveDups()
}

func (n *Node) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.XMLName.Local == "" {
		return nil
	}

	start.Name = n.XMLName
	start.Attr = n.Attr.XML()
	if n.Model != nil {
		return e.EncodeElement(struct {
			Data  Model
			Nodes []*Node
		}{
			Data:  n.Model,
			Nodes: n.Nodes,
		}, start)

	} else {
		return e.EncodeElement(struct {
			Data  string `xml:",chardata"`
			Nodes []*Node
		}{
			Data:  n.Value,
			Nodes: n.Nodes,
		}, start)

	}
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var nodes []*Node
	var done bool
	for !done {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.CharData:
			n.Value = strings.TrimSpace(string(t))
		case xml.StartElement:
			x := NewNode(emptyName)
			x.UnmarshalXML(d, t)
			nodes = append(nodes, x)
		case xml.EndElement:
			done = true
		}
	}
	n.XMLName = start.Name
	n.Attr.From(start.Attr)
	n.Nodes = nodes
	return nil
}
