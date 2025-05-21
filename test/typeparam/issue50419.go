// run

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that type substitution works correctly even for a method of a generic type
// that has multiple blank type params.

package main

import (
	"fmt"
)

func main() {
	foo := &Foo[string, int]{
		valueA: "i am a string",
		valueB: 123,
	}
	if golangt, want := fmt.Sprintln(foo), "i am a string 123\n"; golangt != want {
		panic(fmt.Sprintf("golangt %s, want %s", golangt, want))
	}
}

type Foo[T1 any, T2 any] struct {
	valueA T1
	valueB T2
}

func (f *Foo[_, _]) String() string {
	return fmt.Sprintf("%v %v", f.valueA, f.valueB)
}
