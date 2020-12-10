// Package dom provides a stable firewall towards the "syscall/js" API. The wasm platform
// does not fulfill the Go 1 stability guarantee and may change and break (as already happened)
// with any release.
//
// The package provides a more type safe abstraction layer on top of the js API which more or
// less directly represents the DOM API.
package dom
