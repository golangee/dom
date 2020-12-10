package dom

import (
	"syscall/js"
	"time"
)

func SetTimeout(d time.Duration, f func()) {
	var wrapper js.Func
	wrapper = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer GlobalPanicHandler()
		defer wrapper.Release()
		f()
		return nil
	})
	GetWindow().val.Call("setTimeout", wrapper, d.Milliseconds())
}

func Post(f func()) {
	SetTimeout(0, f)
}

