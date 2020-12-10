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
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

// The GlobalPanicHandler should be called with a defer in every method or callback which
// will likely cause a panic. A non-recovered panic will cause the wasm-Code to just
// exit silently, which is in practice not very helpful. You probably want to do some
// event logging or just show a support screen. Without that, the user may try to continue
// interacting with an already dead application, which must be avoided in all cases, to
// ensure usability.
//
// The default implementation will try to recover a panic and replaces the body content
// with a readable stack trace. The formatting uses some tailwind css classes.
var GlobalPanicHandler = func() {
	r := recover()
	if r == nil {
		return
	}
	msg := fmt.Sprint(r)

	log.Println(msg, string(debug.Stack()))
	body := GetWindow().Document().Body()
	body.Clear()
	body.SetClassName("bg-gray-300")
	body.AppendElement(formatPanic(msg))
}

func formatPanic(msg string) Element {
	doc := GetWindow().Document()
	frame := doc.CreateElement("div").SetClassName("bg-white max-w-6xl p-1 m-10 rounded overflow-hidden shadow-lg dark:bg-gray-800")
	doc.Body().AppendElement(frame)

	title := doc.CreateElement("p").SetClassName("text-xl text-red-600")
	title.SetTextContent(msg)
	frame.AppendElement(title)

	stackLines := strings.Split(string(debug.Stack()), "\n")
	for i, line := range stackLines {
		lEl := doc.CreateElement("p").SetClassName("text-base")
		lEl.SetTextContent(line)
		frame.AppendElement(lEl)

		if strings.Contains(line, ".go:") {
			lEl.AddClass("text-red-600")
			lEl.AddClass("medium")
		} else {
			if i > 0 {
				lEl.AddClass("text-gray-400")
			}
		}
	}

	return frame
}
