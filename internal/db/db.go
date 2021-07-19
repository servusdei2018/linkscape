// Copyright 2021 The Linkscape Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"github.com/patrickmn/go-cache"
)

// DB facilitates interaction with an underlying cache.
// DB and its methods can be safely used by multiple goroutines.
type DB struct {
	*cache.Cache
}

// New instantiates a new database.
func New(cfg *Config) *DB {
	return &DB{cache.New(cfg.Lifetime, cfg.Frequency)}
}

// Get retrieves an entry from the database.
func (db *DB) Get(key string) (*string, bool) {
	val, found := db.Cache.Get(key)
	if !found {
		return nil, false
	}
	return val.(*string), true
}

// Set stores an entry in the database.
func (db *DB) Set(key string, val *string) {
	db.Cache.Set(key, val, cache.DefaultExpiration)
}
