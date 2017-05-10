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
	"github.com/echa/go-xmp/xmp"
)

type QtProApps struct {
	ClipID                       string   `qt:"com.apple.proapps.clipID"                                xmp:"qt:ClipID"`
	OriginalFormat               string   `qt:"com.apple.proapps.originalFormat"                        xmp:"qt:OriginalFormat"` // AVC-Intra 4:2:2, Apple XDCAM HD422 1080i50 (50 Mb/s CBR)
	CustomColor                  string   `qt:"com.apple.proapps.customcolor"                           xmp:"qt:CustomColor"`    // com.sony.s-log3-cine
	CameraName                   string   `qt:"com.apple.proapps.cameraName"                            xmp:"qt:CameraName"`
	MIOCameraName                string   `qt:"com.apple.proapps.mio.cameraName"                        xmp:"qt:MIOCameraName"`
	Scene                        string   `qt:"com.apple.proapps.scene"                                 xmp:"qt:Scene"`
	Shot                         string   `qt:"com.apple.proapps.shot"                                  xmp:"qt:Shot"`
	Angle                        string   `qt:"com.apple.proapps.angle"                                 xmp:"qt:Angle"`
	LogNote                      string   `qt:"com.apple.proapps.logNote"                               xmp:"qt:LogNote"`
	LastUpdateDate               xmp.Date `qt:"com.apple.proapps.lastupdatedate"                        xmp:"qt:Lastupdatedate"`        // 2016-07-21 08:23:35 +0000
	IngestDateDescription        string   `qt:"com.apple.proapps.ingestDate.description"                xmp:"qt:IngestDateDescription"` // 2016-07-25 12:55:32 +0200
	StartTimeCodeFrameCount      int64    `qt:"com.apple.proapps.avchd.startTimeCodeFrameCount"         xmp:"qt:StartTimeCodeFrameCount"`
	DropFrame                    Bool     `qt:"com.apple.proapps.avchd.dropFrame"                       xmp:"qt:DropFrame"`
	VideoFrameDuration           int64    `qt:"com.apple.proapps.avchd.videoFrameDuration"              xmp:"qt:VideoFrameDuration"`
	Category                     string   `qt:"com.apple.proapps.share.category"                        xmp:"qt:Category"`
	ShareID                      string   `qt:"com.apple.proapps.share.id"                              xmp:"qt:ShareID"`
	EpisodeID                    string   `qt:"com.apple.proapps.share.episodeID"                       xmp:"qt:EpisodeID"`
	EpisodeNumber                string   `qt:"com.apple.proapps.share.episodeNumber"                   xmp:"qt:EpisodeNumber"`
	MediaKind                    string   `qt:"com.apple.proapps.share.mediaKind"                       xmp:"qt:MediaKind"`
	Screenwriters                string   `qt:"com.apple.proapps.share.screenWriter"                    xmp:"qt:Screenwriters"`
	SeasonNumber                 string   `qt:"com.apple.proapps.share.seasonNumber"                    xmp:"qt:SeasonNumber"`
	TVNetwork                    string   `qt:"com.apple.proapps.share.tvNetwork"                       xmp:"qt:TVNetwork"`
	Reel                         string   `qt:"com.apple.proapps.reel"                                  xmp:"qt:Reel"`
	CameraID                     string   `qt:"com.apple.proapps.cameraID"                              xmp:"qt:CameraID"`
	CameraManufacturer           string   `qt:"com.apple.proapps.manufacturer"                          xmp:"qt:CameraManufacturer"`
	CameraModel                  string   `qt:"com.apple.proapps.modelname"                             xmp:"qt:CameraModel"`
	CameraSerialNumber           string   `qt:"com.apple.proapps.serialno"                              xmp:"qt:CameraSerialNumber"`
	ClipFileName                 string   `qt:"com.apple.proapps.clipFileName"                          xmp:"qt:ClipFileName"`          // : A001C001_160721_D620
	AscCDL                       string   `qt:"com.apple.proapps.color.asc-cdl"                         xmp:"qt:AscCDL"`                // :
	IsGood                       string   `qt:"com.apple.proapps.isGood"                                xmp:"qt:IsGood"`                // : 0
	DataSize                     string   `qt:"com.apple.proapps.dataSize"                              xmp:"qt:DataSize"`              // : A�c?,
	DisplayFormat                string   `qt:"com.apple.proapps.displayFormat"                         xmp:"qt:DisplayFormat"`         // : 4k
	ShootingRate                 string   `qt:"com.apple.proapps.shootingRate"                          xmp:"qt:ShootingRate"`          // : @9
	VideoBitrate                 string   `qt:"com.apple.proapps.videoBitrate"                          xmp:"qt:VideoBitrate"`          // : 0
	Pulldown                     string   `qt:"com.apple.proapps.pulldown"                              xmp:"qt:Pulldown"`              // : 1
	MediaRate                    string   `qt:"com.apple.proapps.mediaRate"                             xmp:"qt:MediaRate"`             // : @9
	TimecodeFormat               string   `qt:"com.apple.proapps.timecodeFormat"                        xmp:"qt:TimecodeFormat"`        // : 2
	NumberOfAudioChannels        string   `qt:"com.apple.proapps.numberOfAudioChannels"                 xmp:"qt:NumberOfAudioChannels"` // : 0
	SampleRate                   string   `qt:"com.apple.proapps.sampleRate"                            xmp:"qt:SampleRate"`            // : 0
	BitsPerSample                string   `qt:"com.apple.proapps.bitsPerSample"                         xmp:"qt:BitsPerSample"`         // : 0
	StudioAlphaHandling          string   `qt:"com.apple.proapps.studio.alphaHandling"                  xmp:"qt:StudioAlphaHandling"`
	StudioCameraAngle            string   `qt:"com.apple.proapps.studio.angle"                          xmp:"qt:StudioCameraAngle"`
	StudioAnamorphicOverride     string   `qt:"com.apple.proapps.studio.metadataAnamorphicType"         xmp:"qt:StudioAnamorphicOverride"`
	StudioDeinterlace            string   `qt:"com.apple.proapps.studio.metadataDeinterlaceType"        xmp:"qt:StudioDeinterlace"`
	StudioFieldDominanceOverride string   `qt:"com.apple.proapps.studio.metadataFieldDominanceOverride" xmp:"qt:StudioFieldDominanceOverride"`
	StudioLocation               string   `qt:"com.apple.proapps.studio.metadataLocation"               xmp:"qt:StudioLocation"`
	StudioReel                   string   `qt:"com.apple.proapps.studio.reel"                           xmp:"qt:StudioReel"`
	StudioScene                  string   `qt:"com.apple.proapps.studio.scene"                          xmp:"qt:StudioScene"`
	StudioShot                   string   `qt:"com.apple.proapps.studio.shot"                           xmp:"qt:StudioTake"`
}

func (m *QtProApps) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsQuicktime}
}

func (m *QtProApps) Can(nsName string) bool {
	return nsName == NsQuicktime.GetName()
}

func (x *QtProApps) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x *QtProApps) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *QtProApps) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *QtProApps) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *QtProApps) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	}
	return nil
}
