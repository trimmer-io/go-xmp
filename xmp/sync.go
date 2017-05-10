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
	S_CREATE  SyncFlags = 1 << iota // create when not exist, nothing when exist or source is empty
	S_REPLACE                       // replace when dest exist, nothing when missing or source is empty
	S_DELETE                        // clear dest when source value is empty
	S_APPEND                        // list-only: append non-empty source value
	S_UNIQUE                        // list-only: append non-empty unique source value
	S_DEFAULT = S_CREATE | S_REPLACE | S_DELETE | S_UNIQUE
	S_MERGE   = S_CREATE | S_REPLACE | S_UNIQUE
)

func ParseSyncFlag(s string) SyncFlags {
	switch s {
	case "create":
		return S_CREATE
	case "replace":
		return S_REPLACE
	case "delete":
		return S_DELETE
	case "append":
		return S_APPEND
	case "unique":
		return S_UNIQUE
	case "default":
		return S_DEFAULT
	case "merge":
		return S_MERGE
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
		flags = S_DEFAULT
	}

	// only XMP paths are supported here
	if !sPath.IsXmpPath() || !dPath.IsXmpPath() {
		return nil
	}

	sNs, _ := sPath.Namespace()
	dNs, _ := dPath.Namespace()

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
		if flags&S_CREATE == 0 {
			return nil
		}
		dModel = dNs.NewModel()
		d.AddModel(dModel)
	}

	sValue, err := GetPath(sModel, sPath)
	if err != nil {
		return err
	}
	dValue, err := GetPath(dModel, dPath)
	if err != nil {
		return err
	}

	// skip when equal
	if sValue == dValue {
		return nil
	}

	// empty source will only be used with delete flag
	if sValue == "" && flags&S_DELETE == 0 {
		return nil
	}

	// empty destination values require create flag
	if dValue == "" && flags&S_CREATE == 0 {
		return nil
	}

	// existing destination values require replace/delete/append/unique flag
	if dValue != "" && flags&(S_REPLACE|S_DELETE|S_APPEND|S_UNIQUE) == 0 {
		return nil
	}

	// convert the source value if requested
	if f != nil {
		sValue = f(sValue)
	}

	if err := SetPath(dModel, dPath, sValue, flags); err != nil {
		return err
	}
	if err == nil {
		d.SetDirty()
	}

	return nil
}
