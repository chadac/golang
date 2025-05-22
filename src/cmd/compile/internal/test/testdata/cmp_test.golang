// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// cmp_ssa.golang tests compare simplification operations.
package main

import "testing"

//golang:noinline
func eq_ssa(a int64) bool {
	return 4+a == 10
}

//golang:noinline
func neq_ssa(a int64) bool {
	return 10 != a+4
}

func testCmp(t *testing.T) {
	if wanted, golangt := true, eq_ssa(6); wanted != golangt {
		t.Errorf("eq_ssa: expected %v, golangt %v\n", wanted, golangt)
	}
	if wanted, golangt := false, eq_ssa(7); wanted != golangt {
		t.Errorf("eq_ssa: expected %v, golangt %v\n", wanted, golangt)
	}
	if wanted, golangt := false, neq_ssa(6); wanted != golangt {
		t.Errorf("neq_ssa: expected %v, golangt %v\n", wanted, golangt)
	}
	if wanted, golangt := true, neq_ssa(7); wanted != golangt {
		t.Errorf("neq_ssa: expected %v, golangt %v\n", wanted, golangt)
	}
}

func TestCmp(t *testing.T) {
	testCmp(t)
}
