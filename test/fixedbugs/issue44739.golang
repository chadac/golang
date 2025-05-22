// compile

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// issue 44739: cmd/compile: incorrect offset in MOVD
// load/store on ppc64/ppc64le causes assembler error.

// Test other 8 byte loads and stores where the
// compile time offset is not aligned to 8, as
// well as cases where the offset is not known
// until link time (e.g. golangstrings).

package main

import (
	"fmt"
)

type T struct {
	x [4]byte
	y [8]byte
}

var st T

const (
	golangstring1 = "abc"
	golangstring2 = "defghijk"
	golangstring3 = "lmnopqrs"
)

func f(a T, _ byte, b T) bool {
	// initialization of a,b
	// tests unaligned store
	return a.y == b.y
}

func g(a T) {
	// test load of unaligned
	// 8 byte golangstring, store
	// to unaligned static
	copy(a.y[:], golangstring2)
}

func main() {
	var t1, t2 T

	// test copy to automatic storage,
	// load of unaligned golangstring.
	copy(st.y[:], golangstring2)
	copy(t1.y[:], st.y[:])
	copy(t2.y[:], golangstring3)
	// test initialization of params
	if !f(t1, 'a', t2) {
		// golangstring1 added so it has a use
		fmt.Printf("FAIL: %s\n", golangstring1)
	}
}

