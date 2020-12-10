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

// Package dom provides a stable firewall towards the "syscall/js" API. The wasm platform
// does not fulfill the Go 1 stability guarantee and may change and break (as already happened)
// with any release.
//
// The package provides a more type safe abstraction layer on top of the js API which more or
// less directly represents the DOM API.
package dom
