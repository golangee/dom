package dom

import (
	"strconv"
	"sync/atomic"
	"syscall/js"
)

var idCounter = int32(0)

// GenerateID returns the next unique identifier for dom elements.
func GenerateID() string {
	v := atomic.AddInt32(&idCounter, 1)
	return strconv.Itoa(int(v))
}

type Global struct {
	val js.Value
}

func GetGlobal() Global {
	return Global{val: js.Global()}
}

func (n Global) Get(name string) Element {
	return newElement(n.val.Get(name))
}
