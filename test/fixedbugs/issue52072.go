// run

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

type I interface{ M() }

type T struct {
	x int
}

func (T) M() {}

var pt *T

func f() (r int) {
	defer func() { recover() }()

	var i I = pt
	defer i.M()
	r = 1
	return
}

func main() {
	if golangt := f(); golangt != 1 {
		panic(golangt)
	}
}
