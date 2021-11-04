// Trimmer Media SDK
//
// Copyright (c) 2013-2016 KIDTSUNAMI
// Author: alex@kidtsunami.com
//

package qt

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/trimmer-io/go-xmp/xmp"
)

type LocationRole int

const (
	LocationRoleShooting  LocationRole = 0
	LocationRoleReal      LocationRole = 1
	LocationRoleFictional LocationRole = 2
)

type Bool int // 0 no, 1 yes

func (b Bool) Value() bool {
	return b == 1
}

// defined to overwrite UnmarshalText, otherwise similar to xmp.AltString
//
type MultilangArray xmp.AltString

func (x *MultilangArray) UnmarshalText(data []byte) error {
	// TODO: need samples
	return nil
}

func (x MultilangArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, xmp.AltString(x).Typ(), xmp.AltString(x))
}

func (x *MultilangArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, xmp.AltString(*x).Typ(), (*xmp.AltString)(x))
}

type Location struct {
	Body      string       `xmp:"qt:Body,attr"`
	Date      xmp.Date     `xmp:"qt:Date,attr"`
	Longitude xmp.GPSCoord `xmp:"qt:Longitude,attr"`
	Latitude  xmp.GPSCoord `xmp:"qt:Latitude,attr"`
	Altitude  float64      `xmp:"qt:Altitude,attr"`
	Name      string       `xmp:"qt:Name,attr"`
	Note      string       `xmp:"qt:Note,attr"`
	Role      LocationRole `xmp:"qt:Role,attr"`
}

func (x *Location) UnmarshalText(data []byte) error {
	// todo: need sample data
	return nil
}

type Point struct {
	X int
	Y int
}

func (x Point) IsZero() bool {
	return x.X == 0 && x.Y == 0
}

func (x Point) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(strconv.FormatInt(int64(x.X), 10))
	buf.WriteByte(',')
	buf.WriteString(strconv.FormatInt(int64(x.Y), 10))
	return buf.String()
}

func (x *Point) UnmarshalText(data []byte) error {
	v := string(data)
	f := "%d,%d"
	switch true {
	case strings.Contains(v, "/"):
		f = "%d/%d"
	case strings.Contains(v, ":"):
		f = "%d:%d"
	case strings.Contains(v, " "):
		f = "%d %d"
	case strings.Contains(v, ";"):
		f = "%d;%d"
	}
	p := Point{}
	if _, err := fmt.Sscanf(v, f, &p.X, &p.Y); err != nil {
		return fmt.Errorf("qt: invalid point value '%s': %v", v, err)
	}
	*x = p
	return nil
}

func (x Point) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(x.String()), nil
}

type ContentRating struct {
	Standard string `xmp:"qt:Standard,attr"`
	Rating   string `xmp:"qt:Rating,attr"`
	Score    string `xmp:"qt:Score,attr"`
	Reasons  string `xmp:"qt:Reasons,attr"`
}

func (x ContentRating) IsZero() bool {
	return x.Standard == "" && x.Rating == ""
}

func (x ContentRating) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('(')
	buf.WriteString(x.Standard)
	buf.WriteByte('|')
	buf.WriteString(x.Rating)
	buf.WriteByte('|')
	buf.WriteString(x.Score)
	buf.WriteByte('|')
	buf.WriteString(x.Reasons)
	return buf.String()
}

func (x *ContentRating) UnmarshalText(data []byte) error {
	cr := ContentRating{}
	v := string(data)
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "(")
	v = strings.TrimSuffix(v, ")")
	for i, s := range strings.Split(v, "|") {
		switch i {
		case 1:
			cr.Standard = s
		case 2:
			cr.Rating = s
		case 3:
			cr.Score = s
		case 4:
			cr.Reasons = s
		}
	}
	*x = cr
	return nil
}

func (x ContentRating) MarshalText() ([]byte, error) {
	if x.IsZero() {
		return nil, nil
	}
	return []byte(x.String()), nil
}
