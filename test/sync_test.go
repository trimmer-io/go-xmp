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

package main

import (
	"os"
	"testing"

	_ "github.com/trimmer-io/go-xmp/models"
	"github.com/trimmer-io/go-xmp/xmp"
)

var SyncMergeTestcases = map[xmp.Path]string{
	xmp.Path("dc:format"):       "format",    // string
	xmp.Path("dc:title"):        "the_title", // AltString
	xmp.Path("dc:type"):         "type",      // StringArray
	xmp.Path("myns:type"):       "hello",
	xmp.Path("myns:root/child"): "hello2",
}

func TestMergeCreate(T *testing.T) {
	d1 := xmp.NewDocument()
	for n, v := range SyncMergeTestcases {
		if err := d1.SetPath(xmp.PathValue{
			Path:      n,
			Value:     v,
			Namespace: "exotic:ns:42",
			Flags:     xmp.CREATE,
		}); err != nil {
			T.Errorf("%v: %v", n, err)
		}
	}
	d2 := xmp.NewDocument()
	if err := d2.Merge(d1, xmp.MERGE); err != nil {
		T.Errorf("merge failed: %v", err)
	}
	for n, v := range SyncMergeTestcases {
		if val, err := d2.GetPath(n); err != nil {
			T.Errorf("%v: get failed: %v", n, err)
		} else if v != val {
			T.Errorf("%v: invalid contents, expected=%s got=%s", n, v, val)
		}
	}
}

func TestUpdateAfterMerge(T *testing.T) {
	d1 := xmp.NewDocument()
	for n, v := range SyncMergeTestcases {
		if err := d1.SetPath(xmp.PathValue{
			Path:      n,
			Value:     v,
			Namespace: "exotic:ns:42",
			Flags:     xmp.CREATE,
		}); err != nil {
			T.Errorf("%v: %v", n, err)
		}
	}
	d2 := xmp.NewDocument()
	if err := d2.Merge(d1, xmp.MERGE); err != nil {
		T.Errorf("merge failed: %v", err)
	}
	for n, v := range SyncMergeTestcases {
		if err := d1.SetPath(xmp.PathValue{
			Path:  n,
			Value: v + "_updated",
			Flags: xmp.REPLACE,
		}); err != nil {
			T.Errorf("%v: set failed: %v", n, err)
		}
		if val, err := d2.GetPath(n); err != nil {
			T.Errorf("%v: get failed: %v", n, err)
		} else if v != val {
			T.Errorf("%v: invalid contents, expected=%s got=%s", n, v+"_updated", val)
		}
	}
}

func TestMergeDocuments(T *testing.T) {
	for _, v := range testfiles {
		f, err := os.Open(v)
		if err != nil {
			T.Logf("Cannot open sample '%s': %v", v, err)
			continue
		}
		dec := xmp.NewDecoder(f)
		d1 := xmp.NewDocument() //&xmp.Document{}
		if err := dec.Decode(d1); err != nil {
			T.Errorf("%s: %v", v, err)
		}
		f.Close()
		// set the model dirty so ListPaths will call SyncToXmp to make sure all
		// cross-model relations are in sync; this is required to get the same result
		// compared to d2 (after merge) which will also run ListPaths on a dirty document
		d1.SetDirty()
		p1, err := d1.ListPaths()
		if err != nil {
			T.Errorf("%s: %v", v, err)
			continue
		}
		d2 := xmp.NewDocument()
		if err := d2.Merge(d1, xmp.MERGE); err != nil {
			T.Errorf("merge failed for '%s': %v", v, err)
			continue
		}
		p2, err := d2.ListPaths()
		if err != nil {
			T.Errorf("%s: %v", v, err)
			continue
		}
		// compare all paths
		if l1, l2 := len(p1), len(p2); l1 != l2 {
			T.Errorf("merge '%s': invalid number of paths: exp=%d got=%d", v, l1, l2)
			for _, v := range p1.Diff(p2) {
				var op byte
				switch v.Flags {
				case xmp.DELETE:
					op = '-'
				case xmp.CREATE:
					op = '+'
				case xmp.REPLACE:
					op = '#'
				}
				T.Errorf("> %c %s '%s'", op, v.Path.String(), v.Value)
			}
			continue
		}
		for i, x := range p1 {
			if x.Value != p2[i].Value {
				T.Errorf("merge '%s': value mismatch on path '%s': exp=%s got=%s", v, x.Path.String(), x.Value, p2[i].Value)
			}
		}
	}
}
