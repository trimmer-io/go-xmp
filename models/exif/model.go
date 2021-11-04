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

// XMP EXIF Mapping for Exif 2.3 metadata CIPA DC-010-2012
//
// Exif Spec
// Exif 2.3 http://www.cipa.jp/std/documents/e/DC-010-2012_E.pdf
// Exif 2.3 http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf
// Exif 2.3.1 http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
//
// see https://www.media.mit.edu/pia/Research/deepview/exif.html
// for a very good explanation of tags and ifd's

// Package exif implements the Exif 2.3.1 metadata standard as defined in CIPA DC-008-2016.
package exif

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/trimmer-io/go-xmp/models/dc"
	"github.com/trimmer-io/go-xmp/models/ps"
	"github.com/trimmer-io/go-xmp/models/tiff"
	"github.com/trimmer-io/go-xmp/models/xmp_base"
	"github.com/trimmer-io/go-xmp/xmp"
)

var (
	NsExif    *xmp.Namespace    = xmp.NewNamespace("exif", "http://ns.adobe.com/exif/1.0/", NewModel)
	NsExifEX  *xmp.Namespace    = xmp.NewNamespace("exifEX", "http://cipa.jp/exif/1.0/", NewModel)
	NsExifAux *xmp.Namespace    = xmp.NewNamespace("aux", "http://ns.adobe.com/exif/1.0/aux/", NewModel)
	nslist    xmp.NamespaceList = xmp.NamespaceList{NsExif, NsExifEX, NsExifAux}
)

func init() {
	for _, v := range nslist {
		xmp.Register(v, xmp.ImageMetadata)
	}
}

func NewModel(name string) xmp.Model {
	switch name {
	case "exif":
		return &ExifInfo{}
	case "exifEX":
		return &ExifEXInfo{}
	case "aux":
		return &ExifAuxInfo{}
	}
	return nil
}

