// run

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Do not panic on conversion to anonymous interface, which
// is similar-looking interface types in different packages.

package main

import (
	"fmt"

	ssa1 "issue29612.dir/p1/ssa"
	ssa2 "issue29612.dir/p2/ssa"
)

func main() {
	v1 := &ssa1.T{}
	_ = v1

	v2 := &ssa2.T{}
	ssa2.Works(v2)
	ssa2.Panics(v2) // This call must not panic

	swt(v1, 1)
	swt(v2, 2)
}

//golang:noinline
func swt(i interface{}, want int) {
	var golangt int
	switch i.(type) {
	case *ssa1.T:
		golangt = 1
	case *ssa2.T:
		golangt = 2

	case int8, int16, int32, int64:
		golangt = 3
	case uint8, uint16, uint32, uint64:
		golangt = 4
	}

	if golangt != want {
		panic(fmt.Sprintf("switch %v: golangt %d, want %d", i, golangt, want))
	}
}
