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
	"bufio"
	"bytes"
	"io"
)

func Read(r io.Reader) (*Document, error) {
	x := &Document{}
	d := NewDecoder(r)
	if err := d.Decode(x); err != nil {
		return nil, err
	}
	return x, nil
}

func Scan(r io.Reader) (*Document, error) {
	pp, err := ScanPackets(r)
	if err != nil {
		return nil, err
	}
	x := &Document{}
	if err := Unmarshal(pp[0], x); err != nil {
		return nil, err
	}
	return x, nil
}

func ScanPackets(r io.Reader) ([][]byte, error) {
	packets := make([][]byte, 0)
	s := bufio.NewScanner(r)
	s.Split(splitPacket)
	for s.Scan() {
		b := s.Bytes()
		if isXmpPacket(b) {
			x := make([]byte, len(b))
			copy(x, b)
			packets = append(packets, x)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(packets) == 0 {
		return nil, io.EOF
	}
	return packets, nil
}

var packet_start = []byte("<?xpacket begin")
var packet_end = []byte("<?xpacket end")       // plus suffix `"w"?>`
var magic = []byte("W5M0MpCehiHzreSzNTczkc9d") // len 24

func isXmpPacket(b []byte) bool {
	return bytes.Index(b[:51], magic) > -1
}

func splitPacket(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := bytes.Index(data, packet_start)
	if start == -1 {
		ofs := len(data) - len(packet_start)
		if ofs > 0 {
			return ofs, nil, nil
		}
		return len(data), nil, nil
	}
	end := bytes.Index(data[start:], packet_end)
	last := start + end + len(packet_end) + 6
	if end == -1 || last > len(data) {
		if atEOF {
			return len(data), nil, nil
		}
		return 0, nil, nil
	}
	return last, data[start:last], nil
}
