// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// closure.golang tests closure operations.
package main

import "testing"

//golang:noinline
func testCFunc_ssa() int {
	a := 0
	b := func() {
		switch {
		}
		a++
	}
	b()
	b()
	return a
}

func testCFunc(t *testing.T) {
	if want, golangt := 2, testCFunc_ssa(); golangt != want {
		t.Errorf("expected %d, golangt %d", want, golangt)
	}
}

// TestClosure tests closure related behavior.
func TestClosure(t *testing.T) {
	testCFunc(t)
}
