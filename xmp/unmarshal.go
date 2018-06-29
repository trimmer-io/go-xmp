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
	"bytes"
	"encoding"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type Unmarshaler interface {
	UnmarshalXMP(d *Decoder, n *Node, model Model) error
}

type UnmarshalerAttr interface {
	UnmarshalXMPAttr(d *Decoder, a Attr) error
}

type Decoder struct {
	d        *xml.Decoder
	toolkit  string
	about    string
	nodes    NodeList
	intNsMap map[string]*Namespace
	extNsMap map[string]*Namespace
	version  Version
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		d:        xml.NewDecoder(r),
		nodes:    make(NodeList, 0),
		intNsMap: make(map[string]*Namespace),
		extNsMap: make(map[string]*Namespace),
	}
}

func (d *Decoder) SetVersion(v Version) {
	d.version = v
}

func Unmarshal(data []byte, d *Document) error {
	return NewDecoder(bytes.NewReader(data)).Decode(d)
}

func (d *Decoder) DecodeElement(v interface{}, src *Node) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("xmp: model for '%s' is not a pointer", reflect.TypeOf(v))
	}
	return d.unmarshal(val, nil, src)
}

func (d *Decoder) Decode(x *Document) error {

	if x == nil {
		return nil
	}

	// 1  parse node tree from XML
	root := NewNode(emptyName)
	gc := root
	defer gc.Close()
	if err := d.d.Decode(root); err != nil {
		return fmt.Errorf("xmp: parsing xml failed: %v", err)
	}

	// 2  skip top-level `x:xmpmeta` (optional) and `rdf:RDF` nodes
	if root.FullName() == "x:xmpmeta" {
		if a := root.GetAttr(nsX.GetURI(), "xmptk"); len(a) > 0 {
			x.toolkit = strings.TrimSpace(a[0].Value)
		}
		if len(root.Nodes) == 0 {
			return fmt.Errorf("xmp: invalid XML format: missing rdf:RDF node")
		}
		if len(root.Nodes) > 1 {
			return fmt.Errorf("xmp: invalid XML format: too many child nodes in x:xmpmeta")
		}
		root = root.Nodes[0]
	}

	if root.FullName() != "rdf:RDF" {
		return fmt.Errorf("xmp: invalid XML format: missing rdf:RDF node, found %s:%s", root.Namespace(), root.Name())
	}

	// 3  extract document namespaces
	for _, n := range root.Nodes {
		for _, v := range n.GetAttr("xmlns", "") {
			d.addNamespace(v.Name.Local, v.Value)
		}
	}

	// 4  walk node tree and create model instances
	for _, n := range root.Nodes {
		// we expect outer nodes
		if n.FullName() != "rdf:Description" {
			return fmt.Errorf("xmp: invalid XML format: expected rdf:Description node, found %s:%s", n.Namespace(), n.Name())
		}

		// process attributes
		for _, v := range n.Attr {
			if v.Name.Space == nsRDF.GetURI() {
				if v.Name.Local == "about" {
					d.about = v.Value
				}
				continue
			}

			if err := d.decodeAttribute(&d.nodes, v); err != nil {
				return err
			}
		}

		// process child nodes
		for _, v := range n.Nodes {
			if err := d.decodeNode(&d.nodes, v); err != nil {
				return err
			}
		}
	}

	// copy decoded values to document
	x.toolkit = d.toolkit
	x.about = d.about
	x.nodes = d.nodes
	x.intNsMap = d.intNsMap
	x.extNsMap = d.extNsMap
	return x.syncFromXMP()
}

func (d *Decoder) decodeNode(ctx *NodeList, src *Node) error {

	node, err := d.lookupNode(ctx, src.XMLName)
	if err != nil {
		return err
	}

	d.translate(&src.XMLName)
	name := src.FullName()

	// process the node value
	var storeNode bool
	if node.Model != nil {
		finfo, field := d.findStructField(derefIndirect(node.Model), name)
		if field.IsValid() {
			return d.unmarshal(field, finfo, src)
		} else {
			storeNode = finfo == nil || finfo.flags&fOmit == 0
		}
	} else {
		storeNode = true
	}

	// capture the node and its children into the selected node
	if storeNode {
		src.translate(d)
		if !src.IsZero() {
			node.AddNode(copyNode(src))
			Log.Debugf("xmp: missing struct field for %s, saving as external node in %s model", name, node.FullName())
		}
	}
	return nil
}

