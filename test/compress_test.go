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

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang/snappy"
	"github.com/montanaflynn/stats"

	_ "github.com/echa/go-xmp/models"
	"github.com/echa/go-xmp/xmp"
)

// Compression tests
//

type TestStats struct {
	OrigSizes       []int
	XmpSizes        []int
	JsonSizes       []int
	ReadTimes       []time.Duration
	JsonTimes       []time.Duration
	XmpTimes        []time.Duration
	XmpGzipSizes    []int
	XmpGzipTimes    []time.Duration
	JsonGzipSizes   []int
	JsonGzipTimes   []time.Duration
	XmpSnappySizes  []int
	XmpSnappyTimes  []time.Duration
	JsonSnappySizes []int
	JsonSnappyTimes []time.Duration
}

func newStats() *TestStats {
	return &TestStats{
		OrigSizes:       make([]int, len(testfiles)),
		XmpSizes:        make([]int, len(testfiles)),
		JsonSizes:       make([]int, len(testfiles)),
		ReadTimes:       make([]time.Duration, len(testfiles)),
		JsonTimes:       make([]time.Duration, len(testfiles)),
		XmpTimes:        make([]time.Duration, len(testfiles)),
		XmpGzipSizes:    make([]int, len(testfiles)),
		XmpGzipTimes:    make([]time.Duration, len(testfiles)),
		JsonGzipSizes:   make([]int, len(testfiles)),
		JsonGzipTimes:   make([]time.Duration, len(testfiles)),
		XmpSnappySizes:  make([]int, len(testfiles)),
		XmpSnappyTimes:  make([]time.Duration, len(testfiles)),
		JsonSnappySizes: make([]int, len(testfiles)),
		JsonSnappyTimes: make([]time.Duration, len(testfiles)),
	}
}

func TestCompression(T *testing.T) {

	st := newStats()

	// run tests
	for i, v := range testfiles {
		f, err := os.Open(v)
		if err != nil {
			T.Errorf("Cannot open sample '%s': %v", v, err)
		}
		if s, err := f.Stat(); err == nil {
			st.OrigSizes[i] = int(s.Size())
		}

		start := time.Now()
		d := xmp.NewDecoder(f)
		doc := &xmp.Document{}
		if err := d.Decode(doc); err != nil {
			T.Errorf("%s: %v", v, err)
		}
		f.Close()
		st.ReadTimes[i] = time.Since(start)

		// JSON target
		start = time.Now()
		buf, _ = json.Marshal(doc)
		st.JsonTimes[i] = time.Since(start)
		st.JsonSizes[i] = len(buf)
		st.JsonGzipSizes[i], st.JsonGzipTimes[i] = gz(T, buf)
		st.JsonSnappySizes[i], st.JsonSnappyTimes[i] = snap(T, buf)

		// XMP target
		start = time.Now()
		buf, _ := xmp.Marshal(doc)
		st.XmpTimes[i] = time.Since(start)
		st.XmpSizes[i] = len(buf)
		st.XmpGzipSizes[i], st.XmpGzipTimes[i] = gz(T, buf)
		st.XmpSnappySizes[i], st.XmpSnappyTimes[i] = snap(T, buf)

		doc.Close()
	}

	// summarize result
	T.Logf("Compression Results %d        mean               min                  max\n", len(testfiles))
	T.Logf("-----------------------------------------------------------------------------\n")
	printStats(T, "Original sizes", st.OrigSizes, st.OrigSizes)
	printStats(T, "XMP sizes", st.XmpSizes, st.OrigSizes)
	printStats(T, "XMP Gzip sizes", st.XmpGzipSizes, st.OrigSizes)
	printStats(T, "XMP Snappy sizes", st.XmpSnappySizes, st.OrigSizes)
	printStats(T, "JSON sizes", st.JsonSizes, st.OrigSizes)
	printStats(T, "JSON Gzip sizes", st.JsonGzipSizes, st.OrigSizes)
	printStats(T, "JSON Snappy sizes", st.JsonSnappySizes, st.OrigSizes)
	T.Logf("-----------------------------------------------------------------------------\n")
	printTimes(T, "XML->XMP times", st.ReadTimes)
	printTimes(T, "XMP->JSON times", st.JsonTimes)
	printTimes(T, "XMP->XML times", st.XmpTimes)
	printTimes(T, "XMP Gzip times", st.XmpGzipTimes)
	printTimes(T, "XMP Snappy times", st.XmpSnappyTimes)
	printTimes(T, "JSON Gzip times", st.JsonGzipTimes)
	printTimes(T, "JSON Snappy times", st.JsonSnappyTimes)
}

func printStats(T *testing.T, s string, v, o []int) {
	pct := make([]float64, len(v))
	for i, l := 0, len(v); i < l; i++ {
		pct[i] = float64(v[i]) * 100 / float64(o[i])
	}
	val := stats.LoadRawData(v)
	pval := stats.LoadRawData(pct)
	mean, _ := val.Mean()
	max, _ := val.Max()
	min, _ := val.Min()
	pmean, _ := pval.Mean()
	pmax, _ := pval.Max()
	pmin, _ := pval.Min()
	T.Logf("%20s %10d (%5.1f) %10d (%5.1f) %10d (%5.1f)", s, int(mean), pmean, int(min), pmin, int(max), pmax)
}

func printTimes(T *testing.T, s string, v []time.Duration) {
	val := stats.LoadRawData(v)
	mean, _ := val.Mean()
	max, _ := val.Max()
	min, _ := val.Min()
	T.Logf("%20s %18v %18v %18v", s, time.Duration(mean), time.Duration(min), time.Duration(max))
}

func gz(T *testing.T, b []byte) (int, time.Duration) {
	start := time.Now()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(b)
	if err != nil {
		T.Error(err)
	}
	if err := zw.Close(); err != nil {
		T.Error(err)
	}
	return buf.Len(), time.Since(start)
}

func snap(T *testing.T, b []byte) (int, time.Duration) {
	start := time.Now()
	var buf bytes.Buffer
	sw := snappy.NewBufferedWriter(&buf)
	_, err := sw.Write(b)
	if err != nil {
		T.Error(err)
	}
	if err := sw.Close(); err != nil {
		T.Error(err)
	}
	return buf.Len(), time.Since(start)
}
