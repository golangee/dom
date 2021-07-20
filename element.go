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

import (
	"strconv"
	"syscall/js"
)

const EventRelease = "forms-release"

// Element is always used in a value context, to avoid additional GC pressure.
// Therefore our state-handling relies on Javascript messaging (especially the
// release cycles).
type Element struct {
	val js.Value
}

func newElement(val js.Value) Element {
	return Element{val}
}

func (n Element) SetTextContent(v string) Element {
	n.val.Set("textContent", v)
	return n
}

func (n Element) SetInnerHTML(v string) Element {
	n.val.Set("innerHTML", v)
	return n
}

func (n Element) AppendTextNode(t string) Element {
	tn := GetWindow().Document().createTextNode(t)
	n.val.Call("appendChild", tn)
	return n
}

func (n Element) GetTagName() string {
	return n.val.Get("tagName").String()
}

func (n Element) SetClassName(str string) Element {
	//n.val.Set("className", str) // the property is not defined on non-html elements, like svg
	n.SetAttribute("class", str)
	return n
}

func (n Element) AppendElement(aChild Element) Element {
	n.val.Call("appendChild", aChild.val)
	return n
}

func (n Element) Clear() Element {
	for _, element := range n.Children() {
		element.Release()
	}
	return n.SetTextContent("")
}

func (n Element) Style() Style {
	return newStyle(n.val.Get("style"))
}

func (n Element) Children() []Element {
	var res []Element
	arr := n.val.Get("children")
	for i := 0; i < arr.Length(); i++ {
		res = append(res, newElement(arr.Index(i)))
	}
	return res
}

func (n Element) ReplaceWith(o Element) Element {
	n.val.Call("replaceWith", o.val)
	//	n.AppendElement(o)
	return o
}

func (n Element) String() string {
	v := n.val
	switch v.Type() {
	case js.TypeString:
		return v.String()
	case js.TypeBoolean:
		return strconv.FormatBool(v.Bool())
	case js.TypeNumber:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case js.TypeNull:
		return ""
	default:
		return n.val.Get("outerHTML").String()
	}
}

// Set sets the javascript property. Most standard attributes have an according property.
// See also https://javascript.info/dom-attributes-and-properties#html-attributes.
func (n Element) Set(p string, x interface{}) Element {
	n.val.Set(p, x)
	return n
}

// Get returns unwrapped primitives like string, bool or float but for others its return value is undefined and
// may change at any time.
func (n Element) Get(p string) interface{} {
	v := n.val.Get(p)
	return asGoValue(v)
}

func asGoValue(v js.Value) interface{} {
	switch v.Type() {
	case js.TypeString:
		return v.String()
	case js.TypeBoolean:
		return v.Bool()
	case js.TypeNumber:
		return v.Float()
	case js.TypeNull:
		return nil
	default:
		return v
	}
}

// SetAttribute sets the html attribute. There are attributes, which have no corresponding
// javascript property. See https://javascript.info/dom-attributes-and-properties#html-attributes.
func (n Element) SetAttribute(a string, x interface{}) Element {
	n.val.Call("setAttribute", a, x)
	return n
}

// GetAttribute returns unwrapped primitives like string, bool or float but for others its return value is undefined and
// may change at any time.
func (n Element) GetAttribute(p string) interface{} {
	v := n.val.Call("getAttribute", p)
	return asGoValue(v)
}

// RemoveAttribute deletes the html attribute.
func (n Element) RemoveAttribute(a string) Element {
	n.val.Call("removeAttribute", a)
	return n
}

// Call is an unclean abstraction and can invoke attached javascript methods.
// Args which are elements, are internally unwrapped.
func (n Element) Call(name string, args ...interface{}) Element {
	tmp := make([]interface{}, len(args))
	copy(tmp, args)
	for i, arg := range tmp {
		if elem, ok := arg.(Element); ok {
			tmp[i] = elem.val
		}
	}

	res := n.val.Call(name, tmp...)
	return newElement(res)
}

// The keydown event is fired when a key is pressed. See also
// https://developer.mozilla.org/en-US/docs/Web/API/Document/keydown_event
func (n Element) AddKeyListener(typ string, f func(keyCode int)) Releasable {
	return n.addEventListener(typ, false, func(this js.Value, args []js.Value) interface{} {
		f(args[0].Get("keyCode").Int())
		return nil
	})
}

