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
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type Marshaler interface {
	MarshalXMP(e *Encoder, n *Node, model Model) error
}

type MarshalerAttr interface {
	MarshalXMPAttr(e *Encoder, name xml.Name, n *Node) (Attr, error)
}

const (
	Xpacket = 1 << iota
	Xpadding
	eMode = Xpacket | Xpadding
)

type Encoder struct {
	e        *xml.Encoder
	cw       *countWriter
	root     *Node
	version  Version
	nsTagMap map[string]string
	intNsMap map[string]*Namespace
	extNsMap map[string]*Namespace
	flags    int
}

var ErrOverflow = errors.New("xmp: document exceeds size limit")

type countWriter struct {
	n     int64
	limit int64
	w     io.Writer
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func newCountWriter(w io.Writer) *countWriter {
	return &countWriter{w: w}
}

func (w *countWriter) Write(b []byte) (int, error) {
	if w.limit > 0 {
		if w.n >= w.limit {
			return 0, ErrOverflow
		}
		bLen := int64(len(b))
		copied, err := io.CopyN(w.w, bytes.NewReader(b), min(w.limit-w.n, bLen))
		w.n += copied
		if err != nil {
			return int(copied), err
		}
		if copied < bLen {
			return int(copied), ErrOverflow
		}
		return int(copied), nil
	} else {
		n, err := w.w.Write(b)
		w.n += int64(n)
		return n, err
	}
}

func NewEncoder(w io.Writer) *Encoder {
	cw := newCountWriter(w)
	return &Encoder{
		e:        xml.NewEncoder(cw),
		cw:       cw,
		root:     NewNode(xml.Name{}),
		nsTagMap: make(map[string]string),
		intNsMap: make(map[string]*Namespace),
		extNsMap: make(map[string]*Namespace),
		flags:    Xpacket,
	}
}

func (e *Encoder) SetMaxSize(size int64) {
	e.cw.limit = size
}

func (e *Encoder) SetFlags(flags int) {
	e.flags = flags & eMode
}

func (e *Encoder) Indent(prefix, indent string) {
	e.e.Indent(prefix, indent)
}

func Marshal(d *Document) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(d); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func MarshalIndent(d *Document, prefix, indent string) ([]byte, error) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.e.Indent(prefix, indent)
	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// encode to separate rdf:Description objects per namespace
func (e *Encoder) Encode(d *Document) error {

	if d == nil {
		return nil
	}

	// sync individual models to establish correct XMP entries
	if err := d.syncToXMP(); err != nil {
		return err
	}

	// reset nodes and map
	e.root = NewNode(xml.Name{})
	defer e.root.Close()
	e.nsTagMap = make(map[string]string)
	e.intNsMap = d.intNsMap
	e.extNsMap = d.extNsMap

	// 1  build output node tree (model -> nodes+attr with one root node per
	//    XMP namespace)
	for _, n := range d.nodes {
		// 1.1  encode the model (Note: models typically use multiple XMP namespaces)
		//      so we generate wrapper nodes on the fly
		if n.Model != nil {
			if err := e.marshalValue(reflect.ValueOf(n.Model), nil, e.root, true); err != nil {
				return err
			}
		}

		// 1.2  merge external nodes (Note: all ext nodes collected under a
		//      document node belong to the same namespace)
		ns := e.findNs(n.XMLName)
		if ns == nil {
			return fmt.Errorf("xmp: missing namespace for model node %s", n.XMLName.Local)
		}
		node := e.root.Nodes.FindNode(ns)
		if node == nil {
			node = NewNode(n.XMLName)
			e.root.AddNode(node)
		}
		node.Nodes = append(node.Nodes, copyNodes(n.Nodes)...)

		// 1.3  merge external attributes (Note: all ext attr collected under a
		//      document node belong to the same namespace)
		node.Attr = append(node.Attr, n.Attr...)
	}

	// 2  collect root-node namespaces
	for _, n := range e.root.Nodes {
		l := make([]Attr, 0)
		for _, ns := range n.Namespaces(d) {
			l = append(l, ns.GetAttr())
		}

		// add the about attr
		about := aboutAttr
		about.Value = d.about
		l = append(l, about)
		n.Attr = append(l, n.Attr...)
		n.XMLName = rdfDescription
	}

	// 3 remove empty root nodes
	nl := make(NodeList, 0)
	for _, n := range e.root.Nodes {
		if len(n.Attr) <= 2 && len(n.Nodes) == 0 {
			n.Close()
			continue
		}
		nl = append(nl, n)
	}
	e.root.Nodes = nl

	// 4  output XML

	// 4.1 write packet header
	if e.flags&Xpacket > 0 {
		if _, err := e.cw.Write(xmp_packet_header); err != nil {
			return err
		}
	}

	// 4.2 add top-level XMP namespace and toolkit as attributes
	start := xml.StartElement{
		Name: xml.Name{Local: "x:xmpmeta"},
		Attr: make([]xml.Attr, 0),
	}
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:x"}, Value: "adobe:ns:meta/"})
	tk := d.toolkit
	if tk == "" {
		tk = XMP_TOOLKIT_VERSION
	}
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "x:xmptk"}, Value: tk})
	if err := e.e.EncodeToken(start); err != nil {
		return err
	}

	// 4.3 add document tree
	e.root.XMLName = xml.Name{Local: "rdf:RDF"}
	e.root.Attr = []Attr{nsRDF.GetAttr()}
	if err := e.e.Encode(e.root); err != nil {
		return err
	}

	if err := e.e.EncodeToken(start.End()); err != nil {
		return err
	}
	if err := e.e.Flush(); err != nil {
		return err
	}

	// 4.4 write footer including optional padding
	if e.flags&Xpacket > 0 {
		if e.flags&Xpadding > 0 && e.cw.limit > 0 {
			pad := e.cw.limit - e.cw.n - 20

			for i := int64(0); i < pad; i++ {
				if i%80 == 0 {
					if _, err := e.cw.Write([]byte("\n")); err != nil {
						return err
					}
				} else {
					if _, err := e.cw.Write([]byte(" ")); err != nil {
						return err
					}
				}
			}
		}

		if _, err := e.cw.Write(xmp_packet_footer); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) EncodeElement(v interface{}, node *Node) error {
	return e.marshalValue(reflect.ValueOf(v), nil, node, false)
}

