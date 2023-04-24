// Copyright 2023 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// If we decide to employ this auto generation of markdown documentation for
// amtool and alertmanager, this package could potentially be moved to
// prometheus/common. However, it is crucial to note that this functionality is
// tailored specifically to the way in which the Prometheus documentation is
// rendered, and should be avoided for use by third-party users.

package cache

import (
	"sync"
)

const (
	shardSize = 128
)

var DefautlDictCache = NewDictCache()

type dictValue struct {
	value string
}

type DictCache struct {
	dbk   []map[string]*dictValue
	locks []*sync.RWMutex
}

func NewDictCache() *DictCache {
	dc := DictCache{
		dbk:   make([]map[string]*dictValue, shardSize),
		locks: make([]*sync.RWMutex, shardSize),
	}
	for i := 0; i < shardSize; i++ {
		dc.dbk[i] = make(map[string]*dictValue)
		dc.locks[i] = &sync.RWMutex{}
	}
	return &dc
}

func (d *DictCache) shard(key string) int {
	l := len(key)
	if l == 0 {
		return 0
	}
	s := int(key[0])
	s <<= 8
	s += int(key[l-1])
	s %= shardSize
	return s
}

func (d *DictCache) Get(word string) string {
	l := len(word)
	if l == 0 {
		return ""
	}
	s := d.shard(word)

	d.locks[s].RLock()
	if v, ok := d.dbk[s][word]; ok {
		d.locks[s].RUnlock()
		return v.value
	}
	d.locks[s].RUnlock()

	d.locks[s].Lock()
	v := &dictValue{
		value: word,
	}
	d.dbk[s][word] = v
	d.locks[s].Unlock()

	return word
}