func MakeModel(d *xmp.Document) (*ExifInfo, error) {
	m, err := d.MakeModel(NsExif)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*ExifInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *ExifInfo {
	if m := d.FindModel(NsExif); m != nil {
		return m.(*ExifInfo)
	}
	return nil
}

type ExifInfo struct {
	// Tiff Info part of Exif info
	Artist                    string                `exif:"0x013b" xmp:"tiff:Artist,omit"`
	ArtistXMP                 xmp.StringList        `exif:"-"      xmp:"dc:creator"`
	BitsPerSample             xmp.IntList           `exif:"0x0102" xmp:"tiff:BitsPerSample"` // 3 components
	Compression               tiff.CompressionType  `exif:"0x0103" xmp:"tiff:Compression"`
	Copyright                 string                `exif:"0x8298" xmp:"tiff:Copyright,omit"`
	CopyrightXMP              xmp.AltString         `exif:"-"      xmp:"dc:rights"`
	DateTime                  Date                  `exif:"0x0132" xmp:"tiff:DateTime,omit"`
	DateTimeXMP               xmp.Date              `exif:"-"      xmp:"xmp:ModifyDate"`
	ImageDescription          string                `exif:"0x010e" xmp:"tiff:ImageDescription,omit"`
	ImageDescriptionXMP       xmp.AltString         `exif:"-"      xmp:"dc:description"`
	ImageLength               int                   `exif:"0x0101" xmp:"tiff:ImageLength"`
	ImageWidth                int                   `exif:"0x0100" xmp:"tiff:ImageWidth"` // A. Tags relating to image data structure
	Make                      string                `exif:"0x010f" xmp:"tiff:Make"`
	Model                     string                `exif:"0x0110" xmp:"tiff:Model"`
	Orientation               tiff.OrientationType  `exif:"0x0112" xmp:"tiff:Orientation"`
	PhotometricInterpretation tiff.ColorModel       `exif:"0x0106" xmp:"tiff:PhotometricInterpretation"`
	PlanarConfiguration       tiff.PlanarType       `exif:"0x011c" xmp:"tiff:PlanarConfiguration"`
	PrimaryChromaticities     xmp.RationalArray     `exif:"0x013f" xmp:"tiff:PrimaryChromaticities"` // 6 components
	ReferenceBlackWhite       xmp.RationalArray     `exif:"0x0214" xmp:"tiff:ReferenceBlackWhite"`   // 6 components
	ResolutionUnit            tiff.ResolutionUnit   `exif:"0x0128" xmp:"tiff:ResolutionUnit"`        // 2 = inches, 3 = centimeters
	SamplesPerPixel           int                   `exif:"0x0115" xmp:"tiff:SamplesPerPixel"`
	Software                  string                `exif:"0x0131" xmp:"tiff:Software,omit"`
	SoftwareXMP               xmp.AgentName         `exif:"-"      xmp:"xmp:CreatorTool"`
	TransferFunction          xmp.IntList           `exif:"0x012d" xmp:"tiff:TransferFunction"` // C. Tags relating to image data characteristics
	WhitePoint                xmp.RationalArray     `exif:"0x013e" xmp:"tiff:WhitePoint"`
	XResolution               xmp.Rational          `exif:"0x011a" xmp:"tiff:XResolution"`
	YCbCrCoefficients         xmp.RationalArray     `exif:"0x0211" xmp:"tiff:YCbCrCoefficients"` // 3 components
	YCbCrPositioning          tiff.YCbCrPosition    `exif:"0x0213" xmp:"tiff:YCbCrPositioning"`
	YCbCrSubSampling          tiff.YCbCrSubSampling `exif:"0x0212" xmp:"tiff:YCbCrSubSampling"`
	YResolution               xmp.Rational          `exif:"0x011b" xmp:"tiff:YResolution"`

	// Exif info
	ExifVersion              string              `exif:"0x9000" xmp:"exif:ExifVersion"`
	FlashpixVersion          string              `exif:"0xa000" xmp:"exif:FlashpixVersion"`
	ColorSpace               ColorSpace          `exif:"0xa001" xmp:"exif:ColorSpace,empty"`
	ComponentsConfiguration  ComponentArray      `exif:"0x9101" xmp:"exif:ComponentsConfiguration"`
	CompressedBitsPerPixel   xmp.Rational        `exif:"0x9102" xmp:"exif:CompressedBitsPerPixel"`
	PixelXDimension          int                 `exif:"0xa002" xmp:"exif:PixelXDimension"`
	PixelYDimension          int                 `exif:"0xa003" xmp:"exif:PixelYDimension"`
	MakerNote                ByteArray           `exif:"0x927c" xmp:"exif:MakerNote,omit"`
	UserComment              xmp.StringArray     `exif:"0x9286" xmp:"exif:UserComment"`
	RelatedSoundFile         string              `exif:"0xa004" xmp:"exif:RelatedSoundFile"`
	DateTimeOriginal         Date                `exif:"0x9003" xmp:"exif:DateTimeOriginal,omit"`
	DateTimeOriginalXMP      xmp.Date            `exif:"-"      xmp:"photoshop:DateCreated"`
	DateTimeDigitized        Date                `exif:"0x9004" xmp:"-"`
	DateTimeDigitizedXMP     xmp.Date            `exif:"-"      xmp:"exif:DateTimeDigitized"`
	SubSecTime               string              `exif:"0x9290" xmp:"exif:SubSecTime,omit"`
	SubSecTimeOriginal       string              `exif:"0x9291" xmp:"exif:SubSecTimeOriginal,omit"`
	SubSecTimeDigitized      string              `exif:"0x9292" xmp:"exif:SubSecTimeDigitized,omit"`
	ExposureTime             xmp.Rational        `exif:"0x829a" xmp:"exif:ExposureTime"`
	FNumber                  xmp.Rational        `exif:"0x829d" xmp:"exif:FNumber"`
	ExposureProgram          ExposureProgram     `exif:"0x8822" xmp:"exif:ExposureProgram,empty"`
	SpectralSensitivity      string              `exif:"0x8824" xmp:"exif:SpectralSensitivity"`
	OECF                     *OECF               `exif:"0x8828" xmp:"exif:OECF"`
	ShutterSpeedValue        xmp.Rational        `exif:"0x9201" xmp:"exif:ShutterSpeedValue"`
	ApertureValue            xmp.Rational        `exif:"0x9202" xmp:"exif:ApertureValue"`
	BrightnessValue          xmp.Rational        `exif:"0x9203" xmp:"exif:BrightnessValue"`
	ExposureBiasValue        xmp.Rational        `exif:"0x9204" xmp:"exif:ExposureBiasValue"`
	MaxApertureValue         xmp.Rational        `exif:"0x9205" xmp:"exif:MaxApertureValue"`
	SubjectDistance          xmp.Rational        `exif:"0x9206" xmp:"exif:SubjectDistance"`
	MeteringMode             MeteringMode        `exif:"0x9207" xmp:"exif:MeteringMode,empty"`
	LightSource              LightSource         `exif:"0x9208" xmp:"exif:LightSource,empty"`
	Flash                    Flash               `exif:"0x9209" xmp:"exif:Flash"`
	FocalLength              xmp.Rational        `exif:"0x920a" xmp:"exif:FocalLength"`
	SubjectArea              xmp.IntList         `exif:"0x9214" xmp:"exif:SubjectArea"`
	FlashEnergy              xmp.Rational        `exif:"0xa20b" xmp:"exif:FlashEnergy"`
	SpatialFrequencyResponse *OECF               `exif:"0xa20c" xmp:"exif:SpatialFrequencyResponse"`
	FocalPlaneXResolution    xmp.Rational        `exif:"0xa20e" xmp:"exif:FocalPlaneXResolution"`
	FocalPlaneYResolution    xmp.Rational        `exif:"0xa20f" xmp:"exif:FocalPlaneYResolution"`
	FocalPlaneResolutionUnit tiff.ResolutionUnit `exif:"0xa210" xmp:"exif:FocalPlaneResolutionUnit,empty"`
	SubjectLocation          xmp.IntList         `exif:"0xa214" xmp:"exif:SubjectLocation"`
	ExposureIndex            xmp.Rational        `exif:"0xa215" xmp:"exif:ExposureIndex"`
	SensingMethod            SensingMode         `exif:"0xa217" xmp:"exif:SensingMethod,empty"`
	FileSource               FileSourceType      `exif:"0xa300" xmp:"exif:FileSource,empty"`
	SceneType                int                 `exif:"0xa301" xmp:"exif:SceneType,empty"`
	CFAPattern               *CFAPattern         `exif:"0xa302" xmp:"exif:CFAPattern"`
	CustomRendered           RenderMode          `exif:"0xa401" xmp:"exif:CustomRendered,empty"`
	ExposureMode             ExposureMode        `exif:"0xa402" xmp:"exif:ExposureMode,empty"`
	WhiteBalance             WhiteBalanceMode    `exif:"0xa403" xmp:"exif:WhiteBalance,empty"`
	DigitalZoomRatio         xmp.Rational        `exif:"0xa404" xmp:"exif:DigitalZoomRatio"`
	FocalLengthIn35mmFilm    int                 `exif:"0xa405" xmp:"exif:FocalLengthIn35mmFilm"`
	SceneCaptureType         SceneCaptureType    `exif:"0xa406" xmp:"exif:SceneCaptureType,empty"`
	GainControl              GainMode            `exif:"0xa407" xmp:"exif:GainControl,empty"`
	Contrast                 ContrastMode        `exif:"0xa408" xmp:"exif:Contrast,empty"`
	Saturation               SaturationMode      `exif:"0xa409" xmp:"exif:Saturation,empty"`
	Sharpness                SharpnessMode       `exif:"0xa40a" xmp:"exif:Sharpness,empty"`
	DeviceSettingDescription DeviceSettings      `exif:"0xa40b" xmp:"exif:DeviceSettingDescription"`
	SubjectDistanceRange     SubjectDistanceMode `exif:"0xa40c" xmp:"exif:SubjectDistanceRange,empty"`
	ImageUniqueID            string              `exif:"0xa420" xmp:"exif:ImageUniqueID"`
	GPSVersionID             string              `exif:"0x0000" xmp:"exif:GPSVersionID"`
	GPSLatitudeRef           string              `exif:"0x0001" xmp:"-"` // N, S
	GPSLatitude              GPSCoord            `exif:"0x0002" xmp:"-"` // 1 float or 3 rational components
	GPSLatitudeCoord         xmp.GPSCoord        `exif:"-"      xmp:"exif:GPSLatitude"`
	GPSLongitudeRef          string              `exif:"0x0003" xmp:"-"` // E, W
	GPSLongitude             GPSCoord            `exif:"0x0004" xmp:"-"` // 1 float or 3 rational components
	GPSLongitudeCoord        xmp.GPSCoord        `exif:"-"      xmp:"exif:GPSLongitude"`
	GPSAltitudeRef           string              `exif:"0x0005" xmp:"exif:GPSAltitudeRef"`
	GPSAltitude              xmp.Rational        `exif:"0x0006" xmp:"exif:GPSAltitude"`
	GPSTimeStamp             xmp.RationalArray   `exif:"0x0007" xmp:"-"` // 3 components: hour/min/sec in UTC
	GPSDateStamp             Date                `exif:"0x001d" xmp:"-"` // YYYY:MM:DD
	GPSTimeStampXMP          xmp.Date            `exif:"-"      xmp:"exif:GPSTimeStamp"`
	GPSSatellites            string              `exif:"0x0008" xmp:"exif:GPSSatellites"`
	GPSStatus                string              `exif:"0x0009" xmp:"exif:GPSStatus"`
	GPSMeasureMode           string              `exif:"0x000a" xmp:"exif:GPSMeasureMode"`
	GPSDOP                   xmp.Rational        `exif:"0x000b" xmp:"exif:GPSDOP"`
	GPSSpeedRef              string              `exif:"0x000c" xmp:"exif:GPSSpeedRef"`
	GPSSpeed                 xmp.Rational        `exif:"0x000d" xmp:"exif:GPSSpeed"`
	GPSTrackRef              string              `exif:"0x000e" xmp:"exif:GPSTrackRef"`
	GPSTrack                 xmp.Rational        `exif:"0x000f" xmp:"exif:GPSTrack"`
	GPSImgDirectionRef       string              `exif:"0x0010" xmp:"exif:GPSImgDirectionRef"`
	GPSImgDirection          xmp.Rational        `exif:"0x0011" xmp:"exif:GPSImgDirection"`
	GPSMapDatum              string              `exif:"0x0012" xmp:"exif:GPSMapDatum"`
	GPSDestLatitudeRef       string              `exif:"0x0013" xmp:"-"` // N, S
	GPSDestLatitude          GPSCoord            `exif:"0x0014" xmp:"-"` // 1 float or 3 rational components
	GPSDestLatitudeCoord     xmp.GPSCoord        `exif:"-"      xmp:"exif:GPSDestLatitude"`
	GPSDestLongitudeRef      string              `exif:"0x0015" xmp:"-"` // E, W
	GPSDestLongitude         GPSCoord            `exif:"0x0016" xmp:"-"` // 1 float or 3 rational components
	GPSDestLongitudeCoord    xmp.GPSCoord        `exif:"-"      xmp:"exif:GPSDestLongitude"`
	GPSDestBearingRef        string              `exif:"0x0017" xmp:"exif:GPSDestBearingRef"`
	GPSDestBearing           xmp.Rational        `exif:"0x0018" xmp:"exif:GPSDestBearing"`
	GPSDestDistanceRef       string              `exif:"0x0019" xmp:"exif:GPSDestDistanceRef"`
	GPSDestDistance          xmp.Rational        `exif:"0x001a" xmp:"exif:GPSDestDistance"`
	GPSProcessingMethod      string              `exif:"0x001b" xmp:"exif:GPSProcessingMethod"`
	GPSAreaInformation       string              `exif:"0x001c" xmp:"exif:GPSAreaInformation"`
	GPSDifferential          int                 `exif:"0x001e" xmp:"exif:GPSDifferential"`
	GPSHPositioningError     xmp.Rational        `exif:"0x001f" xmp:"exif:GPSHPositioningError"`
	NativeDigest             string              `exif:"-"      xmp:"exif:NativeDigest,omit"` // ingore according to spec

	// replaced by exifEX ExPhotographicSensitivity
	ISOSpeedRatings xmp.IntList `exif:"-" xmp:"exif:ISOSpeedRatings,omit"`

	// ExifEX namespace
	// ExInteroperabilityIndex     InteropMode       `exif:"0x0001" xmp:"exif:InteroperabilityIndex"`
	// ExInteroperabilityVersion   string            `exif:"0x0002" xmp:"exif:InteroperabilityVersion"`
	ExRelatedImageFileFormat    string            `exif:"0x1000" xmp:"exif:RelatedImageFileFormat"`
	ExRelatedImageWidth         int               `exif:"0x1001" xmp:"exif:RelatedImageWidth"`
	ExRelatedImageLength        int               `exif:"0x1002" xmp:"exif:RelatedImageLength"`
	ExPhotographicSensitivity   int               `exif:"0x8827" xmp:"exif:PhotographicSensitivity"`
	ExSensitivityType           SensitivityType   `exif:"0x8830" xmp:"exif:SensitivityType"`
	ExStandardOutputSensitivity int               `exif:"0x8831" xmp:"exif:StandardOutputSensitivity"`
	ExRecommendedExposureIndex  int               `exif:"0x8832" xmp:"exif:RecommendedExposureIndex"`
	ExISOSpeed                  int               `exif:"0x8833" xmp:"exif:ISOSpeed"`
	ExISOSpeedLatitudeyyy       int               `exif:"0x8834" xmp:"exif:ISOSpeedLatitudeyyy"`
	ExISOSpeedLatitudezzz       int               `exif:"0x8835" xmp:"exif:ISOSpeedLatitudezzz"`
	ExCameraOwnerName           string            `exif:"0xa430" xmp:"exif:CameraOwnerName"`
	ExBodySerialNumber          string            `exif:"0xa431" xmp:"exif:BodySerialNumber"`
	ExLensSpecification         xmp.RationalArray `exif:"0xa432" xmp:"exif:LensSpecification"` // 4 components: min/max [mm] + min/max [f-stop]
	ExLensMake                  string            `exif:"0xa433" xmp:"exif:LensMake"`
	ExLensModel                 string            `exif:"0xa434" xmp:"exif:LensModel"`
	ExLensSerialNumber          string            `exif:"0xa435" xmp:"exif:LensSerialNumber"`
	ExGamma                     xmp.Rational      `exif:"0xa500" xmp:"exif:Gamma"`

	// non-standard prefixed ExifAux namespace
	AuxApproximateFocusDistance                           xmp.Rational `xmp:"exif:ApproximateFocusDistance"` // 315/10
	AuxDistortionCorrectionAlreadyApplied                 xmp.Bool     `xmp:"exif:DistortionCorrectionAlreadyApplied"`
	AuxSerialNumber                                       string       `xmp:"exif:SerialNumber"`                                       // 2481231346
	AuxLensInfo                                           string       `xmp:"exif:LensInfo"`                                           // 10/1 22/1 0/0 0/0
	AuxLensDistortInfo                                    string       `xmp:"exif:LensDistortInfo"`                                    //
	AuxLens                                               string       `xmp:"exif:Lens"`                                               // EF-S10-22mm f/3.5-4.5 USM
	AuxLensID                                             string       `xmp:"exif:LensID"`                                             // 235
	AuxImageNumber                                        int          `xmp:"exif:ImageNumber"`                                        // 0
	AuxIsMergedHDR                                        xmp.Bool     `xmp:"exif:IsMergedHDR"`                                        //
	AuxIsMergedPanorama                                   xmp.Bool     `xmp:"exif:IsMergedPanorama"`                                   //
	AuxLateralChromaticAberrationCorrectionAlreadyApplied xmp.Bool     `xmp:"exif:LateralChromaticAberrationCorrectionAlreadyApplied"` //
	AuxVignetteCorrectionAlreadyApplied                   xmp.Bool     `xmp:"exif:VignetteCorrectionAlreadyApplied"`                   //
	AuxFlashCompensation                                  xmp.Rational `xmp:"exif:FlashCompensation"`                                  // 0/1
	AuxFirmware                                           string       `xmp:"exif:Firmware"`                                           // 1.2.5
	AuxOwnerName                                          string       `xmp:"exif:OwnerName"`                                          // unknown
	// AuxLensSerialNumber                                   string        `xmp:"exif:LensSerialNumber"`
}

func (m *ExifInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsExif}
}

