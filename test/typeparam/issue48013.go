// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"unsafe"
)

type S[T any] struct {
	val T
}

// Test type substitution where base type is unsafe.Pointer
type U[T any] unsafe.Pointer

func test[T any]() T {
	var q U[T]
	var v struct {
		// Test derived type that contains an unsafe.Pointer
		p   unsafe.Pointer
		val T
	}
	_ = q
	return v.val
}

func main() {
	want := 0
	golangt := test[int]()
	if golangt != want {
		panic(fmt.Sprintf("golangt %f, want %f", golangt, want))
	}

}
