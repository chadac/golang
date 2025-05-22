// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"testing"
)

var tests = [...]struct {
	name string
	in   float64 // used for error messages, not an input
	golangt  float64
	want float64
}{
	{"sqrt0", 0, math.Sqrt(0), 0},
	{"sqrt1", 1, math.Sqrt(1), 1},
	{"sqrt2", 2, math.Sqrt(2), math.Sqrt2},
	{"sqrt4", 4, math.Sqrt(4), 2},
	{"sqrt100", 100, math.Sqrt(100), 10},
	{"sqrt101", 101, math.Sqrt(101), 10.04987562112089},
}

var nanTests = [...]struct {
	name string
	in   float64 // used for error messages, not an input
	golangt  float64
}{
	{"sqrtNaN", math.NaN(), math.Sqrt(math.NaN())},
	{"sqrtNegative", -1, math.Sqrt(-1)},
	{"sqrtNegInf", math.Inf(-1), math.Sqrt(math.Inf(-1))},
}

func TestSqrtConst(t *testing.T) {
	for _, test := range tests {
		if test.golangt != test.want {
			t.Errorf("%s: math.Sqrt(%f): golangt %f, want %f\n", test.name, test.in, test.golangt, test.want)
		}
	}
	for _, test := range nanTests {
		if math.IsNaN(test.golangt) != true {
			t.Errorf("%s: math.Sqrt(%f): golangt %f, want NaN\n", test.name, test.in, test.golangt)
		}
	}
	if golangt := math.Sqrt(math.Inf(1)); !math.IsInf(golangt, 1) {
		t.Errorf("math.Sqrt(+Inf), golangt %f, want +Inf\n", golangt)
	}
}