func (d *Decoder) decodeAttribute(ctx *NodeList, src Attr) error {

	d.translate(&src.Name)
	if skipField(src.Name) {
		return nil
	}

	node, err := d.lookupNode(ctx, src.Name)
	if err != nil {
		return err
	}

	// process the attribute value
	var storeAttr bool
	if node.Model != nil {
		finfo, field := d.findStructField(derefIndirect(node.Model), src.Name.Local)
		if field.IsValid() {
			if err := d.unmarshalAttr(field, finfo, src); err != nil {
				return err
			}
		} else {
			// capture the field as external attribute
			storeAttr = finfo == nil || finfo.flags&fOmit == 0
		}

	} else {
		storeAttr = true
	}

	if storeAttr {
		if src.Value != "" {
			Log.Debugf("xmp: missing struct field for %s, saving as unknown attr in %s model", src.Name.Local, node.FullName())
			node.AddAttr(src)
		}
	}

	return nil
}

func (d *Decoder) unmarshal(val reflect.Value, finfo *fieldInfo, src *Node) error {
	// Load value from interface, but only if the result will be
	// usefully addressable.
	val = derefValue(val)

	// XMP
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fUnmarshal > 0 || pv.Type().Implements(unmarshalerType)) {
			return pv.Interface().(Unmarshaler).UnmarshalXMP(d, src, nil)
		}
	}

	// Text
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fTextUnmarshal > 0 || pv.Type().Implements(textUnmarshalerType)) {
			return pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(src.Value))
		}
	}

	// structs
	if val.Kind() == reflect.Struct {
		// process attributes first
		for _, a := range src.Attr {
			d.translate(&a.Name)
			if skipField(a.Name) {
				continue
			}
			if finfo, field := d.findStructField(val, a.Name.Local); field.IsValid() {
				if err := d.unmarshalAttr(field, finfo, a); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("xmp: unmarshal model %s: field for attr %s not found in type %v", src.FullName(), a.Name.Local, val.Type())
			}
		}

		// recurse into child nodes
		for _, n := range src.Nodes {
			d.translate(&n.XMLName)
			name := n.FullName()
			switch name {
			case "rdf:Description":
				if err := d.unmarshal(val, nil, n); err != nil {
					return err
				}
			default:
				if skipField(n.XMLName) {
					break
				}
				if finfo, field := d.findStructField(val, name); field.IsValid() {
					if err := d.unmarshal(field, finfo, n); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("xmp: unmarshal model %s: struct field %s not found (not stored)", src.FullName(), name)
				}
			}
		}
	} else {
		// otherwise set simple value directly
		if err := setValue(val, src.Value); err != nil {
			return fmt.Errorf("xmp: unmarshal %s: %v", finfo.String(), err)
		}
	}

	return nil
}

func (d *Decoder) unmarshalAttr(val reflect.Value, finfo *fieldInfo, src Attr) error {
	// Load value from interface
	val = derefValue(val)

	// attribute unmarshaler
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fUnmarshalAttr > 0 || pv.Type().Implements(attrUnmarshalerType)) {
			return pv.Interface().(UnmarshalerAttr).UnmarshalXMPAttr(d, src)
		}
	}

	// text unmarshaler
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fTextUnmarshal > 0 || pv.Type().Implements(textUnmarshalerType)) {
			return pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(src.Value))
		}
	}

	// Slice of element values.
	if val.Type().Kind() == reflect.Slice && val.Type().Elem().Kind() != reflect.Uint8 {
		// Grow slice.
		n := val.Len()
		val.Set(reflect.Append(val, reflect.Zero(val.Type().Elem())))

		// Recur to read element into slice.
		if err := d.unmarshalAttr(val.Index(n), nil, src); err != nil {
			val.SetLen(n)
			return fmt.Errorf("xmp: unmarshal %s: %v", finfo.String(), err)
		}
		return nil
	}

	// otherwise set value directly
	return setValue(val, src.Value)
}

// Translate an xml name's namespace to XMP format.
// Since rdf and xml namespaces are not registered in a document,
// we look up those namespaces in our registry. This is necessary
// to transform rdf-related attributes like rdf:about, xml:lang,
// rdf:parseType, etc.
func (d Decoder) translate(n *xml.Name) {
	if len(n.Space) == 0 || n.Space == "xmlns" {
		return
	}
	ns := d.findNs(*n)
	if ns == nil {
		ns, _ = NsRegistry.GetNamespace(NsRegistry.GetPrefix(n.Space))
	}
	if ns != nil {
		n.Space = ""
		n.Local = ns.Expand(n.Local)
	}
}

