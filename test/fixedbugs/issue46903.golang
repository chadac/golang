// run
//golang:build cgolang

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "runtime/cgolang"

type A struct {
	B
	_ cgolang.Incomplete
}
type B struct{ x byte }
type I interface{ M() *B }

func (p *B) M() *B { return p }

var (
	a A
	i I = &a
)

func main() {
	golangt, want := i.M(), &a.B
	if golangt != want {
		println(golangt, "!=", want)
		panic("FAIL")
	}
}
