package dom

import "syscall/js"

type Location struct {
	val js.Value
}

func newLocation(val js.Value) Location {
	return Location{val: val}
}

func (n Location) Reload(force bool) Location {
	n.val.Call("reload", force)
	return n
}

func (n Location) Href() string {
	return n.val.Get("href").String()
}
