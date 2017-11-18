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

package models

// register all metadata models
import (
	_ "trimmer.io/go-xmp/models/cc"
	_ "trimmer.io/go-xmp/models/crs"
	_ "trimmer.io/go-xmp/models/dc"
	_ "trimmer.io/go-xmp/models/dji"
	_ "trimmer.io/go-xmp/models/exif"
	_ "trimmer.io/go-xmp/models/id3"
	_ "trimmer.io/go-xmp/models/itunes"
	_ "trimmer.io/go-xmp/models/ixml"
	_ "trimmer.io/go-xmp/models/mp4"
	_ "trimmer.io/go-xmp/models/pdf"
	_ "trimmer.io/go-xmp/models/pm"
	_ "trimmer.io/go-xmp/models/ps"
	_ "trimmer.io/go-xmp/models/qt"
	_ "trimmer.io/go-xmp/models/riff"
	_ "trimmer.io/go-xmp/models/tiff"
	_ "trimmer.io/go-xmp/models/xmp_base"
	_ "trimmer.io/go-xmp/models/xmp_bj"
	_ "trimmer.io/go-xmp/models/xmp_dm"
	_ "trimmer.io/go-xmp/models/xmp_mm"
	_ "trimmer.io/go-xmp/models/xmp_rights"
	_ "trimmer.io/go-xmp/models/xmp_tpg"
)