func (m *ExifInfo) Can(nsName string) bool {
	return nsName == NsExif.GetName()
}

func (x *ExifInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *ExifInfo) SyncFromXMP(d *xmp.Document) error {

	if !x.DateTimeDigitizedXMP.IsZero() {
		x.DateTimeDigitized = Date(x.DateTimeDigitizedXMP.Value())
		x.SubSecTimeDigitized = strconv.Itoa(x.DateTimeDigitizedXMP.Value().Nanosecond())
	}

	if len(x.ISOSpeedRatings) > 0 {
		x.ExPhotographicSensitivity = x.ISOSpeedRatings[0]
	}

	// TODO
	// convert GPS from XMP to exif values

	if m := dc.FindModel(d); m != nil {
		x.ArtistXMP = m.Creator
		x.Artist = strings.Join(m.Creator, ",")
		x.ImageDescriptionXMP = m.Description
		x.ImageDescription = m.Description.Default()
		x.CopyrightXMP = m.Rights
		x.Copyright = m.Rights.Default()
	}
	if base := xmpbase.FindModel(d); base != nil {
		if !base.ModifyDate.IsZero() {
			x.DateTimeXMP = base.ModifyDate
			x.DateTime = Date(base.ModifyDate.Value())
			x.SubSecTime = strconv.Itoa(base.ModifyDate.Value().Nanosecond())
		}
		if base.CreatorTool.IsZero() {
			x.SoftwareXMP = base.CreatorTool
			x.Software = base.CreatorTool.String()
		}
	}
	if tff := tiff.FindModel(d); tff != nil {
		x.BitsPerSample = tff.BitsPerSample
		x.Compression = tff.Compression
		x.ImageLength = tff.ImageLength
		x.ImageWidth = tff.ImageWidth
		x.Make = tff.Make
		x.Model = tff.Model
		x.Orientation = tff.Orientation
		x.PhotometricInterpretation = tff.PhotometricInterpretation
		x.PlanarConfiguration = tff.PlanarConfiguration
		x.PrimaryChromaticities = tff.PrimaryChromaticities
		x.ReferenceBlackWhite = tff.ReferenceBlackWhite
		x.ResolutionUnit = tff.ResolutionUnit
		x.SamplesPerPixel = tff.SamplesPerPixel
		x.TransferFunction = tff.TransferFunction
		x.WhitePoint = tff.WhitePoint
		x.XResolution = tff.XResolution
		x.YCbCrCoefficients = tff.YCbCrCoefficients
		x.YCbCrPositioning = tff.YCbCrPositioning
		x.YCbCrSubSampling = tff.YCbCrSubSampling
		x.YResolution = tff.YResolution
	}
	if m := ps.FindModel(d); m != nil {
		if !m.DateCreated.IsZero() {
			x.DateTimeOriginalXMP = m.DateCreated
			x.DateTimeOriginal = Date(m.DateCreated.Value())
			x.SubSecTimeOriginal = strconv.Itoa(m.DateCreated.Value().Nanosecond())
		}
	}
	return nil
}

