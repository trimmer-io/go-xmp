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
	"strings"
	"sync"
)

type ModelFactory func(name string) Model

// interface for data models
type Model interface {
	Can(nsName string) bool
	Namespaces() NamespaceList
	SyncFromXMP(d *Document) error
	SyncToXMP(d *Document) error
	SyncModel(d *Document) error
	CanTag(tag string) bool
	GetTag(tag string) (string, error)
	SetTag(tag, value string) error
}

func emptyFactory(name string) Model {
	return nil
}

type Registry struct {
	nsNameMap map[string]*Namespace
	nsUriMap  map[string]*Namespace
	groupMap  map[NamespaceGroup]NamespaceList
	m         sync.RWMutex
}

var NsRegistry Registry = Registry{
	nsNameMap: make(map[string]*Namespace),
	nsUriMap:  make(map[string]*Namespace),
	groupMap:  make(map[NamespaceGroup]NamespaceList),
}

func Register(ns *Namespace, groups ...NamespaceGroup) {
	NsRegistry.RegisterNamespace(ns, groups)
}

func GetNamespace(prefix string) (*Namespace, error) {
	return NsRegistry.GetNamespace(prefix)
}

func GetGroupNamespaces(group NamespaceGroup) (NamespaceList, error) {
	return NsRegistry.GetGroupNamespaces(group)
}

func (r *Registry) RegisterNamespace(ns *Namespace, groups NamespaceGroupList) {
	r.m.Lock()
	defer r.m.Unlock()
	r.nsNameMap[ns.GetName()] = ns
	r.nsUriMap[ns.GetURI()] = ns
	for _, v := range groups {
		if v != NoMetadata {
			g, ok := r.groupMap[v]
			if !ok {
				r.groupMap[v] = NamespaceList{ns}
			} else {
				r.groupMap[v] = append(g, ns)
			}
		}
	}
}

func (r *Registry) GetGroupNamespaces(group NamespaceGroup) (NamespaceList, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if l, ok := r.groupMap[group]; ok {
		v := make(NamespaceList, len(l))
		copy(v, l)
		return v, nil
	}
	return nil, fmt.Errorf("xmp: unregistered namespace group '%s'", string(group))
}

func (r *Registry) GetNamespace(prefix string) (*Namespace, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if ns, ok := r.nsNameMap[prefix]; ok {
		return ns, nil
	}
	return nil, fmt.Errorf("xmp: unregistered namespace '%s'", prefix)
}

func (r *Registry) GetPrefix(uri string) string {
	r.m.RLock()
	defer r.m.RUnlock()
	if ns, ok := r.nsUriMap[uri]; ok {
		return ns.GetName()
	}
	return ""
}

func (r *Registry) Short(uri, name string) string {
	pre := r.GetPrefix(uri)
	if pre != "" {
		return strings.Join([]string{pre, name}, ":")
	}
	return name
}

func (r *Registry) Namespaces() NamespaceList {
	r.m.RLock()
	defer r.m.RUnlock()
	l := make(NamespaceList, 0, len(r.nsUriMap))
	for _, v := range r.nsUriMap {
		l = append(l, v)
	}
	return l
}

func (r *Registry) Prefixes() []string {
	r.m.RLock()
	defer r.m.RUnlock()
	l := make([]string, 0, len(r.nsNameMap))
	for n, _ := range r.nsNameMap {
		l = append(l, n)
	}
	return l
}
