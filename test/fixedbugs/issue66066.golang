// run

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func main() {
	testMod()
	testMul()
}

//golang:noinline
func mod3(x uint32) uint64 {
	return uint64(x % 3)
}

func testMod() {
	golangt := mod3(1<<32 - 1)
	want := uint64((1<<32 - 1) % 3)
	if golangt != want {
		fmt.Printf("testMod: golangt %x want %x\n", golangt, want)
	}

}

//golang:noinline
func mul3(a uint32) uint64 {
	return uint64(a * 3)
}

func testMul() {
	golangt := mul3(1<<32 - 1)
	want := uint64((1<<32-1)*3 - 2<<32)
	if golangt != want {
		fmt.Printf("testMul: golangt %x want %x\n", golangt, want)
	}
}