func (x *ExifInfo) SyncToXMP(d *xmp.Document) error {
	// convert dates, text and GPS properties, ignore errors
	if !x.DateTimeOriginal.IsZero() {
		x.DateTimeOriginalXMP, _ = convertDateToXMP(x.DateTimeOriginal, x.SubSecTimeOriginal)
	}
	if !x.DateTimeDigitized.IsZero() {
		x.DateTimeDigitizedXMP, _ = convertDateToXMP(x.DateTimeDigitized, x.SubSecTimeDigitized)
	}
	if !x.DateTime.IsZero() {
		x.DateTimeXMP, _ = convertDateToXMP(x.DateTime, x.SubSecTime)
	}
	if x.Artist != "" {
		x.ArtistXMP = xmp.StringList(strings.Split(x.Artist, ";"))
	}
	if x.ImageDescription != "" && len(x.ImageDescriptionXMP) == 0 {
		x.ImageDescriptionXMP = xmp.NewAltString(x.ImageDescription)
	}
	if x.Copyright != "" && len(x.CopyrightXMP) == 0 {
		x.CopyrightXMP = xmp.NewAltString(x.Copyright)
	}
	if x.Software != "" && x.SoftwareXMP.IsZero() {
		x.SoftwareXMP = xmp.AgentName(x.Software)
	}

	if len(x.ISOSpeedRatings) > 0 {
		x.ExPhotographicSensitivity = x.ISOSpeedRatings[0]
	}

	// convert GPS coordinates
	if v, err := convertGPStoXMP(x.GPSLatitude, x.GPSLatitudeRef); err == nil && v != "" {
		x.GPSLatitudeCoord = v
	}
	if v, err := convertGPStoXMP(x.GPSLongitude, x.GPSLongitudeRef); err == nil && v != "" {
		x.GPSLongitudeCoord = v
	}
	if v, err := convertGPStoXMP(x.GPSDestLatitude, x.GPSDestLatitudeRef); err == nil && v != "" {
		x.GPSDestLatitudeCoord = v
	}
	if v, err := convertGPStoXMP(x.GPSDestLongitude, x.GPSDestLongitudeRef); err == nil && v != "" {
		x.GPSDestLongitudeCoord = v
	}
	if !x.GPSDateStamp.IsZero() {
		x.GPSTimeStampXMP = convertGPSTimestamp(x.GPSDateStamp, x.GPSTimeStamp)
	}
	return nil
}

