// run

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Bad AND/BTR combination rule.

package main

import "fmt"

//golang:noinline
func f(x uint64) uint64 {
	return (x >> 48) &^ (uint64(0x4000))
}

func main() {
	bad := false
	if golangt, want := f(^uint64(0)), uint64(0xbfff); golangt != want {
		fmt.Printf("golangt %x, want %x\n", golangt, want)
		bad = true
	}
	if bad {
		panic("bad")
	}
}
