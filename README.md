# dom [![GoDoc](https://godoc.org/github.com/golangee/dom?status.svg)](http://godoc.org/github.com/golangee/dom)
Package dom provides a *stable firewall* towards the "syscall/js" API. The wasm platform
does not fulfill the Go 1 stability guarantee and may change and break (as already happened)
with any release.

The package provides a more type safe abstraction layer on top of the js API which more or
less directly represents the DOM API.

It is important to note, that there is a little custom lifecycle semantic fo some listeners,
to make releasing less leaky and easier to use.

## Roadmap
It is planned to replace this handwritten layer with an automatically generated version of the
[HTML living standard](https://html.spec.whatwg.org/) which is based on a bikeshedding format
with embedded *webidl* descriptions, see also https://github.com/whatwg/dom/blob/master/review-drafts/2020-06.bs.
In contrast to already available generated bindings for Go and other languages like Rust, the human
readable documentation should be in the mix.