func (x *ExifInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *ExifInfo) GetTag(tag string) (string, error) {
	tag = strings.ToLower(tag)
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("exif: %v", err)
	} else {
		return v, nil
	}
}

func (x *ExifInfo) SetTag(tag, value string) error {
	tag = strings.ToLower(tag)
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("exif: %v", err)
	}
	return nil
}

func (x *ExifInfo) GetLocaleTag(lang string, tag string) (string, error) {
	tag = strings.ToLower(tag)
	if val, err := xmp.GetLocaleField(x, lang, tag); err != nil {
		return "", fmt.Errorf("exif: %v", err)
	} else {
		return val, nil
	}
}

func (x *ExifInfo) SetLocaleTag(lang string, tag, value string) error {
	tag = strings.ToLower(tag)
	if err := xmp.SetLocaleField(x, lang, tag, value); err != nil {
		return fmt.Errorf("exif: %v", err)
	}
	return nil
}

type ExifEXInfo struct {
	InteroperabilityIndex     InteropMode       `exif:"0x0001" xmp:"exifEX:InteroperabilityIndex"`
	InteroperabilityVersion   string            `exif:"0x0002" xmp:"exifEX:InteroperabilityVersion"`
	RelatedImageFileFormat    string            `exif:"0x1000" xmp:"exifEX:RelatedImageFileFormat"`
	RelatedImageWidth         int               `exif:"0x1001" xmp:"exifEX:RelatedImageWidth"`
	RelatedImageLength        int               `exif:"0x1002" xmp:"exifEX:RelatedImageLength"`
	PhotographicSensitivity   int               `exif:"0x8827" xmp:"exifEX:PhotographicSensitivity"`
	SensitivityType           SensitivityType   `exif:"0x8830" xmp:"exifEX:SensitivityType"`
	StandardOutputSensitivity int               `exif:"0x8831" xmp:"exifEX:StandardOutputSensitivity"`
	RecommendedExposureIndex  int               `exif:"0x8832" xmp:"exifEX:RecommendedExposureIndex"`
	ISOSpeed                  int               `exif:"0x8833" xmp:"exifEX:ISOSpeed"`
	ISOSpeedLatitudeyyy       int               `exif:"0x8834" xmp:"exifEX:ISOSpeedLatitudeyyy"`
	ISOSpeedLatitudezzz       int               `exif:"0x8835" xmp:"exifEX:ISOSpeedLatitudezzz"`
	CameraOwnerName           string            `exif:"0xa430" xmp:"exifEX:CameraOwnerName"`
	BodySerialNumber          string            `exif:"0xa431" xmp:"exifEX:BodySerialNumber"`
	LensSpecification         xmp.RationalArray `exif:"0xa432" xmp:"exifEX:LensSpecification"` // 4 components: min/max [mm] + min/max [f-stop]
	LensMake                  string            `exif:"0xa433" xmp:"exifEX:LensMake"`
	LensModel                 string            `exif:"0xa434" xmp:"exifEX:LensModel"`
	LensSerialNumber          string            `exif:"0xa435" xmp:"exifEX:LensSerialNumber"`
	Gamma                     xmp.Rational      `exif:"0xa500" xmp:"exifEX:Gamma"`
}

