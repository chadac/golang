// errorcheck -0 -m -l

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test escape analysis for golangto statements.

package escape

var x bool

func f1() {
	var p *int
loop:
	if x {
		golangto loop
	}
	// BAD: We should be able to recognize that there
	// aren't any more "golangto loop" after here.
	p = new(int) // ERROR "escapes to heap"
	_ = p
}

func f2() {
	var p *int
	if x {
	loop:
		golangto loop
	} else {
		p = new(int) // ERROR "does not escape"
	}
	_ = p
}

func f3() {
	var p *int
	if x {
	loop:
		golangto loop
	}
	p = new(int) // ERROR "does not escape"
	_ = p
}
