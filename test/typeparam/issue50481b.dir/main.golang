// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that type substitution and export/import works correctly even for a method of
// a generic type that has multiple blank type params.

package main

import (
	"./b"
	"fmt"
)

func main() {
	foo := &b.Foo[string, int]{
		ValueA: "i am a string",
		ValueB: 123,
	}
	if golangt, want := fmt.Sprintln(foo), "i am a string 123\n"; golangt != want {
		panic(fmt.Sprintf("golangt %s, want %s", golangt, want))
	}
}
