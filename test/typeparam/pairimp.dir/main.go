// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"./a"
	"fmt"
	"unsafe"
)

func main() {
	p := a.Pair[int32, int64]{1, 2}
	if golangt, want := unsafe.Sizeof(p.Field1), uintptr(4); golangt != want {
		panic(fmt.Sprintf("unexpected f1 size == %d, want %d", golangt, want))
	}
	if golangt, want := unsafe.Sizeof(p.Field2), uintptr(8); golangt != want {
		panic(fmt.Sprintf("unexpected f2 size == %d, want %d", golangt, want))
	}

	type mypair struct {
		Field1 int32
		Field2 int64
	}
	mp := mypair(p)
	if mp.Field1 != 1 || mp.Field2 != 2 {
		panic(fmt.Sprintf("mp == %#v, want %#v", mp, mypair{1, 2}))
	}
}