func (e *Encoder) marshalValue(val reflect.Value, finfo *fieldInfo, node *Node, withWrapper bool) error {
	// name := "-"
	// if finfo != nil {
	// 	name = finfo.name
	// }
	// log.Debugf("xmp: marshalValue %v %v %s\n", val.Kind(), val.Type(), name)

	if !val.IsValid() {
		return nil
	}

	// check empty
	if finfo != nil && finfo.flags&fEmpty == 0 && isEmptyValue(val) {
		// log.Debugf("xmp: skipping empty value for %s\n", name)
		return nil
	}

	// Drill into interfaces and pointers.
	// This can turn into an infinite loop given a cyclic chain,
	// but it matches the Go 1 behavior.
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			// log.Debugf("xmp: skipping nil value for %s\n", name)
			return nil
		}
		val = val.Elem()
	}

	kind := val.Kind()
	typ := val.Type()

	// call marshaler interface
	if val.CanInterface() && (finfo != nil && finfo.flags&fMarshal > 0 || typ.Implements(marshalerType)) {
		// log.Debugf("xmp: marshalValue calling MarshalXMP on %v for %s\n", val.Type(), name)
		return val.Interface().(Marshaler).MarshalXMP(e, node, nil)
	}

	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fMarshal > 0 || pv.Type().Implements(marshalerType)) {
			// log.Debugf("xmp: marshalValue calling MarshalXMP on %v for %s\n", val.Type(), name)
			return pv.Interface().(Marshaler).MarshalXMP(e, node, nil)
		}
	}

	// Check for text marshaler and marshal as node value
	if val.CanInterface() && (finfo != nil && finfo.flags&fTextMarshal > 0 || typ.Implements(textMarshalerType)) {
		// log.Debugf("xmp: marshalValue calling MarshalText on %v for %s\n", val.Type(), name)
		b, err := val.Interface().(encoding.TextMarshaler).MarshalText()
		if err != nil || b == nil {
			return err
		}
		node.Value = string(b)
		return nil
	}

	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && (finfo != nil && finfo.flags&fTextMarshal > 0 || pv.Type().Implements(textMarshalerType)) {
			// log.Debugf("xmp: marshalValue calling MarshalText on %v for %s\n", val.Type(), name)
			b, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil || b == nil {
				return err
			}
			node.Value = string(b)
			return nil
		}
	}

	// XMP arrays require special treatment. Most arrays should have MarshalXML
	// methods defined, but in case an extenstion developer forgets, we add this
	// transparently, guessing the correct RDF array type from Go array/slice.
	if (kind == reflect.Slice || kind == reflect.Array) && typ.Elem().Kind() != reflect.Uint8 {
		atype := ArrayTypeOrdered
		if kind == reflect.Array {
			atype = ArrayTypeUnordered
		}
		// log.Debugf("xmp: marshalValue calling MarshalArray on %v for %s\n", typ, name)
		return MarshalArray(e, node, atype, val.Interface())
	}

	// simple values are just fine
	if node != nil && kind != reflect.Struct {
		// log.Debugf("xmp: marshalValue adding simple value to node %s type %v\n", node.XMLName.Local, typ)
		s, b, err := marshalSimple(typ, val)
		if err != nil {
			return err
		}
		if b != nil {
			s = string(b)
		}
		node.Value = s
		return nil
	}

	// handle structs
	tinfo, err := getTypeInfo(typ, "xmp")
	if err != nil {
		return err
	}

	// encode struct attributes
	for _, finfo := range tinfo.fields {
		if finfo.flags&fOmit > 0 {
			continue
		}

		if finfo.flags&fAttr == 0 {
			continue
		}

		// version must always match
		if !e.version.Between(finfo.minVersion, finfo.maxVersion) {
			// log.Debugf("xmp: marshalValue attr field %s version %v - %v does not match %v\n", finfo.name, finfo.minVersion, finfo.maxVersion, e.version)
			continue
		}

		fv := finfo.value(val)

		if (fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr) && fv.IsNil() {
			// log.Debugf("xmp: marshalValue attr field %s is nil\n", finfo.name)
			continue
		}

		if finfo.flags&fEmpty == 0 && isEmptyValue(fv) {
			// log.Debugf("xmp: marshalValue attr field %s is empty\n", finfo.name)
			continue
		}

		// find or create output node for storing the attribute/node contents
		var dest *Node
		if withWrapper || node == nil {
			ns := e.findNs(NewName(finfo.name))
			if ns == nil {
				return fmt.Errorf("xmp: missing namespace for attr field %s", finfo.name)
			}
			dest = node.Nodes.FindNode(ns)
			if dest == nil {
				// log.Debugf("xmp: marshalValue creating new node for attr field %s type %s\n", finfo.name, ns.GetName())
				dest = NewNode(ns.XMLName(""))
				node.AddNode(dest)
			}
		} else {
			dest = node
		}

		name := NewName(finfo.name)
		if err := e.marshalAttr(dest, name, fv); err != nil {
			return err
		}
	}

	// encode struct fields
	var haveField bool
	for _, finfo := range tinfo.fields {
		if finfo.flags&fOmit > 0 {
			continue
		}

		if finfo.flags&fElement == 0 {
			// log.Debugf("xmp: marshalValue field %s is not an element\n", finfo.name)
			continue
		}

		// version must always match
		if !e.version.Between(finfo.minVersion, finfo.maxVersion) {
			// log.Debugf("xmp: marshalValue node field %s version %v - %v does not match %v\n", finfo.name, finfo.minVersion, finfo.maxVersion, e.version)
			continue
		}

		fv := finfo.value(val)

		if (fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr) && fv.IsNil() {
			// log.Debugf("xmp: marshalValue node field %s is nil\n", finfo.name)
			continue
		}

		if finfo.flags&fEmpty == 0 && isEmptyValue(fv) {
			// log.Debugf("xmp: marshalValue node field %s is empty\n", finfo.name)
			continue
		}

		// find or create output node for storing the attribute/node contents
		var dest *Node
		if withWrapper || node == nil {
			ns := e.findNs(NewName(finfo.name))
			if ns == nil {
				return fmt.Errorf("xmp: marshalValue missing namespace for %s", finfo.name)
			}
			dest = node.Nodes.FindNode(ns)
			if dest == nil {
				// log.Debugf("xmp: marshalValue creating new node for node field %s type %s\n", finfo.name, ns.GetName())
				dest = NewNode(ns.XMLName(""))
				node.AddNode(dest)
			}
		} else {
			dest = node
		}

		// create a new node for this resource
		resNode := NewNode(NewName(finfo.name))
		dest.AddNode(resNode)

		if err := e.marshalValue(fv, &finfo, resNode, false); err != nil {
			return err
		}
		haveField = true
	}

	if haveField && node.XMLName != rdfDescription {
		// for all kinds of structs add rdf:resource as node attribute
		node.AddAttr(rdfResourceAttr)
	}

	return nil
}

