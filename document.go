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

package dom

import "syscall/js"

type Document struct {
	val js.Value
}

func newDocument(val js.Value) Document {
	return Document{val}
}

func GetDocument() Document {
	return GetWindow().Document()
}

func (n Document) Body() Element {
	return newElement(n.val.Get("body"))
}

// In an HTML document, the document.createElement() method creates the HTML element specified by tagName, or an HTMLUnknownElement if tagName isn't recognized.
func (n Document) createTextNode(name string) js.Value {
	return n.val.Call("createTextNode", name)
}

// In an HTML document, the document.createElement() method creates the HTML element specified by tagName, or an HTMLUnknownElement if tagName isn't recognized.
func (n Document) CreateElement(name string) Element {
	v := n.val.Call("createElement", name)
	return newElement(v)
}

// In an HTML document, the document.createElement() method creates the HTML element specified by tagName, or an HTMLUnknownElement if tagName isn't recognized.
func (n Document) CreateElementNS(ns string, name string) Element {
	v := n.val.Call("createElementNS", ns, name)
	return newElement(v)
}

func (n Document) DocumentElement() Element {
	body := n.val.Get("documentElement")
	return newElement(body)
}

// GeElementById follows https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementById.
// If not found, returns an Element whose Element.IsNull method will return true.
func (n Document) GetElementById(id string) Element {
	return newElement(n.val.Call("getElementById", id))
}
