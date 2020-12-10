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

// A Rect represents a DOMRect see also https://developer.mozilla.org/en-US/docs/Web/API/DOMRect. It seems undefined,
// what nature the coordinates are.
type Rect struct {
	val js.Value
}

func newRect(val js.Value) Rect {
	return Rect{val}
}

func (n Rect) GetX() int {
	return n.val.Get("x").Int()
}

func (n Rect) SetX(x int) Rect {
	n.val.Set("x", x)
	return n
}

func (n Rect) GetY() int {
	return n.val.Get("y").Int()
}

func (n Rect) SetY(y int) Rect {
	n.val.Set("Y", y)
	return n
}

func (n Rect) GetWidth() int {
	return n.val.Get("width").Int()
}

func (n Rect) SetWidth(width int) Rect {
	n.val.Set("width", width)
	return n
}

func (n Rect) GetHeight() int {
	return n.val.Get("height").Int()
}

func (n Rect) SetHeight(height int) Rect {
	n.val.Set("height", height)
	return n
}

func (n Rect) GetTop() int {
	return n.val.Get("top").Int()
}

func (n Rect) SetTop(t int) Rect {
	n.val.Set("top", t)
	return n
}

func (n Rect) GetRight() int {
	return n.val.Get("right").Int()
}

func (n Rect) SetRight(r int) Rect {
	n.val.Set("right", r)
	return n
}

func (n Rect) GetBottom() int {
	return n.val.Get("bottom").Int()
}

func (n Rect) SetBottom(b int) Rect {
	n.val.Set("bottom", b)
	return n
}

func (n Rect) GetLeft() int {
	return n.val.Get("left").Int()
}

func (n Rect) SetLeft(l int) Rect {
	n.val.Set("left", l)
	return n
}