func isEmptyValue(v reflect.Value) bool {
	if v.CanInterface() && v.Type().Implements(zeroType) {
		return v.Interface().(Zero).IsZero()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (e *Encoder) marshalAttr(node *Node, name xml.Name, val reflect.Value) error {
	if val.CanInterface() && val.Type().Implements(attrMarshalerType) {
		// log.Debugf("xmp: marshalAttr calling MarshalXmpAttr on %v for %s\n", val.Type(), name.Local)
		attr, err := val.Interface().(MarshalerAttr).MarshalXMPAttr(e, name, node)
		if err != nil {
			return err
		}
		if attr.Name.Local != "" {
			node.AddAttr(attr)
		}
		return nil
	}

	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(attrMarshalerType) {
			// log.Debugf("xmp: marshalAttr calling MarshalXmpAttr on %v for %s\n", val.Type(), name.Local)
			attr, err := pv.Interface().(MarshalerAttr).MarshalXMPAttr(e, name, node)
			if err != nil {
				return err
			}
			if attr.Name.Local != "" {
				node.AddAttr(attr)
			}
			return nil
		}
	}

	if val.CanInterface() && val.Type().Implements(textMarshalerType) {
		// log.Debugf("xmp: marshalAttr calling MarshalText on %v for %s\n", val.Type(), name.Local)
		b, err := val.Interface().(encoding.TextMarshaler).MarshalText()
		if err != nil || b == nil {
			return err
		}
		node.AddAttr(Attr{Name: name, Value: string(b)})
		return nil
	}

	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(textMarshalerType) {
			// log.Debugf("xmp: marshalAttr calling MarshalText on %v for %s\n", val.Type(), name.Local)
			b, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil || b == nil {
				return err
			}
			node.AddAttr(Attr{Name: name, Value: string(b)})
			return nil
		}
	}

	// Dereference or skip nil pointer, interface values.
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			// log.Debugf("xmp: marshalAttr field %s is nil\n", name.Local)
			return nil
		}
		val = val.Elem()
	}

	s, b, err := marshalSimple(val.Type(), val)
	if err != nil {
		return err
	}
	if b != nil {
		s = string(b)
	}
	node.AddAttr(Attr{Name: name, Value: s})
	return nil
}