// AddEventListener is internally very complex, because it keeps a global callback
// reference to connect the wasm and the javascript context. The wasm side must keep
// a global un-collectable function and the javascript side does the same. This makes
// event handling currently very expensive. Always ensure that
// you call Release on this Element to free all resources.
func (n Element) AddEventListener(typ string, once bool, listener func()) Releasable {
	return n.addEventListener(typ, once, func(this js.Value, args []js.Value) interface{} {
		args[0].Call("stopPropagation")
		if listener == nil {
			return nil
		}

		listener()
		return nil
	})
}

func (n Element) addEventListener(typ string, once bool, listener func(this js.Value, args []js.Value) interface{}) Releasable {
	defer GlobalPanicHandler()

	alreadyReleased := false
	var actualFunc, releaseFunc js.Func

	unregisterJS := func() {
		defer GlobalPanicHandler()

		if !alreadyReleased {
			alreadyReleased = true
			n.val.Call("removeEventListener", typ, actualFunc)
			n.val.Call("removeEventListener", EventRelease, releaseFunc)
			actualFunc.Release()
			releaseFunc.Release()
		}
	}

	actualFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer GlobalPanicHandler()

		if once {
			unregisterJS()
		}
		return listener(this, args)
	})

	releaseFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer GlobalPanicHandler()

		if !alreadyReleased {
			unregisterJS()
		}
		return nil
	})

	n.val.Call("addEventListener", EventRelease, releaseFunc)
	n.val.Call("addEventListener", typ, actualFunc)

	return actualFunc // TODO actually we should also release the release-funcs
}

// Release is part of our custom lifecycle. We need a manual destructor, because we
// have currently two running contexts: our wasm program and the browsers javascript interpreter.
// This is important to remove callbacks and other attached resources, which would otherwise
// never be freed, due to global or cyclic references. All contained children element, will
// also receive a Release call.
func (n Element) Release() {
	for _, element := range n.Children() {
		element.Release()
	}

	event := js.Global().Get("Event").New(EventRelease)
	n.val.Call("dispatchEvent", event)

	//	log.NewLogger().Print(ecs.Msg("releasing: "+n.String()))
}

func (n Element) AddReleaseListener(f func()) Element {
	var fun js.Func
	fun = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer GlobalPanicHandler()

		args[0].Call("stopPropagation")
		f()

		n.val.Call("removeEventListener", EventRelease, fun)
		fun.Release()

		//log.NewLogger().Print(ecs.Msg("plain releaser listener releasing: "+n.String()))
		return nil
	})

	n.val.Call("addEventListener", EventRelease, fun)
	return n
}

func (n Element) AddClass(v string) Element {
	n.val.Get("classList").Call("add", v)
	return n
}

func (n Element) HasClass(v string) bool {
	return n.val.Get("classList").Call("contains", v).Bool()
}

func (n Element) RemoveClass(v string) Element {
	n.val.Get("classList").Call("remove", v)
	return n
}

// Removes this Element from its parent and Release it. See also
// https://developer.mozilla.org/de/docs/Web/API/ChildNode/remove.
func (n Element) Remove() Element {
	n.val.Call("remove")
	n.Release()
	return n
}

// LastChild returns last child node. See also https://developer.mozilla.org/de/docs/Web/API/Node/lastChild.
func (n Element) LastChild() Element {
	return newElement(n.val.Get("lastChild"))
}

// Equal has the Javascript == semantic on the Element (equal reference)
func (n Element) Equal(o Element) bool {
	return n.val.Equal(o.val)
}

// IsNull returns true, if the represented Element is actually a null value.
func (n Element) IsNull() bool {
	return n.val.IsNull()
}

// GetID returns the empty string or the according unique element id. See also
// https://developer.mozilla.org/en-US/docs/Web/API/Element/id.
func (n Element) GetID() string {
	return n.val.Get("id").String()
}

// SetID update the according unique element id. See also
// https://developer.mozilla.org/en-US/docs/Web/API/Element/id.
func (n Element) SetID(id string) {
	n.val.Set("id", id)
}

// See https://developer.mozilla.org/en-US/docs/Web/API/Element/getBoundingClientRect.
func (n Element) GetBoundingClientRect() Rect {
	return newRect(n.val.Call("getBoundingClientRect"))
}
