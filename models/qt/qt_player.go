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

package qt

import (
	"fmt"
	"trimmer.io/go-xmp/xmp"
)

type QtPlayer struct {
	AudioBalance     float32 `qt:"com.apple.quicktime.player.movie.audio.balance"     xmp:"qt:AudioBalance"`
	AudioBass        float32 `qt:"com.apple.quicktime.player.movie.audio.bass"        xmp:"qt:AudioBass"`
	AudioGain        float32 `qt:"com.apple.quicktime.player.movie.audio.gain"        xmp:"qt:AudioGain"`
	AudioMute        Bool    `qt:"com.apple.quicktime.player.movie.audio.mute"        xmp:"qt:AudioMute"`
	AudioPitchshift  float32 `qt:"com.apple.quicktime.player.movie.audio.pitchshift"  xmp:"qt:AudioPitchshift"`
	AudioTreble      float32 `qt:"com.apple.quicktime.player.movie.audio.treble"      xmp:"qt:AudioTreble"`
	VisualBrightness float32 `qt:"com.apple.quicktime.player.movie.visual.brightness" xmp:"qt:VisualBrightness"`
	VisualColor      float32 `qt:"com.apple.quicktime.player.movie.visual.color"      xmp:"qt:VisualColor"`
	VisualContrast   float32 `qt:"com.apple.quicktime.player.movie.visual.contrast"   xmp:"qt:VisualContrast"`
	VisualTint       float32 `qt:"com.apple.quicktime.player.movie.visual.tint"       xmp:"qt:VisualTint"`
	Version          string  `qt:"com.apple.quicktime.player.version"                 xmp:"qt:Version"`
}

func (m *QtPlayer) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsQuicktime}
}

func (m *QtPlayer) Can(nsName string) bool {
	return nsName == NsQuicktime.GetName()
}

func (x *QtPlayer) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *QtPlayer) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x *QtPlayer) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *QtPlayer) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *QtPlayer) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *QtPlayer) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsQuicktime.GetName(), err)
	}
	return nil
}
