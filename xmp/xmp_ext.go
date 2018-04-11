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
)

// Raw nodes are used to keep custom data models inside a container. This
// is used by extension nodes in the Adobe UMC SDK, but may be used in
// private models as well.
type Extension Node

func (x *Extension) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return (*Node)(x).MarshalXML(e, start)
}

func (x *Extension) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return (*Node)(x).UnmarshalXML(d, start)
}

func (x *Extension) IsZero() bool {
	return (*Node)(x).IsZero()
}

func (x *Extension) FindNodeByName(name string) *Node {
	for _, n := range x.Nodes {
		if n.Name() == name {
			return n
		}
		if n.Model != nil && n.Model.Can(name) {
			return n
		}
	}
	return nil
}

func (x Extension) MarshalXMP(e *Encoder, node *Node, m Model) error {
	for _, ext := range x.Nodes {
		if ext.Model != nil {
			if err := e.EncodeElement(ext.Model, node); err != nil {
				return err
			}
		}

		// unwrap one extra layer that we keep for identifying extensions
		if x.XMLName.Local == "" {
			for _, child := range ext.Nodes {
				if child.Model != nil {
					if err := e.EncodeElement(child.Model, node); err != nil {
						return err
					}
				} else {
					node.Nodes = append(node.Nodes, copyNode(child))
				}
			}
		} else {
			node.Nodes = append(node.Nodes, copyNode(ext))
		}
	}
	return nil
}

func (x *Extension) UnmarshalXMPAttr(d *Decoder, a Attr) error {
	return d.decodeAttribute(&x.Nodes, a)
}

func (x *Extension) UnmarshalXMP(d *Decoder, src *Node, m Model) error {
	for _, v := range src.Attr {
		if err := d.decodeAttribute(&x.Nodes, v); err != nil {
			return err
		}
	}

	for _, child := range src.Nodes {
		// recurse if we see a new description element (used in xmpMM:pantry)
		if child.FullName() == "rdf:Description" {
			if err := x.UnmarshalXMP(d, child, m); err != nil {
				return err
			}
			continue
		}
		if err := d.decodeNode(&x.Nodes, child); err != nil {
			// when decoding fails, store as raw child node
			x.Nodes = append(x.Nodes, copyNode(child))
		}
	}
	return nil
}

// used for xmpMM:Pantry bags
//
// <xmpMM:Pantry>
//   <rdf:Bag>
//     <rdf:li>
//       <rdf:Description>...</rdf:Description>
//     </rdf:li>
//   </rdf:Bag>
// </xmpMM:Pantry>
//
type ExtensionArray []*Extension

func (x ExtensionArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func (x ExtensionArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	if len(x) == 0 {
		return nil
	}
	bag := NewNode(xml.Name{Local: "rdf:Bag"})
	node.Nodes = append(node.Nodes, bag)
	for _, v := range x {
		li := NewNode(xml.Name{Local: "rdf:li"})
		elem := NewNode(rdfDescription)
		li.Nodes = append(li.Nodes, elem)
		if err := e.EncodeElement(v, elem); err != nil {
			return err
		}
		bag.Nodes = append(bag.Nodes, li)
	}
	return nil
}

func (x *ExtensionArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	return UnmarshalArray(d, node, x.Typ(), x)
}

// Generates the following XMP/JSON structures from an array instead of
// xmp:<rdf:Bag> / json:[]
//
// <iXML:extension>
//    <PRIVATE-NAME-1 rdf:parseType="Resource">
//        <PRIVATE-FIELD/>
//    </PRIVATE-NAME-1>
//    <PRIVATE-NAME-2 rdf:parseType="Resource">
//        <PRIVATE-FIELD/>
//    </PRIVATE-NAME-2>
// </iXML:extension>
//
// iXML:extension: {
//		"PRIVATE-NAME-1": {
//	        "PRIVATE-FIELD": "",
//      },
//		"PRIVATE-NAME-2": {
//	        "PRIVATE-FIELD": "",
//      }
// }
//
type NamedExtensionArray []*Extension

func (x NamedExtensionArray) Typ() ArrayType {
	return ArrayTypeUnordered
}

func (x *NamedExtensionArray) FindNodeByName(name string) *Node {
	for _, v := range *x {
		n := (*Node)(v)
		if n.Name() == name {
			return n
		}
		if n.Model != nil && n.Model.Can(name) {
			return n
		}
	}
	return nil
}

func (x NamedExtensionArray) MarshalXMP(e *Encoder, node *Node, m Model) error {
	for _, v := range x {
		ext := NewNode(NewName(v.XMLName.Local))
		ext.AddAttr(rdfResourceAttr)
		ext.Nodes = append(ext.Nodes, copyNodes(v.Nodes)...)
		ext.Attr = append(ext.Attr, v.Attr...)
		node.AddNode(ext)
	}
	return nil
}

func (x *NamedExtensionArray) UnmarshalXMP(d *Decoder, node *Node, m Model) error {
	for _, v := range node.Nodes {
		v.translate(d)
		ext := (*Extension)(NewNode(NewName(v.XMLName.Local)))
		ext.Nodes = append(ext.Nodes, copyNodes(v.Nodes)...)
		ext.Attr = append(ext.Attr, v.Attr...)
		*x = append(*x, ext)
	}
	return nil
}
