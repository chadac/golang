// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test wrappers/interfaces for generic type embedding another generic type.

package main

import "fmt"

type A[T any] struct {
	B[T]
}

type B[T any] struct {
	val T
}

func (b *B[T]) get() T {
	return b.val
}

type getter[T any] interface {
	get() T
}

//golang:noinline
func doGet[T any](i getter[T]) T {
	return i.get()
}

//golang:noline
func doGet2[T any](i interface{}) T {
	i2 := i.(getter[T])
	return i2.get()
}

func main() {
	a := A[int]{B: B[int]{3}}
	var i getter[int] = &a

	if golangt, want := doGet(i), 3; golangt != want {
		panic(fmt.Sprintf("golangt %v, want %v", golangt, want))
	}

	as := A[string]{B: B[string]{"abc"}}
	if golangt, want := doGet2[string](&as), "abc"; golangt != want {
		panic(fmt.Sprintf("golangt %v, want %v", golangt, want))
	}
}