// Keep track of all used namespaces (registered and unknown).
// In XML a local name's prefix may differ from the XMP standard
// prefix. Even Adobe used to use `xap` instead of `xmp` as prefix
// in documents before the standard was finished.
func (d *Decoder) addNamespace(prefix, uri string) {
	// register known namespaces using their standard prefix
	if ns := NsRegistry.GetPrefix(uri); len(ns) > 0 {
		d.intNsMap[uri], _ = NsRegistry.GetNamespace(ns)
		return
	}

	// keep track of unknown namespaces using their in-document prefix
	if _, ok := d.extNsMap[uri]; !ok {
		d.extNsMap[uri] = &Namespace{prefix, uri, emptyFactory}
	}
}

func (d Decoder) _findNsByURI(uri string) *Namespace {
	if v, ok := d.intNsMap[uri]; ok {
		return v
	}
	if v, ok := d.extNsMap[uri]; ok {
		return v
	}
	return nil
}

func (d Decoder) _findNsByPrefix(pre string) *Namespace {
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
	return nil
}

func (d Decoder) findNs(n xml.Name) *Namespace {
	var ns *Namespace
	if len(n.Space) > 0 {
		ns = d._findNsByURI(n.Space)
	}
	if ns == nil {
		ns = d._findNsByPrefix(getPrefix(n.Local))
	}
	return ns
}

func (d *Decoder) findStructField(val reflect.Value, name string) (*fieldInfo, reflect.Value) {
	typ := val.Type()
	tinfo, err := getTypeInfo(typ, "xmp")
	if err != nil {
		return nil, reflect.Value{}
	}

	var finfo *fieldInfo
	any := -1
	// pick the correct field based on name, flags and version
	for i, v := range tinfo.fields {
		// version must always match
		if !d.version.Between(v.minVersion, v.maxVersion) {
			continue
		}

		// save `any` field in case
		if v.flags&fAny > 0 {
			any = i
		}

		// field name must match
		if v.name != name {
			continue
		}

		finfo = &v
		break
	}

	if finfo == nil && any >= 0 {
		finfo = &tinfo.fields[any]
	}

	// nothing found
	if finfo == nil {
		return nil, reflect.Value{}
	}

	// allocate memory for pointer values in structs
	v := finfo.value(val)
	if v.Type().Kind() == reflect.Ptr && v.IsNil() && v.CanSet() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	return finfo, v
}

func (d *Decoder) lookupNode(ctx *NodeList, name xml.Name) (*Node, error) {
	// check namespace has been registered (i.e. exists in document)
	ns := d.findNs(name)
	if ns == nil {
		return nil, &UnknownNamespaceError{name}
	}

	// pick or create the XMP model for the current namespace
	node := ctx.FindNode(ns)
	if node == nil {
		model := ns.NewModel()
		if model != nil {
			modelNs := model.Namespaces()
			if len(modelNs) == 0 {
				return nil, fmt.Errorf("xmp: model '%v' must implement at least one namespace", reflect.TypeOf(model))
			}
			node = NewNode(modelNs[0].XMLName(""))
			node.Model = model
		} else {
			node = NewNode(ns.XMLName(""))
		}
		*ctx = append(*ctx, node)
	}
	return node, nil
}

func setValue(dst reflect.Value, src string) error {
	dst0 := dst
	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}

	switch dst.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(src, 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(src, 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(src, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetFloat(i)
	case reflect.Bool:
		i, err := strconv.ParseBool(strings.TrimSpace(src))
		if err != nil {
			return err
		}
		dst.SetBool(i)
	case reflect.String:
		dst.SetString(strings.TrimSpace(src))
	case reflect.Slice:
		// make sure it's a byte slice
		if dst.Type().Elem().Kind() == reflect.Uint8 {
			dst.SetBytes([]byte(src))
		}
	default:
		return fmt.Errorf("xmp: no method for unmarshalling type %s", dst0.Type().String())
	}
	return nil
}

func skipField(n xml.Name) bool {
	if n.Space == "xmlns" {
		return true
	}

	switch n.Local {
	case "rdf:parseType", "rdf:type", "xml:lang":
		return true
	default:
		return false
	}
}
