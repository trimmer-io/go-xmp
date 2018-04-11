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
	"strings"
)

func (n Node) Dump(d int) {
	pfx := strings.Repeat(">", d)
	if len(pfx) > 0 {
		pfx += " "
	}
	if n.Model != nil {
		fmt.Printf("%sNODE MODEL %s (%d attr, %d children)\n", pfx, n.XMLName.Local, len(n.Attr), len(n.Nodes))
	} else {
		fmt.Printf("%sNODE EXT  %s (%d attr, %d children)\n", pfx, n.XMLName.Local, len(n.Attr), len(n.Nodes))
	}
	for i, v := range n.Attr {
		fmt.Printf("%s ATTR %2d %s := %s\n", pfx, i, v.Name.Local, v.Value)
	}
	for i, v := range n.Nodes {
		fmt.Printf("%d ", i)
		v.Dump(d + 1)
	}
}

func (d *Document) Dump() {
	d.DumpNamespaces()
	p, _ := d.ListPaths()
	for _, v := range p {
		fmt.Printf("%s = %s\n", v.Path.String(), v.Value)
	}
	// for _, v := range d.Nodes {
	// 	v.Dump(0)
	// }
}

func (d *Document) DumpNamespaces() {
	for n, v := range d.intNsMap {
		fmt.Printf("NS INT %s %+v\n", n, v)
	}
	for n, v := range d.extNsMap {
		fmt.Printf("NS EXT %s %+v\n", n, v)
	}
}

func DumpStats() {
	fmt.Println("xmp: Node Pool Allocs ", npAllocs)
	fmt.Println("xmp: Node Pool Frees  ", npFrees)
	fmt.Println("xmp: Node Pool Hits   ", npHits)
	fmt.Println("xmp: Node Pool Returns", npReturns)
	fmt.Println("xmp: Node Pool InUse  ", npAllocs+npHits-npReturns-npFrees)
	fmt.Println("xmp: Node Pool InPool ", Max64(0, npReturns-npHits))
}
