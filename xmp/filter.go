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
	"reflect"
	"strings"
)

type Filter struct {
	include NamespaceList
	exclude NamespaceList
}

func NewFilter(include, exclude NamespaceList) *Filter {
	return &Filter{
		include: include,
		exclude: exclude,
	}
}

func ParseFilterStrict(s string) (*Filter, error) {
	f := NewFilter(nil, nil)
	if len(s) == 0 {
		return f, nil
	}
	for _, v := range strings.Split(s, ",") {
		if v == "" {
			continue
		}
		switch v[0] {
		case '-':
			v = v[1:]
			// 1st try parsing as group
			if g := ParseNamespaceGroup(v); g != NoMetadata {
				f.exclude = append(f.exclude, g.Namespaces()...)
				break
			}
			// 2nd try as namespace
			if ns, err := GetNamespace(v); err == nil {
				f.exclude = append(f.exclude, ns)
			} else {
				return nil, fmt.Errorf("no registered xmp namespace %s", v)
			}
		case '+':
			v = v[1:]
			fallthrough
		default:
			// 1st try parsing as group
			if g := ParseNamespaceGroup(v); g != NoMetadata {
				f.include = append(f.include, g.Namespaces()...)
				break
			}
			// 2nd try as namespace
			if ns, err := GetNamespace(v); err == nil {
				f.include = append(f.include, ns)
			} else {
				return nil, fmt.Errorf("no registered xmp namespace %s", v)
			}
		}
	}
	return f, nil
}

func ParseFilter(s string) *Filter {
	f := NewFilter(nil, nil)
	if len(s) == 0 {
		return f
	}
	for _, v := range strings.Split(s, ",") {
		if v == "" {
			continue
		}
		switch v[0] {
		case '-':
			v = v[1:]
			// 1st try parsing as group
			if g := ParseNamespaceGroup(v); g != NoMetadata {
				f.exclude = append(f.exclude, g.Namespaces()...)
				break
			}
			// 2nd try as namespace
			if ns, err := GetNamespace(v); err == nil {
				f.exclude = append(f.exclude, ns)
			} else {
				f.exclude = append(f.exclude, &Namespace{Name: v})
			}
		case '+', ' ': // treat space as '+' to because of URL unescaping
			v = v[1:]
			fallthrough
		default:
			// 1st try parsing as group
			if g := ParseNamespaceGroup(v); g != NoMetadata {
				f.include = append(f.include, g.Namespaces()...)
				break
			}
			// 2nd try as namespace
			if ns, err := GetNamespace(v); err == nil {
				f.include = append(f.include, ns)
			} else {
				f.include = append(f.include, &Namespace{Name: v})
			}
		}
	}
	return f
}

func (x *Filter) UnmarshalText(b []byte) error {
	f := ParseFilter(string(b))
	*x = *f
	return nil
}

func (x Filter) String() string {
	s := make([]string, 0, len(x.include)+len(x.exclude))
	for _, v := range x.include {
		s = append(s, fmt.Sprintf("+%s", v.GetName()))
	}
	for _, v := range x.exclude {
		s = append(s, fmt.Sprintf("-%s", v.GetName()))
	}
	return strings.Join(s, ",")
}

func (x Filter) Apply(d *Document) bool {
	var removed bool
	if len(x.include) > 0 {
		r := d.FilterNamespaces(x.include)
		removed = removed || r
	}
	if len(x.exclude) > 0 {
		r := d.RemoveNamespaces(x.exclude)
		removed = removed || r
	}
	return removed
}

// converter function for Gorilla schema; will become unnecessary once
// https://github.com/gorilla/schema/issues/57 is fixed
//
// register with decoder.RegisterConverter(Filter(""), ConvertFilter)
func ConvertFilter(value string) reflect.Value {
	v := ParseFilter(value)
	return reflect.ValueOf(*v)
}
