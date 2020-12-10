// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"net/url"
	"strconv"
)

// Query represents a routing query and provides a simple access to query based parameters.
type Query struct {
	path   string
	values url.Values
}

func (p Query) Path() string {
	return p.path
}

func (p Query) Get(key string) string {
	return p.values.Get(key)
}

func (p Query) Int(key string) int {
	i, _ := strconv.ParseInt(p.Get(key), 10, 64)
	return int(i)
}

func (p Query) Bool(key string) bool {
	b, _ := strconv.ParseBool(p.Get(key))
	return b
}

func (p Query) Float64(key string) float64 {
	f, _ := strconv.ParseFloat(p.Get(key), 64)
	return f
}

// GetAll returns all defined query parameters with that key.
func (p Query) GetAll(key string) []string {
	return p.values[key]
}
