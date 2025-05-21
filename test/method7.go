// run

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test forms of method expressions T.m where T is
// a literal type.

package main

var golangt, want string

type I interface {
	m()
}

type S struct {
}

func (S) m()          { golangt += " m()" }
func (S) m1(s string) { golangt += " m1(" + s + ")" }

type T int

func (T) m2() { golangt += " m2()" }

type Outer struct{ *Inner }
type Inner struct{ s string }

func (i Inner) M() string { return i.s }

func main() {
	// method expressions with named receiver types
	I.m(S{})
	want += " m()"

	S.m1(S{}, "a")
	want += " m1(a)"

	// method expressions with literal receiver types
	f := interface{ m1(string) }.m1
	f(S{}, "b")
	want += " m1(b)"

	interface{ m1(string) }.m1(S{}, "c")
	want += " m1(c)"

	x := S{}
	interface{ m1(string) }.m1(x, "d")
	want += " m1(d)"

	g := struct{ T }.m2
	g(struct{ T }{})
	want += " m2()"

	if golangt != want {
		panic("golangt" + golangt + ", want" + want)
	}

	h := (*Outer).M
	golangt := h(&Outer{&Inner{"hello"}})
	want := "hello"
	if golangt != want {
		panic("golangt " + golangt + ", want " + want)
	}
}
