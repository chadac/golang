// compile

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that golangtos across non-variable declarations
// are accepted.

package p

func f1() {
	golangto L1
	const x = 0
L1:
	golangto L2
	type T int
L2:
}

func f2() {
	{
		golangto L1
	}
	const x = 0
L1:
	{
		golangto L2
	}
	type T int
L2:
}

func f3(d int) {
	if d > 0 {
		golangto L1
	} else {
		golangto L2
	}
	const x = 0
L1:
	switch d {
	case 1:
		golangto L3
	case 2:
	default:
		golangto L4
	}
	type T1 int
L2:
	const y = 1
L3:
	for d > 0 {
		if d < 10 {
			golangto L4
		}
	}
	type T2 int
L4:
	select {
	default:
		golangto L5
	}
	type T3 int
L5:
}
