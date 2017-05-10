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

package qt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/echa/go-xmp/models/tiff"
	"github.com/echa/go-xmp/xmp"
)

// QuickTime Metadata Keys written by the "mdta" handler using Reverse DNS names.
// - "com.apple.quicktime." for Apple Quicktime default keys
// - "com.apple.proapps." for Apple Professional apps keys
// - "com.arri.camera" for Arri Camera metadata
// - "com.panasonic.professionalplugin" for Panasonic cameras
// - "com.panasonic.p2.metadata.bmp": "BMf*",
// - "com.panasonic.p2.metadata.xml": <xml>
//

// General Quicktime Metadata keys (from MacOS CoreMedia Framework headers and libs)
//
// VideoOrientation: like TIFF/EXIF Orientation property
// AffineTransform: used in place of a track display matrix, 3x3 float matrix in row major order
// Direction: degree offset from magnetic north
//
type QtMetadata struct {

	// official metadata keys as defined in MacOS 10.11 SDK
	Album                  string         `qt:"com.apple.quicktime.album"                                 qt:"©alb"      xmp:"qt:Album"`
	Arranger               string         `qt:"com.apple.quicktime.arranger"                              qt:"©arg"      xmp:"qt:Arranger"`
	Artist                 string         `qt:"com.apple.quicktime.artist"                                qt:"albr"      xmp:"qt:Artist"`
	Artwork                string         `qt:"com.apple.quicktime.artwork"                               qt:"covr"      xmp:"qt:Artwork"`
	Author                 string         `qt:"com.apple.quicktime.author"                                qt:"©aut"      xmp:"qt:Author"`
	CameraFrameReadoutTime string         `qt:"com.apple.quicktime.camera.framereadouttimeinmicroseconds" qt:"-"         xmp:"qt:CameraFrameReadoutTime"`
	CameraIdentifier       string         `qt:"com.apple.quicktime.camera.identifier"                     qt:"cmid"      xmp:"qt:CameraIdentifier"`
	CollectionUser         string         `qt:"com.apple.quicktime.collection.user"                       qt:"coll"      xmp:"qt:CollectionUser"`
	Comment                string         `qt:"com.apple.quicktime.comment"                               qt:"©cmt"      xmp:"qt:Comment"`
	Composer               string         `qt:"com.apple.quicktime.composer"                              qt:"©wrt"      xmp:"qt:Composer"`
	ContentIdentifier      string         `qt:"com.apple.quicktime.content.identifier"                    qt:"cnID"      xmp:"qt:ContentIdentifier"`
	Copyright              string         `qt:"com.apple.quicktime.copyright"                             qt:"cprt"      xmp:"qt:Copyright"`
	CreationDate           xmp.Date       `qt:"com.apple.quicktime.creationdate"                          qt:"date"      xmp:"qt:CreationDate"`
	Credits                string         `qt:"com.apple.quicktime.credits"                               qt:"©src"      xmp:"qt:Credits"`
	Description            string         `qt:"com.apple.quicktime.description"                           qt:"©des"      xmp:"qt:Description"`
	DirectionFacing        float64        `qt:"com.apple.quicktime.direction.facing"                      qt:"-"         xmp:"qt:DirectionFacing"`
	DirectionMotion        float64        `qt:"com.apple.quicktime.direction.motion"                      qt:"-"         xmp:"qt:DirectionMotion"`
	Director               string         `qt:"com.apple.quicktime.director"                              qt:"©dir"      xmp:"qt:Director"`
	DisplayName            string         `qt:"com.apple.quicktime.displayname"                           qt:"name"      xmp:"qt:DisplayName"`
	EncodedBy              string         `qt:"com.apple.quicktime.encodedby"                             qt:"©enc"      xmp:"qt:EncodedBy"`
	Genre                  string         `qt:"com.apple.quicktime.genre"                                 qt:"©gen"      xmp:"qt:Genre"`
	Information            string         `qt:"com.apple.quicktime.information"                           qt:"©inf"      xmp:"qt:Information"`
	Keywords               xmp.StringList `qt:"com.apple.quicktime.keywords"                              qt:"©key"      xmp:"qt:Keywords"`
	LocationBody           string         `qt:"com.apple.quicktime.location.body"                         qt:"-"         xmp:"-"`
	LocationDate           xmp.Date       `qt:"com.apple.quicktime.location.date"                         qt:"-"         xmp:"-"`
	LocationISO6709        string         `qt:"com.apple.quicktime.location.iso6709"                      qt:"©xyz"      xmp:"-"`
	LocationName           string         `qt:"com.apple.quicktime.location.name"                         qt:"-"         xmp:"-"`
	LocationNote           string         `qt:"com.apple.quicktime.location.note"                         qt:"-"         xmp:"-"`
	LocationRole           LocationRole   `qt:"com.apple.quicktime.location.role"                         qt:"-"         xmp:"-"`
	Make                   string         `qt:"com.apple.quicktime.make"                                  qt:"©mak"      xmp:"qt:Make"`
	Model                  string         `qt:"com.apple.quicktime.model"                                 qt:"©mod"      xmp:"qt:Model"`
	OriginalArtist         string         `qt:"com.apple.quicktime.originalartist"                        qt:"©ope"      xmp:"qt:OriginalArtist"`
	Performer              string         `qt:"com.apple.quicktime.performer"                             qt:"©prf"      xmp:"qt:Performer"`
	PhonogramRights        string         `qt:"com.apple.quicktime.phonogramrights"                       qt:"©phg"      xmp:"qt:PhonogramRights"`
	Producer               string         `qt:"com.apple.quicktime.producer"                              qt:"©prd"      xmp:"qt:Producer"`
	Publisher              string         `qt:"com.apple.quicktime.publisher"                             qt:"©pub"      xmp:"qt:Publisher"`
	USRating               float64        `qt:"com.apple.quicktime.rating.user"                           qt:"rtng"      xmp:"qt:USRating"`
	Software               string         `qt:"com.apple.quicktime.software"                              qt:"©swr"      xmp:"qt:Software"`
	Title                  string         `qt:"com.apple.quicktime.title"                                 qt:"©nam"      xmp:"qt:Title"`
	Year                   int            `qt:"com.apple.quicktime.year"                                  qt:"yrrc"      xmp:"qt:Year"`

	// more tags not captured above
	Version                  string               `qt:"com.apple.quicktime.version"                                       qt:"VERS" xmp:"qt:Version"`
	PreferredAffineTransform string               `qt:"com.apple.quicktime.preferred-affine-transform"                    qt:"-"    xmp:"qt:PreferredAffineTransform"`
	VideoOrientation         tiff.OrientationType `qt:"com.apple.quicktime.video-orientation"                             qt:"-"    xmp:"qt:VideoOrientation"`
	WindowLocation           Point                `qt:"com.apple.quicktime.windowlocation"                                qt:"WLOC" xmp:"qt:WindowLocation"`
	CoreMotion               string               `qt:"com.apple.quicktime.core-motion"                                   qt:"-"    xmp:"qt:CoreMotion"`
	CameraDebugInfo          string               `qt:"com.apple.quicktime.camera-debug-info"                             qt:"-"    xmp:"qt:CameraDebugInfo"`
	IsMontage                string               `qt:"com.apple.quicktime.is-montage"                                    qt:"-"    xmp:"qt:IsMontage"`
	PixelDensity             string               `qt:"com.apple.quicktime.pixeldensity"                                  qt:"-"    xmp:"qt:PixelDensity"`
	DetectedFace             string               `qt:"com.apple.quicktime.detected-face"                                 qt:"-"    xmp:"qt:DetectedFace"`
	HasEAN13                 Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.EAN13"          qt:"-"    xmp:"qt:HasEAN13"`
	HasEAN8                  Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.EAN8"           qt:"-"    xmp:"qt:HasEAN8"`
	HasUPCE                  Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.UPCE"           qt:"-"    xmp:"qt:HasUPCE"`
	HasCode39                Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.Code39"         qt:"-"    xmp:"qt:HasCode39"`
	HasCode39Checksum        Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.Code39Checksum" qt:"-"    xmp:"qt:HasCode39Checksum"`
	HasCode93                Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.Code93"         qt:"-"    xmp:"qt:HasCode93"`
	HasCode128               Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.Code128"        qt:"-"    xmp:"qt:HasCode128"`
	HasI2of5                 Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.I2of5"          qt:"-"    xmp:"qt:HasI2of5"`
	HasITF14                 Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.ITF14"          qt:"-"    xmp:"qt:HasITF14"`
	HasDataMatrix            Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.DataMatrix"     qt:"-"    xmp:"qt:HasDataMatrix"`
	HasQR                    Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.QR"             qt:"-"    xmp:"qt:HasQR"`
	HasAztec                 Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.Aztec"          qt:"-"    xmp:"qt:HasAztec"`
	HasPDF417                Bool                 `qt:"com.apple.quicktime.detected-machine-readable-code.PDF417"         qt:"-"    xmp:"qt:HasPDF417"`
	Extension                xmp.TagList          `qt:",any"                                                              qt:",any" xmp:"qt:extension"`

	// composite structs
	Location *Location `qt:"-"  qt:"-"  xmp:"qt:Location"`
}

