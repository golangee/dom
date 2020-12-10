package dom

import "syscall/js"

// A Releasable cleans up references and the resource must not be used afterwards anymore.
type Releasable interface {
	Release()
}

type ReleaseFunc func()

func (f ReleaseFunc) Release() {
	f()
}

type Window struct {
	val js.Value
}

func newWindow() Window {
	n := js.Global().Get("window")
	return Window{n}
}

func GetWindow() Window {
	return newWindow()
}

// Document returns a reference to the document contained in the window.
func (n Window) Document() Document {
	return newDocument(n.val.Get("document"))
}

func (n Window) OnHashChange(f func()) Releasable {
	fun := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f()
		return nil
	})
	n.val.Set("onhashchange", fun)

	return fun
}

func (n Window) HashChange(f func()) Releasable {
	fun := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f()
		return nil
	})
	n.val.Set("hashchange", fun)

	return fun
}

func (n Window) OnPopState(f func()) Releasable {
	fun := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f()
		return nil
	})
	n.val.Set("onpopstate", fun)

	return fun
}

func (n Window) Location() Location {
	return newLocation(n.val.Get("location"))
}

func (n Window) SetLocation(url string) Window {
	n.val.Set("location", url)
	return n
}

// AddEventListener registers a new listener. Note that there is no automatic release management as for Element.
func (n Window) AddEventListener(typ string, listener func()) Releasable {
	actualFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer GlobalPanicHandler()

		listener()
		return nil
	})

	n.val.Call("addEventListener", typ, actualFunc)

	return ReleaseFunc(func() {
		n.val.Call("removeEventListener", typ, actualFunc)
		actualFunc.Release()
	})
}

// MatchesMedia is the javascript equivalent to css media queries. criteria is for example
//  - (min-width:800px)
//  - (min-width:800px) or (orientation: landscape)
//  - (max-width: 800px)
func (n Window) MatchesMedia(criteria string) bool {
	return n.val.Call("matchMedia", criteria).Get("matches").Bool()
}

func (n Window) InnerWidth() int {
	return n.val.Get("innerWidth").Int()
}

func (n Window) InnerHeight() int {
	return n.val.Get("innerHeight").Int()
}
