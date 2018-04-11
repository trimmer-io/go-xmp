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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "trimmer.io/go-xmp/models"
	"trimmer.io/go-xmp/models/dc"
	"trimmer.io/go-xmp/models/xmp_base"
	"trimmer.io/go-xmp/xmp"
)

var (
	debug bool
	quiet bool
	fjson bool
	fxmp  bool
	fpath bool
	forig bool
	fall  bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "enable debugging")
	flag.BoolVar(&quiet, "quiet", false, "don't output anything")
	flag.BoolVar(&fjson, "json", false, "enable JSON output")
	flag.BoolVar(&fxmp, "xmp", false, "enable XMP output")
	flag.BoolVar(&fpath, "path", false, "enable XMP/Path output")
	flag.BoolVar(&forig, "orig", false, "enable original XMP output")
	flag.BoolVar(&fall, "all", false, "ouput all embedded xmp documents")
}

func fail(v interface{}) {
	fmt.Printf("Error: %s in file %s\n", v, flag.Arg(0))
	os.Exit(1)
}

func out(b []byte) {
	if quiet {
		return
	}
	fmt.Println(string(b))
}

func marshal(d *xmp.Document) []byte {
	b, err := xmp.MarshalIndent(d, "", "  ")
	if err != nil {
		fail(err)
	}
	return b
}

func unmarshal(v []byte) *xmp.Document {
	d := &xmp.Document{}
	if err := xmp.Unmarshal(v, d); err != nil {
		fail(err)
	}
	return d
}

func main() {
	flag.Parse()

	if debug {
		xmp.SetLogLevel(xmp.LogLevelDebug)
	}

	// output original when no option is selected
	if !fjson && !fxmp && !fpath && !forig && !quiet {
		forig = true
	}

	var b []byte

	if flag.NArg() > 0 {
		filename := flag.Arg(0)
		f, err := os.Open(filename)
		if err != nil {
			fail(err)
		}
		defer f.Close()

		switch filepath.Ext(filename) {
		case ".xmp":
			b, err = ioutil.ReadAll(f)
			if err != nil {
				fail(err)
			}

		default:
			bb, err := xmp.ScanPackets(f)
			if err != nil && err != io.EOF {
				fail(err)
			}
			if err == io.EOF {
				return
			}
			if forig && fall {
				for _, b := range bb {
					out(b)
				}
				return
			}
			b = bb[0]
		}

	} else {
		// fill the document with some info
		s := xmp.NewDocument()
		s.AddModel(&xmpbase.XmpBase{
			CreatorTool: xmp.Agent,
			CreateDate:  xmp.Now(),
			ModifyDate:  xmp.Now(),
			Thumbnails: xmpbase.ThumbnailArray{
				xmpbase.Thumbnail{
					Format: "image/jpeg",
					Width:  10,
					Height: 10,
					Image:  []byte("not a real image"),
				},
			},
		})
		s.AddModel(&dc.DublinCore{
			Format:  "image/jpeg",
			Title:   xmp.NewAltString("demo"),
			Creator: xmp.NewStringList("Alexander Eichhorn"),
			Description: xmp.NewAltString(
				xmp.AltItem{Value: "Go-XMP Demo Model", Lang: "en", IsDefault: true},
				xmp.AltItem{Value: "Go-XMP Beispiel Modell", Lang: "de", IsDefault: false},
			),
		})
		b = marshal(s)
		s.Close()
	}

	if forig {
		out(b)
	}

	model := unmarshal(b)

	if fjson {
		b, err := json.MarshalIndent(model, "", "  ")
		if err != nil {
			fail(err)
		}
		out(b)
	}

	if fxmp {
		out(marshal(model))
	}

	if fpath {
		l, err := model.ListPaths()
		if err != nil {
			fail(err)
		}
		for _, v := range l {
			fmt.Printf("%s = %s\n", v.Path.String(), v.Value)
		}
	}

	model.Close()

	if debug {
		xmp.DumpStats()
	}
}
