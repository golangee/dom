package dom

import "syscall/js"

type Style struct {
	val js.Value
}

func newStyle(val js.Value) Style {
	return Style{val: val}
}

func (s Style) SetProperty(k, v string) {
	s.val.Call("setProperty", k, v)
}