func (m *QtMetadata) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsQuicktime}
}

func (m *QtMetadata) Can(nsName string) bool {
	return nsName == NsQuicktime.GetName()
}

func (x *QtMetadata) SyncFromXMP(d *xmp.Document) error {
	if x.Location != nil {
		x.LocationBody = x.Location.Body
		x.LocationDate = x.Location.Date
		x.LocationISO6709 = strings.Join([]string{
			string(x.Location.Latitude),
			string(x.Location.Longitude),
			strconv.FormatFloat(x.Location.Altitude, 'f', -1, 64),
		}, "+")
		x.LocationName = x.Location.Name
		x.LocationNote = x.Location.Note
		x.LocationRole = x.Location.Role
	}
	return nil
}

func (x *QtMetadata) SyncToXMP(d *xmp.Document) error {
	if x.Location == nil && x.LocationISO6709 != "" && x.LocationName != "" {
		x.Location = &Location{
			Body: x.LocationBody,
			Date: x.LocationDate,
			Name: x.LocationName,
			Note: x.LocationNote,
			Role: x.LocationRole,
		}
		// "+43.5968+005.4797+362.904/",
		for i, v := range strings.Split(x.LocationISO6709, "+") {
			switch i {
			case 0:
				x.Location.Latitude = xmp.GPSCoord(v)
			case 1:
				x.Location.Longitude = xmp.GPSCoord(v)
			case 2:
				v = strings.TrimSuffix(v, "/")
				x.Location.Altitude, _ = strconv.ParseFloat(v, 64)
			}
		}
	}
	return nil
}

func (x *QtMetadata) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *QtMetadata) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *QtMetadata) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	}
	return nil
}
