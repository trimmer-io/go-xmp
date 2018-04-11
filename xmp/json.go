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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
)

type jsonDocument struct {
	About      string                     `json:"about,omitempty"`
	Toolkit    string                     `json:"toolkit,omitempty"`
	Namespaces map[string]string          `json:"namespaces"`
	Models     map[string]json.RawMessage `json:"models"`
}

type jsonOutDocument struct {
	About      string                 `json:"about,omitempty"`
	Toolkit    string                 `json:"toolkit,omitempty"`
	Namespaces map[string]string      `json:"namespaces"`
	Models     map[string]interface{} `json:"models"`
}

func (d *Document) MarshalJSON() ([]byte, error) {
	// sync individual models to establish correct XMP entries
	if err := d.syncToXMP(); err != nil {
		return nil, err
	}

	out := &jsonOutDocument{
		About:      d.about,
		Toolkit:    d.toolkit,
		Namespaces: make(map[string]string),
		Models:     make(map[string]interface{}),
	}

	if out.Toolkit == "" {
		out.Toolkit = XMP_TOOLKIT_VERSION
	}

	// We're using the regular XMP decoder with a JSON boilerplate.
	e := NewEncoder(nil)
	e.intNsMap = d.intNsMap
	e.extNsMap = d.extNsMap
	defer e.root.Close()

	// 1  build output node tree (model -> nodes+attr with one root node per
	//    XMP namespace)
	for _, n := range d.nodes {
		// 1.1  encode the model (Note: models typically use multiple XMP namespaces)
		//      so we generate wrapper nodes on the fly
		if n.Model != nil {
			if err := e.marshalValue(reflect.ValueOf(n.Model), nil, e.root, true); err != nil {
				return nil, err
			}
		}

		// 1.2  merge external nodes (Note: all ext nodes collected under a
		//      document node belong to the same namespace)
		ns := e.findNs(n.XMLName)
		if ns == nil {
			return nil, fmt.Errorf("xmp: missing namespace for model node %s\n", n.XMLName.Local)
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
		for _, v := range n.Namespaces(d) {
			if v == nsX || v == nsXML || v == nsRDF {
				continue
			}
			out.Namespaces[v.GetName()] = v.GetURI()
		}
	}

	// 3 convert node tree to json, ignore empty root nodes
	for _, n := range e.root.Nodes {
		if n.IsZero() {
			continue
		}
		if m, err := nodeToJson(n); err != nil {
			return nil, err
		} else {
			// add model under it's namespace name
			out.Models[n.Namespace()] = m
		}
	}

	b, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	// 4 output as json
	return b, nil
}

func nodeToJson(n *Node) (interface{}, error) {
	// process value leaf nodes (Note: value and model are mutually exclusive)
	if n.Value != "" {
		return &n.Value, nil
	}

	out := make(map[string]interface{})

	// add node attributes
	for _, v := range n.Attr {
		name := v.Name.Local
		// Note:  stripping rdf attributes from external child nodes
		//        that are relinked during marshal with all attributes
		if name == "rdf:parseType" {
			continue
		}
		out[name] = v.Value
	}

	// add node children
	for _, v := range n.Nodes {
		switch v.FullName() {
		case "rdf:Bag", "rdf:Seq":
			// skip outer node and unwrap inner li nodes as array elements
			a := make([]interface{}, 0)
			for _, vv := range v.Nodes {
				if vv.Value != "" {
					// simple lists (i.e. string lists or numbers etc)
					a = append(a, &vv.Value)
				} else {
					// object lists
					// m := make(map[string]interface{})
					if m, err := nodeToJson(vv); err != nil {
						return nil, err
					} else {
						a = append(a, m)
					}
				}
			}
			return a, nil

		case "rdf:Alt":
			// determine content type (strings or objects)
			if len(v.Nodes) > 0 && len(v.Nodes[0].Nodes) > 0 {
				// object type: skip outer node and insert interface
				a := make([]interface{}, 0)
				for _, vv := range v.Nodes {
					if m, err := nodeToJson(vv); err != nil {
						return nil, err
					} else {
						a = append(a, m)
					}
				}
				return a, nil
			} else {
				// string type: skip outer node and insert ArrayItems
				var defLangContent string
				var defLangIndex int = -1
				a := make([]*AltItem, 0)
				for i, vv := range v.Nodes {
					// string arrays
					ai := &AltItem{
						Value: vv.Value,
					}
					langAttr := vv.GetAttr("", "lang")
					if len(langAttr) > 0 {
						ai.Lang = langAttr[0].Value
					}
					if ai.Lang == "x-default" {
						ai.IsDefault = true
						defLangContent = ai.Value
						defLangIndex = i
					}
					a = append(a, ai)
				}
				// mark the default item
				if defLangContent != "" && len(a) > 1 {
					for i, l := 0, len(a); i < l; i++ {
						if a[i].Value == defLangContent {
							a[i].IsDefault = true
						}
					}
					// remove the duplicate default element
					if defLangIndex > -1 {
						a = append(a[:defLangIndex], a[defLangIndex+1:]...)
					}
				}
				return a, nil
			}

		default:
			if m, err := nodeToJson(v); err != nil {
				return nil, err
			} else {
				out[v.FullName()] = m
			}
		}
	}

	return out, nil
}

// The difference to regular JSON decoding is that we must take care of
// unknown properties that Go's JSON decoder simply ignores. In our XMP
// model we capture everything unknown in child nodes. Note: due to internal
// namespace lookups we cannot use node attributes, because attribute unmarshal
// will fail for all unknown namespaces.
func (d *Document) UnmarshalJSON(data []byte) error {
	in := &jsonDocument{
		Namespaces: make(map[string]string),
		Models:     make(map[string]json.RawMessage),
	}

	if err := json.Unmarshal(data, in); err != nil {
		return fmt.Errorf("xmp: json unmarshal failed: %v", err)
	}

	// We're using the regular XMP decoder with a JSON boilerplate.
	dec := NewDecoder(nil)

	// register namespaces
	dec.about = in.About
	dec.toolkit = in.Toolkit
	for prefix, uri := range in.Namespaces {
		dec.addNamespace(prefix, uri)
	}

	// build node tree from JSON models
	root := NewNode(emptyName)
	defer root.Close()
	for name, b := range in.Models {
		node := NewNode(xml.Name{Local: name})
		root.Nodes = append(root.Nodes, node)
		content := make(map[string]interface{})
		if err := json.Unmarshal(b, &content); err != nil {
			return fmt.Errorf("xmp: json unmarshal model '%s' failed: %v", name, err)
		}
		for n, v := range content {
			jsonToNode(n, v, node)
		}
	}

	// run node tree through xmp unmarshaler
	for _, n := range root.Nodes {
		// process attributes
		for _, v := range n.Attr {
			if err := dec.decodeAttribute(&dec.nodes, v); err != nil {
				return err
			}
		}
		// process child nodes
		for _, v := range n.Nodes {
			if err := dec.decodeNode(&dec.nodes, v); err != nil {
				return err
			}
		}
	}

	// copy decoded values to document
	d.toolkit = dec.toolkit
	d.about = dec.about
	d.nodes = dec.nodes
	d.intNsMap = dec.intNsMap
	d.extNsMap = dec.extNsMap
	return d.syncFromXMP()
}

func jsonToNode(name string, v interface{}, node *Node) {
	switch {
	case isSimpleValue(v):
		var s string
		switch val := v.(type) {
		case string:
			s = val
		case float64:
			s = strconv.FormatFloat(val, 'f', -1, 64)
		case bool:
			s = strconv.FormatBool(val)
		case nil:
			return
		}
		if name != "" && name != "rdf:value" {
			// add simple values as child nodes with string value
			attrNode := NewNode(xml.Name{Local: name})
			attrNode.Value = s
			node.Nodes = append(node.Nodes, attrNode)
		} else {
			// add as node value when no name is given (i.e. used in string arrays)
			node.Value = s
		}

	case isArrayValue(v):
		// arrays of arrays are not supported in XMP
		if name == "" {
			return
		}

		// add arrays as Seq/li or Alt/li child nodes
		containerNode := NewNode(xml.Name{Local: name})
		node.Nodes = append(node.Nodes, containerNode)
		typ := getArrayType(v)
		anode := NewNode(xml.Name{Space: nsRDF.GetURI(), Local: string(typ)})
		containerNode.Nodes = append(containerNode.Nodes, anode)

		for _, av := range v.([]interface{}) {
			linode := NewNode(xml.Name{Space: nsRDF.GetURI(), Local: "li"})
			anode.Nodes = append(anode.Nodes, linode)
			if typ == ArrayTypeAlternative {
				item := av.(map[string]interface{})
				def := item["isDefault"].(bool)
				lang := item["lang"].(string)
				val := item["value"].(string)

				// add single node when default || !default && lang != ""
				if def || !def && lang != "" {
					if def {
						linode.Attr = append(linode.Attr, Attr{
							Name:  xml.Name{Local: "xml:lang"},
							Value: "x-default",
						})
						linode.Value = val
					} else {
						linode.Attr = append(linode.Attr, Attr{
							Name:  xml.Name{Local: "xml:lang"},
							Value: lang,
						})
						linode.Value = val
					}
				}

				// add second node for default && lang != ""
				if def && lang != "" {
					linode = NewNode(xml.Name{Space: nsRDF.GetURI(), Local: "li"})
					anode.Nodes = append(anode.Nodes, linode)
					linode.Attr = append(linode.Attr, Attr{
						Name:  xml.Name{Local: "xml:lang"},
						Value: lang,
					})
					linode.Value = val
				}
			} else {
				jsonToNode("", av, linode)
			}
		}

	case isObjectValue(v):
		// add objects as child nodes (recursive)
		onode := node
		if name != "" {
			onode = NewNode(xml.Name{Local: name})
			node.Nodes = append(node.Nodes, onode)
		}
		for on, ov := range v.(map[string]interface{}) {
			if isArrayItemType(v) && on == "value" {
				jsonToNode("", ov, onode)
			} else {
				jsonToNode(on, ov, onode)
			}
		}
	}
}

// JSON unmarshal helpers
func isSimpleValue(v interface{}) bool {
	switch v.(type) {
	case string, float64, bool, nil:
		return true
	default:
		return false
	}
}

func isObjectValue(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}:
		return true
	default:
		return false
	}
}

func isArrayValue(v interface{}) bool {
	switch v.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

func isArrayItemType(v interface{}) bool {
	// alternative array item
	if isObjectValue(v) {
		val := v.(map[string]interface{})
		if len(val) > 3 {
			return false
		}
		if _, ok := val["value"]; !ok {
			return false
		}
		if _, ok := val["lang"]; !ok {
			return false
		}
		if _, ok := val["isDefault"]; !ok {
			return false
		}
		return true
	}
	return false
}

func getArrayType(v interface{}) ArrayType {
	slice, ok := v.([]interface{})
	if !ok || len(slice) == 0 {
		return ArrayTypeOrdered
	}
	if isArrayItemType(slice[0]) {
		return ArrayTypeAlternative
	}
	return ArrayTypeOrdered
}
