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

// Sync across different models by specifying paths to values within each
// model. Paths must begin with a registered namespace short name followed
// by a colon (:). Paths must be the same strings as used in type's XMP tag,
// e.g
//
//     XMP Tag                      XMP Path
//     `xmp:"xmpDM:cameraAngle"`    xmpDM:cameraAngle
//
// To refer nested structs or arrays, paths may contain multiple segments
// separated by a forward slash (/). Array elements may be addressed using
// square brackets after the XMP name of the array, e.g.
//
//     XMP Array
//     dc:creator[0]       // the first array element in XMP lists and bags
//     dc:description[en]  // a language entry in XMP alternative arrays
//
// Sync flags
//
// Syncing can be controlled by flags. Flags are binary and can be combined to
// achieve a desired result. Some useful combinations are
//
//   S_CREATE|S_REPLACE|S_DELETE
// 	   create dest when not exist, overwrite when exist, clear when exist and
//     source is empty
//
//   S_CREATE|S_DELETE
// 	   create dest when not exist, clear when exist and source is empty
//
//   S_CREATE|S_REPLACE
//     create when not exist, overwrite when exist unless source is empty
//
// You may find other combinations helpful as well. When working with lists
// as sync destination, combine the above with S_APPEND or S_UNIQUE to extend
// a list. When S_DELETE is set and the source value is empty, the list will
// be cleared. With S_REPLACE, a list will be replaced by the source value,
// that is, afterwards the list contains a single entry only. Precedence order
// for slices/arrays when multiple flags are set is UNIQUE > APPEND > REPLACE.
//
// Examples
//   "xmpDM:cameraAngle" <-> "trim:Shot/Angle"
//   "tiff:Artist" <-> "dc:creator"
//
//

package xmp

import (
	"fmt"
	"strings"
)

type SyncDesc struct {
	Source  Path
	Dest    Path
	Flags   SyncFlags
	Convert ConverterFunc
}

type SyncDescList []*SyncDesc

type ConverterFunc func(string) string

type SyncFlags int

const (
	CREATE  SyncFlags = 1 << iota // create when not exist, nothing when exist or source is empty
	REPLACE                       // replace when dest exist, nothing when missing or source is empty
	DELETE                        // clear dest when source value is empty
	APPEND                        // list-only: append non-empty source value
	UNIQUE                        // list-only: append non-empty unique source value
	NOFAIL                        // don't fail when state+op+flags don't match
	DEFAULT = CREATE | REPLACE | DELETE | UNIQUE
	MERGE   = CREATE | UNIQUE | NOFAIL
	EXTEND  = CREATE | REPLACE | UNIQUE | NOFAIL
	ADD     = CREATE | UNIQUE | NOFAIL
)

func ParseSyncFlag(s string) SyncFlags {
	switch s {
	case "create":
		return CREATE
	case "replace":
		return REPLACE
	case "delete":
		return DELETE
	case "append":
		return APPEND
	case "unique":
		return UNIQUE
	case "nofail":
		return NOFAIL
	case "default":
		return DEFAULT
	case "merge":
		return MERGE
	case "extend":
		return EXTEND
	case "add":
		return ADD
	default:
		return 0
	}
}

func ParseSyncFlags(s string) (SyncFlags, error) {
	var flags SyncFlags
	for _, v := range strings.Split(s, ",") {
		f := ParseSyncFlag(strings.ToLower(v))
		if f == 0 {
			return 0, fmt.Errorf("invalid xmp sync flag '%s'", v)
		}
		flags |= f
	}
	return flags, nil
}

func (f *SyncFlags) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*f = 0
		return nil
	}
	if flags, err := ParseSyncFlags(string(data)); err != nil {
		return err
	} else {
		*f = flags
	}
	return nil
}

// Model is optional and exists for performance reasons to avoid lookups when
// a sync destination is a specific model only (useful to implement SyncFromXMP()
// in models)
func (d *Document) SyncMulti(desc SyncDescList, m Model) error {
	for _, v := range desc {
		if err := d.Sync(v.Source, v.Dest, v.Flags, m, v.Convert); err != nil {
			return err
		}
	}
	return nil
}

func (d *Document) Sync(sPath, dPath Path, flags SyncFlags, v Model, f ConverterFunc) error {
	// use default flags when zero
	if flags == 0 {
		flags = DEFAULT
	}

	// only XMP paths are supported here
	if !sPath.IsXmpPath() || !dPath.IsXmpPath() {
		return nil
	}

	sNs, _ := sPath.Namespace(d)
	dNs, _ := dPath.Namespace(d)

	// skip when either namespace does not exist
	if sNs == nil || dNs == nil {
		return nil
	}

	// skip when sPath model does not exist
	sModel := d.FindModel(sNs)
	if sModel == nil {
		return nil
	}

	dModel := v
	if dModel != nil {
		// dPath model must match dPath namespace
		if !dPath.MatchNamespace(dModel.Namespaces()[0]) {
			return nil
		}
	} else {
		dModel = d.FindModel(dNs)
	}

	// create dPath model
	if dModel == nil {
		if flags&CREATE == 0 {
			return nil
		}
		dModel = dNs.NewModel()
		d.AddModel(dModel)
	}

	sValue, err := GetModelPath(sModel, sPath)
	if err != nil {
		if flags&NOFAIL > 0 {
			return nil
		}
		return err
	}
	dValue, err := GetModelPath(dModel, dPath)
	if err != nil {
		if flags&NOFAIL > 0 {
			return nil
		}
		return err
	}

	// skip when equal
	if sValue == dValue {
		return nil
	}

	// empty source will only be used with delete flag
	if sValue == "" && flags&DELETE == 0 {
		return nil
	}

	// empty destination values require create flag
	if dValue == "" && flags&CREATE == 0 {
		return nil
	}

	// existing destination values require replace/delete/append/unique flag
	if dValue != "" && flags&(REPLACE|DELETE|APPEND|UNIQUE) == 0 {
		return nil
	}

	// convert the source value if requested
	if f != nil {
		sValue = f(sValue)
	}

	if err = SetModelPath(dModel, dPath, sValue, flags); err != nil {
		if flags&NOFAIL > 0 {
			return nil
		}
		return err
	}
	d.SetDirty()
	return nil
}

func (d *Document) Merge(b *Document, flags SyncFlags) error {
	p, err := b.ListPaths()
	if err != nil {
		return err
	}

	// copy namespaces
	for n, v := range b.intNsMap {
		d.intNsMap[n] = v
	}
	for n, v := range b.extNsMap {
		d.extNsMap[n] = v
	}

	// copy content
	for _, v := range p {
		v.Flags = flags
		if err := d.SetPath(v); err != nil {
			return err
		}
	}

	d.SetDirty()
	return nil
}
