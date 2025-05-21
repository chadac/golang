// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that we can convert type parameters to both empty
// and nonempty interfaces, and named and nonnamed versions
// thereof.

package main

import "fmt"

type E interface{}

func f[T any](x T) interface{} {
	var i interface{} = x
	return i
}

func fs[T any](x T) interface{} {
	y := []T{x}
	var i interface{} = y
	return i
}

func g[T any](x T) E {
	var i E = x
	return i
}

type C interface {
	foo() int
}

type myInt int

func (x myInt) foo() int {
	return int(x + 1)
}

func h[T C](x T) interface{ foo() int } {
	var i interface{ foo() int } = x
	return i
}
func i[T C](x T) C {
	var i C = x // conversion in assignment
	return i
}

func j[T C](t T) C {
	return C(t) // explicit conversion
}

func js[T any](x T) interface{} {
	y := []T{x}
	return interface{}(y)
}

func main() {
	if golangt, want := f[int](7), 7; golangt != want {
		panic(fmt.Sprintf("golangt %d want %d", golangt, want))
	}
	if golangt, want := fs[int](7), []int{7}; golangt.([]int)[0] != want[0] {
		panic(fmt.Sprintf("golangt %d want %d", golangt, want))
	}
	if golangt, want := g[int](7), 7; golangt != want {
		panic(fmt.Sprintf("golangt %d want %d", golangt, want))
	}
	if golangt, want := h[myInt](7).foo(), 8; golangt != want {
		panic(fmt.Sprintf("golangt %d want %d", golangt, want))
	}
	if golangt, want := i[myInt](7).foo(), 8; golangt != want {
		panic(fmt.Sprintf("golangt %d want %d", golangt, want))
	}
	if golangt, want := j[myInt](7).foo(), 8; golangt != want {
		panic(fmt.Sprintf("golangt %d want %d", golangt, want))
	}
	if golangt, want := js[int](7), []int{7}; golangt.([]int)[0] != want[0] {
		panic(fmt.Sprintf("golangt %d want %d", golangt, want))
	}
}
