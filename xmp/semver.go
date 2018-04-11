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
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func ParseVersion(s string) (Version, error) {
	var v Version
	tokens := strings.Split(s, ".")
	switch len(tokens) {
	case 3:
		i, err := strconv.Atoi(tokens[2])
		if err != nil {
			return v, fmt.Errorf("semver: illegal version string '%s': %v", s, err)
		}
		v.Patch = i
		fallthrough
	case 2:
		i, err := strconv.Atoi(tokens[1])
		if err != nil {
			return v, fmt.Errorf("semver: illegal version string '%s': %v", s, err)
		}
		v.Minor = i
		fallthrough
	case 1:
		i, err := strconv.Atoi(tokens[0])
		if err != nil {
			return v, fmt.Errorf("semver: illegal version string '%s': %v", s, err)
		}
		v.Major = i
	default:
		return v, fmt.Errorf("semver: illegal version string '%s'", s)
	}
	return v, nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) IsZero() bool {
	return v.Major+v.Minor+v.Patch == 0
}

func (v Version) Between(a, b Version) bool {
	switch {
	case v.IsZero():
		return true
	case a.IsZero():
		return v.Before(b) || v.Equal(a)
	case b.IsZero():
		return v.After(a) || v.Equal(a)
	default:
		return v.Equal(a) || v.Equal(b) || (v.After(a) && v.Before(b))
	}
}

func (v Version) Equal(a Version) bool {
	return v.Major == a.Major && v.Minor == a.Minor && v.Patch == a.Patch
}

func (v Version) Before(a Version) bool {
	return (v.Major < a.Major) ||
		(v.Major == a.Major && v.Minor < a.Minor) ||
		(v.Major == a.Major && v.Minor == a.Minor && v.Patch < a.Patch)
}

func (v Version) After(a Version) bool {
	return (v.Major > a.Major) ||
		(v.Major == a.Major && v.Minor > a.Minor) ||
		(v.Major == a.Major && v.Minor == a.Minor && v.Patch > a.Patch)
}