func marshalSimple(typ reflect.Type, val reflect.Value) (string, []byte, error) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil, nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil, nil
	case reflect.String:
		return val.String(), nil, nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil, nil
	case reflect.Array:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// [...]byte
		var bytes []byte
		if val.CanAddr() {
			bytes = val.Slice(0, val.Len()).Bytes()
		} else {
			bytes = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(bytes), val)
		}
		return "", bytes, nil
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// []byte
		return "", val.Bytes(), nil
	}
	return "", nil, fmt.Errorf("xmp: no method for marshalling type %s (%v)\n", typ.String(), val.Kind())
}

func (e Encoder) _findNsByURI(uri string) *Namespace {
	if v, ok := e.intNsMap[uri]; ok {
		return v
	}
	if v, ok := e.extNsMap[uri]; ok {
		return v
	}
	return nil
}

func (e Encoder) _findNsByPrefix(pre string) *Namespace {
	for _, v := range e.intNsMap {
		if v.GetName() == pre {
			return v
		}
	}
	for _, v := range e.extNsMap {
		if v.GetName() == pre {
			return v
		}
	}
	if ns, err := NsRegistry.GetNamespace(pre); err == nil {
		return ns
	}
	return nil
}

func (e Encoder) findNs(n xml.Name) *Namespace {
	var ns *Namespace
	if n.Space != "" {
		ns = e._findNsByURI(n.Space)
	}
	if ns == nil {
		ns = e._findNsByPrefix(getPrefix(n.Local))
	}
	return ns
}

func ToString(t interface{}) string {
	val := reflect.Indirect(reflect.ValueOf(t))
	if !val.IsValid() {
		return ""
	}
	if val.Type().Implements(stringerType) {
		return t.(fmt.Stringer).String()
	}
	if s, err := ToRawString(val.Interface()); err == nil {
		return s
	}
	return fmt.Sprintf("%v", val.Interface())
}

func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func ToRawString(t interface{}) (string, error) {
	val := reflect.Indirect(reflect.ValueOf(t))
	typ := val.Type()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Array:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// [...]byte
		var b []byte
		if val.CanAddr() {
			b = val.Slice(0, val.Len()).Bytes()
		} else {
			b = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(b), val)
		}
		if !isBase64(string(b)) {
			return base64.StdEncoding.EncodeToString(b), nil
		}
		return string(b), nil
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// []byte
		b := val.Bytes()
		if !isBase64(string(b)) {
			return base64.StdEncoding.EncodeToString(b), nil
		}
		return string(b), nil
	}
	return "", fmt.Errorf("xmp: no method for converting type %s (%v) to string", typ.String(), val.Kind())
}