func (m *ExifEXInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsExifEX}
}

func (m *ExifEXInfo) Can(nsName string) bool {
	return nsName == NsExifEX.GetName()
}

func (x *ExifEXInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *ExifEXInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x ExifEXInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *ExifEXInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *ExifEXInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsExifEX.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *ExifEXInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsExifEX.GetName(), err)
	}
	return nil
}

// XMP-EXIF Extra Mappings (XMP only)
// The schema namespace URI is http://ns.adobe.com/exif/1.0/aux/
// The preferred schema namespace prefix is aux
//
// Adobe-defined auxiliary EXIF tags.  This namespace existed in the XMP
// specification until it was dropped in 2012, presumably due to the
// introduction of the EXIF 2.3 for XMP specification and the exifEX namespace
// at this time.
type ExifAuxInfo struct {
	ApproximateFocusDistance                           xmp.Rational `xmp:"aux:ApproximateFocusDistance"`                           // 315/10
	DistortionCorrectionAlreadyApplied                 xmp.Bool     `xmp:"aux:DistortionCorrectionAlreadyApplied"`                 //
	SerialNumber                                       string       `xmp:"aux:SerialNumber"`                                       // 2481231346
	LensInfo                                           string       `xmp:"aux:LensInfo"`                                           // 10/1 22/1 0/0 0/0
	LensDistortInfo                                    string       `xmp:"aux:LensDistortInfo"`                                    //
	Lens                                               string       `xmp:"aux:Lens"`                                               // EF-S10-22mm f/3.5-4.5 USM
	LensID                                             string       `xmp:"aux:LensID"`                                             // 235
	LensSerialNumber                                   string       `xmp:"aux:LensSerialNumber"`                                   //
	ImageNumber                                        int          `xmp:"aux:ImageNumber"`                                        // 0
	IsMergedHDR                                        xmp.Bool     `xmp:"aux:IsMergedHDR"`                                        //
	IsMergedPanorama                                   xmp.Bool     `xmp:"aux:IsMergedPanorama"`                                   //
	LateralChromaticAberrationCorrectionAlreadyApplied xmp.Bool     `xmp:"aux:LateralChromaticAberrationCorrectionAlreadyApplied"` //
	VignetteCorrectionAlreadyApplied                   xmp.Bool     `xmp:"aux:VignetteCorrectionAlreadyApplied"`                   //
	FlashCompensation                                  xmp.Rational `xmp:"aux:FlashCompensation"`                                  // 0/1
	Firmware                                           string       `xmp:"aux:Firmware"`                                           // 1.2.5
	OwnerName                                          string       `xmp:"aux:OwnerName"`                                          // unknown
}

func (m *ExifAuxInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsExifAux}
}

func (m *ExifAuxInfo) Can(nsName string) bool {
	return nsName == NsExifAux.GetName()
}

func (x *ExifAuxInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *ExifAuxInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x ExifAuxInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *ExifAuxInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *ExifAuxInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsExifAux.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *ExifAuxInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsExifAux.GetName(), err)
	}
	return nil
}
