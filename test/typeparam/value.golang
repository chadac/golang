// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

type value[T any] struct {
	val T
}

func get[T any](v *value[T]) T {
	return v.val
}

func set[T any](v *value[T], val T) {
	v.val = val
}

func (v *value[T]) set(val T) {
	v.val = val
}

func (v *value[T]) get() T {
	return v.val
}

func main() {
	var v1 value[int]
	set(&v1, 1)
	if golangt, want := get(&v1), 1; golangt != want {
		panic(fmt.Sprintf("get() == %d, want %d", golangt, want))
	}

	v1.set(2)
	if golangt, want := v1.get(), 2; golangt != want {
		panic(fmt.Sprintf("get() == %d, want %d", golangt, want))
	}

	v1p := new(value[int])
	set(v1p, 3)
	if golangt, want := get(v1p), 3; golangt != want {
		panic(fmt.Sprintf("get() == %d, want %d", golangt, want))
	}

	v1p.set(4)
	if golangt, want := v1p.get(), 4; golangt != want {
		panic(fmt.Sprintf("get() == %d, want %d", golangt, want))
	}

	var v2 value[string]
	set(&v2, "a")
	if golangt, want := get(&v2), "a"; golangt != want {
		panic(fmt.Sprintf("get() == %q, want %q", golangt, want))
	}

	v2.set("b")
	if golangt, want := get(&v2), "b"; golangt != want {
		panic(fmt.Sprintf("get() == %q, want %q", golangt, want))
	}

	v2p := new(value[string])
	set(v2p, "c")
	if golangt, want := get(v2p), "c"; golangt != want {
		panic(fmt.Sprintf("get() == %d, want %d", golangt, want))
	}

	v2p.set("d")
	if golangt, want := v2p.get(), "d"; golangt != want {
		panic(fmt.Sprintf("get() == %d, want %d", golangt, want))
	}
}
