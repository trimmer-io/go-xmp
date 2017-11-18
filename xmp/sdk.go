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

// Package xmp implements the Extensible Metadata Platform (XMP) defined in ISO 16684-1:2011(E).
package xmp // import "trimmer.io/go-xmp/xmp"

import (
	"encoding/xml"
	"strings"
)

const XMP_TOOLKIT_VERSION = "Go XMP SDK 1.0"

var Agent AgentName = XMP_TOOLKIT_VERSION

func SetAgent(org, software, version, extra string) {
	Agent = AgentName(strings.Join([]string{org, software, version, extra}, " "))
}

var (
	// Core XMP namespaces
	nsX   = &Namespace{"x", "adobe:ns:meta/", nil}
	nsXML = &Namespace{"xml", "http://www.w3.org/XML/1998/namespace", nil}
	nsRDF = &Namespace{"rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#", nil}

	// Common Structures
	nsStArea = &Namespace{"stArea", "http://ns.adobe.com/xmp/sType/Area#", nil}

	// private resources
	emptyName       = xml.Name{}
	rdfResourceAttr = Attr{Name: xml.Name{Local: "rdf:parseType"}, Value: "Resource"}
	rdfDescription  = xml.Name{Local: "rdf:Description"}
	aboutAttr       = Attr{Name: xml.Name{Local: "rdf:about"}, Value: ""}
)

var (
	xmp_packet_header    = []byte("<?xpacket begin=\"\" id=\"W5M0MpCehiHzreSzNTczkc9d\"?>\n")
	xmp_packet_footer    = []byte("\n<?xpacket end=\"w\"?>")
	xmp_packet_footer_ro = []byte("\n<?xpacket end=\"r\"?>")
)

func init() {
	for _, v := range []*Namespace{
		nsX,
		nsXML,
		nsRDF,
		nsStArea,
	} {
		NsRegistry.RegisterNamespace(v, nil)
	}
}
