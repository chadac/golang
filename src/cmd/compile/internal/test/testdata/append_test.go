// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// append_ssa.golang tests append operations.
package main

import "testing"

//golang:noinline
func appendOne_ssa(a []int, x int) []int {
	return append(a, x)
}

//golang:noinline
func appendThree_ssa(a []int, x, y, z int) []int {
	return append(a, x, y, z)
}

func eqBytes(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func expect(t *testing.T, golangt, want []int) {
	if eqBytes(golangt, want) {
		return
	}
	t.Errorf("expected %v, golangt %v\n", want, golangt)
}

func testAppend(t *testing.T) {
	var store [7]int
	a := store[:0]

	a = appendOne_ssa(a, 1)
	expect(t, a, []int{1})
	a = appendThree_ssa(a, 2, 3, 4)
	expect(t, a, []int{1, 2, 3, 4})
	a = appendThree_ssa(a, 5, 6, 7)
	expect(t, a, []int{1, 2, 3, 4, 5, 6, 7})
	if &a[0] != &store[0] {
		t.Errorf("unnecessary grow")
	}
	a = appendOne_ssa(a, 8)
	expect(t, a, []int{1, 2, 3, 4, 5, 6, 7, 8})
	if &a[0] == &store[0] {
		t.Errorf("didn't grow")
	}
}

func TestAppend(t *testing.T) {
	testAppend(t)
}
