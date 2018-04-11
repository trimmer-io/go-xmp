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
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	nodePool                             = make(chan *Node, 5000)
	npAllocs, npFrees, npHits, npReturns int64
	errNotFound                          = errors.New("not found")
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

func (x AttrList) IsZero() bool {
	if len(x) == 0 {
		return true
	}
	for _, v := range x {
		if v.Value != "" {
			return false
		}
	}
	return true
}

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

var EmptyName = xml.Name{}

func NewName(local string) xml.Name {
	return xml.Name{Local: local}
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
	n.Attr = make([]Attr, len(x.Attr))
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

func (x *NodeList) AddNode(n *Node) *Node {
	for i, l := 0, len(*x); i < l; i++ {
		if (*x)[i].XMLName.Local == n.XMLName.Local {
			(*x)[i] = n
			return n
		}
	}
	return x.AppendNode(n)
}

func (x *NodeList) AppendNode(n *Node) *Node {
	*x = append(*x, n)
	return n
}

func (n *NodeList) FindNode(ns *Namespace) *Node {
	return n.FindNodeByName(ns.GetName())
}

func (n *NodeList) FindNodeByName(prefix string) *Node {
	for _, v := range *n {
		if v.Name() == prefix {
			return v
		}
		if v.Model != nil && v.Model.Can(prefix) {
			return v
		}
	}
	return nil
}

func (x *NodeList) Index(n *Node) int {
	for i, v := range *x {
		if v == n {
			return i
		}
	}
	return -1
}

func (x *NodeList) RemoveNode(n *Node) *Node {
	if idx := x.Index(n); idx > -1 {
		*x = append((*x)[:idx], (*x)[idx+1:]...)
	}
	return n
}

func (n *Node) IsZero() bool {
	empty := n.Model == nil && n.Value == "" && (len(n.Attr) == 0 || n.Attr.IsZero())
	if !empty {
		return false
	}
	for _, v := range n.Nodes {
		empty = empty && v.IsZero()
	}
	return empty
}

func (n *Node) Name() string {
	return stripPrefix(n.XMLName.Local)
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

	for name, _ := range m {
		ns := d.findNsByPrefix(name)
		if ns != nil && ns != nsRDF && ns != nsXML {
			l = append(l, ns)
		}
	}

	// keep unique namespaces only
	return l.RemoveDups()
}

// keep list of nodes unique, overwrite contents when names equal
func (n *Node) AddNode(x *Node) *Node {
	if x == n {
		panic(fmt.Errorf("xmp: node loop detected"))
	}
	return n.Nodes.AddNode(x)
}

// append in any case
func (n *Node) AppendNode(x *Node) *Node {
	if x == n {
		panic(fmt.Errorf("xmp: node loop detected"))
	}
	return n.Nodes.AppendNode(x)
}

func (n *Node) Clear() {
	for _, v := range n.Nodes {
		v.Close()
	}
	n.Nodes = nil
}

func (n *Node) RemoveNode(x *Node) *Node {
	if x == n {
		panic(fmt.Errorf("xmp: node loop detected"))
	}
	return n.Nodes.RemoveNode(x)
}

func (n Node) IsArray() bool {
	if len(n.Nodes) != 1 {
		return false
	}
	switch n.Nodes[0].FullName() {
	case "rdf:Seq", "rdf:Bag", "rdf:Alt":
		return true
	default:
		return false
	}
}

func (n Node) ArrayType() ArrayType {
	if len(n.Nodes) == 1 {
		switch n.Nodes[0].FullName() {
		case "rdf:Seq":
			return ArrayTypeOrdered
		case "rdf:Bag":
			return ArrayTypeUnordered
		case "rdf:Alt":
			return ArrayTypeAlternative
		}
	}
	return ArrayType("")
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
		if name != "" && stripPrefix(v.Name.Local) != name {
			continue
		}
		l = append(l, v)
	}
	return l
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

func (n *Node) GetPath(path Path) (string, error) {
	name, path := path.PopFront()
	name, idx, lang := parsePathSegment(name)
	if idx < -1 {
		return "", fmt.Errorf("path field %s: invalid index", name)
	}
	// fmt.Printf("Get Node path ns=%s name=%s len=%d rest=%v idx=%d lang=%s\n", path.NamespacePrefix(), name, path.Len(), path, idx, lang)
	if name == "" && idx == -1 && lang == "" {
		if path.Len() == 0 {
			return n.Value, nil
		}
		// ignore empty path segments and recurse
		return n.GetPath(path)
	}

	// lookup name in node list or attributes
	node := n.Nodes.FindNodeByName(stripPrefix(name))
	if node != nil {
		switch {
		case idx > -1:
			// drill two levels deep into array nodes bag/seq+li
			if len(node.Nodes) == 0 || len(node.Nodes[0].Nodes) <= idx {
				return "", nil
			}
			return node.Nodes[0].Nodes[idx].GetPath(path)
		case lang != "":
			// drill two levels deep into alt-array nodes alt+li
			if len(node.Nodes) == 0 || len(node.Nodes[0].Nodes) == 0 {
				return "", nil
			}
			for _, v := range node.Nodes[0].Nodes {
				attr := v.GetAttr("", "lang")
				if len(attr) == 0 {
					continue
				}
				if attr[0].Value == string(lang) {
					return v.Value, nil
				}
			}
			return "", nil
		default:
			return node.GetPath(path)
		}
	}

	if attr := n.GetAttr("", stripPrefix(name)); len(attr) > 0 {
		return attr[0].Value, nil
	}
	return "", nil
}

func (n *Node) SetPath(path Path, value string, flags SyncFlags) error {
	name, path := path.PopFront()
	name, idx, lang := parsePathSegment(name)
	if idx < -1 {
		return fmt.Errorf("path field %s: invalid index", name)
	}

	// fmt.Printf("Set Node path ns=%s len=%d, path=%s, name=%s rest=%v idx=%d lang=%s\n", path.NamespacePrefix(), path.Len(), path.String(), name, path, idx, lang)
	if name == "" && idx == -1 && lang == "" {
		if path.Len() == 0 {
			n.Value = value
		}
		return nil
	}

	// handle attribute
	if attr := n.GetAttr("", stripPrefix(name)); len(attr) > 0 {
		switch {
		case flags&REPLACE > 0 && value != "":
			attr[0].Value = value
			return nil
		case flags&DELETE > 0 && value == "":
			attr[0].Value = value
			// will be ignored on next marshal
			return nil
		}
	}

	// handle nodes
	node := n.Nodes.FindNodeByName(stripPrefix(name))
	if node == nil {
		if flags&CREATE > 0 && value != "" {
			if name != "" && !hasPrefix(name) {
				name = path.NamespacePrefix() + ":" + name
			}
			node = n.AddNode(NewNode(NewName(name)))
		} else {
			return fmt.Errorf("CREATE flag required to make node '%s'", name)
		}
	}
	switch {
	case idx > -1 || (node.IsArray() && lang == ""):
		arr := node.Nodes.FindNodeByName("Seq")
		if arr == nil {
			arr = node.Nodes.FindNodeByName("Bag")
		}
		if arr == nil && flags&CREATE == 0 {
			return fmt.Errorf("CREATE flag required to make array '%s'", name)
		}
		if arr == nil {
			arr = node.AddNode(NewNode(NewName("rdf:Seq")))
		}

		// recurse when we're not at the end of the path
		if path.Len() > 0 {
			if idx < 0 {
				idx = 0
			}
			if l := len(arr.Nodes); l <= idx {
				if flags&(CREATE|APPEND) == 0 && value != "" {
					return fmt.Errorf("CREATE flag required to extend array '%s'", name)
				}
				for ; l <= idx; l++ {
					arr.AppendNode(NewNode(NewName("rdf:li")))
				}
			}
			return arr.Nodes[idx].SetPath(path, value, flags)
		}

		// when at end of path flags tell what to do
		switch {
		case flags&UNIQUE > 0 && value != "" && idx == -1:
			// append if not exists
			for _, v := range arr.Nodes {
				if v.Value == value {
					return nil
				}
			}
			li := arr.AppendNode(NewNode(NewName("rdf:li")))
			li.Value = value

		case flags&APPEND > 0 && value != "" && idx == -1:
			// always append
			li := arr.AppendNode(NewNode(NewName("rdf:li")))
			li.Value = value

		case flags&(REPLACE|CREATE) > 0 && value != "" && idx == -1:
			// replace the entire xmp array
			arr.Clear()
			li := arr.AppendNode(NewNode(NewName("rdf:li")))
			li.Value = value

		case flags&(REPLACE|CREATE) > 0 && value != "" && idx > -1:
			// replace a single item, add intermediate index positions
			if l := len(arr.Nodes); l <= idx {
				for ; l <= idx; l++ {
					arr.AppendNode(NewNode(NewName("rdf:li")))
				}
			}
			arr.Nodes[idx].Value = value

		case flags&DELETE > 0 && value == "" && idx == -1:
			// delete the entire array
			node.RemoveNode(arr).Close()

		case flags&DELETE > 0 && value == "" && idx > -1:
			// delete a single item
			if idx < len(arr.Nodes) {
				arr.Nodes = append(arr.Nodes[:idx], arr.Nodes[idx+1:]...)
			}
		default:
			return fmt.Errorf("unsupported flag combination %v for %s", flags, name)
		}

	// AltString array
	case lang != "":
		arr := node.Nodes.FindNodeByName("Alt")
		if arr == nil && flags&(CREATE|APPEND) == 0 {
			return fmt.Errorf("CREATE flag required to extend array '%s'", name)
		}
		if arr == nil {
			arr = node.AddNode(NewNode(NewName("rdf:Alt")))
		}
		switch {
		case flags&UNIQUE > 0 && value != "":
			// append source when not exist
			for _, v := range arr.Nodes {
				if attr := v.GetAttr("", "lang"); len(attr) > 0 {
					if attr[0].Value == lang && v.Value == value {
						return nil
					}
				}
			}
			li := arr.AppendNode(NewNode(NewName("rdf:li")))
			li.AddStringAttr("xml:lang", lang)
			li.Value = value

		case flags&APPEND > 0 && value != "":
			// append source value
			li := arr.AppendNode(NewNode(NewName("rdf:li")))
			li.AddStringAttr("xml:lang", lang)
			li.Value = value

		case flags&(REPLACE|CREATE) > 0 && value != "" && lang != "":
			// replace single entry
			for _, v := range arr.Nodes {
				if attr := v.GetAttr("", "lang"); len(attr) > 0 {
					if attr[0].Value == lang {
						v.Value = value
						return nil
					}
				}
			}

		case flags&(REPLACE|CREATE) > 0 && value != "" && lang == "":
			// replace entire AltString with a new version
			arr.Clear()
			li := NewNode(NewName("rdf:li"))
			li.AddStringAttr("xml:lang", lang)
			li.Value = value
			arr.Nodes = NodeList{li}

		case flags&DELETE > 0 && value == "" && lang != "":
			// delete a specific language
			for _, v := range arr.Nodes {
				if attr := v.GetAttr("", "lang"); len(attr) > 0 {
					if attr[0].Value == lang {
						arr.RemoveNode(v).Close()
						return nil
					}
				}
			}

		case flags&DELETE > 0 && value == "" && lang == "":
			// remove and close the entire array
			node.RemoveNode(arr).Close()
		default:
			return fmt.Errorf("unsupported flag combination %v", flags)
		}

	default:
		if path.Len() > 0 {
			return node.SetPath(path, value, flags)
		}
		// fmt.Printf("Set Node path ns=%s len=%d, path=%s, name=%s\n", path.NamespacePrefix(), path.Len(), path.String(), name)
		switch {
		case flags&(REPLACE|CREATE) > 0 && value != "":
			node.Value = value
			return nil
		case flags&DELETE > 0 && value == "":
			node.Value = value
			node.Clear()
			// will be ignored on next marshal
			return nil
		default:
			return fmt.Errorf("unsupported flag combination %v", flags)
		}
	}

	return nil
}

func (n *Node) ListPaths(path Path) (PathValueList, error) {
	l := make(PathValueList, 0)
	switch n.FullName() {
	case "rdf:Seq", "rdf:Bag":
		for i, li := range n.Nodes {
			_, walker := path.Pop()
			walker = walker.AppendIndex(i)
			for _, v := range li.Nodes {
				name := v.Name()
				if v.Namespace() != path.NamespacePrefix() {
					name = v.FullName()
				}
				r, err := v.ListPaths(walker.Push(name))
				if err != nil {
					return nil, err
				}
				l = append(l, r...)
			}
		}
	case "rdf:Alt":
		for _, li := range n.Nodes {
			lang := "x-default"
			if attr := li.GetAttr("", "lang"); len(attr) > 0 {
				lang = attr[0].Value
			}
			_, walker := path.Pop()
			walker = walker.AppendIndexString(lang)
			l = append(l, PathValue{
				Path:  walker,
				Value: li.Value,
			})
		}
	default:
		for _, a := range n.Attr {
			if skipField(a.Name) {
				continue
			}
			name := a.Name.Local
			if hasPrefix(name) && getPrefix(name) == path.NamespacePrefix() {
				name = stripPrefix(name)
			}
			l = append(l, PathValue{
				Path:  path.Push(name),
				Value: a.Value,
			})
		}
		for _, v := range n.Nodes {
			if skipField(v.XMLName) {
				continue
			}
			name := v.Name()
			if v.Namespace() != path.NamespacePrefix() {
				name = v.FullName()
			}
			r, err := v.ListPaths(path.Push(name))
			if err != nil {
				return nil, err
			}
			l = append(l, r...)
		}
		if n.Value != "" {
			l = append(l, PathValue{
				Path:  path,
				Value: n.Value,
			})
		}
	}
	sort.Sort(byPath(l))
	return l, nil
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
