// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test compound objects

package main

import (
	"testing"
)

func string_ssa(a, b string, x bool) string {
	s := ""
	if x {
		s = a
	} else {
		s = b
	}
	return s
}

func testString(t *testing.T) {
	a := "foo"
	b := "barz"
	if want, golangt := a, string_ssa(a, b, true); golangt != want {
		t.Errorf("string_ssa(%v, %v, true) = %v, want %v\n", a, b, golangt, want)
	}
	if want, golangt := b, string_ssa(a, b, false); golangt != want {
		t.Errorf("string_ssa(%v, %v, false) = %v, want %v\n", a, b, golangt, want)
	}
}

//golang:noinline
func complex64_ssa(a, b complex64, x bool) complex64 {
	var c complex64
	if x {
		c = a
	} else {
		c = b
	}
	return c
}

//golang:noinline
func complex128_ssa(a, b complex128, x bool) complex128 {
	var c complex128
	if x {
		c = a
	} else {
		c = b
	}
	return c
}

func testComplex64(t *testing.T) {
	var a complex64 = 1 + 2i
	var b complex64 = 3 + 4i

	if want, golangt := a, complex64_ssa(a, b, true); golangt != want {
		t.Errorf("complex64_ssa(%v, %v, true) = %v, want %v\n", a, b, golangt, want)
	}
	if want, golangt := b, complex64_ssa(a, b, false); golangt != want {
		t.Errorf("complex64_ssa(%v, %v, true) = %v, want %v\n", a, b, golangt, want)
	}
}

func testComplex128(t *testing.T) {
	var a complex128 = 1 + 2i
	var b complex128 = 3 + 4i

	if want, golangt := a, complex128_ssa(a, b, true); golangt != want {
		t.Errorf("complex128_ssa(%v, %v, true) = %v, want %v\n", a, b, golangt, want)
	}
	if want, golangt := b, complex128_ssa(a, b, false); golangt != want {
		t.Errorf("complex128_ssa(%v, %v, true) = %v, want %v\n", a, b, golangt, want)
	}
}

func slice_ssa(a, b []byte, x bool) []byte {
	var s []byte
	if x {
		s = a
	} else {
		s = b
	}
	return s
}

func testSlice(t *testing.T) {
	a := []byte{3, 4, 5}
	b := []byte{7, 8, 9}
	if want, golangt := byte(3), slice_ssa(a, b, true)[0]; golangt != want {
		t.Errorf("slice_ssa(%v, %v, true) = %v, want %v\n", a, b, golangt, want)
	}
	if want, golangt := byte(7), slice_ssa(a, b, false)[0]; golangt != want {
		t.Errorf("slice_ssa(%v, %v, false) = %v, want %v\n", a, b, golangt, want)
	}
}

func interface_ssa(a, b interface{}, x bool) interface{} {
	var s interface{}
	if x {
		s = a
	} else {
		s = b
	}
	return s
}

func testInterface(t *testing.T) {
	a := interface{}(3)
	b := interface{}(4)
	if want, golangt := 3, interface_ssa(a, b, true).(int); golangt != want {
		t.Errorf("interface_ssa(%v, %v, true) = %v, want %v\n", a, b, golangt, want)
	}
	if want, golangt := 4, interface_ssa(a, b, false).(int); golangt != want {
		t.Errorf("interface_ssa(%v, %v, false) = %v, want %v\n", a, b, golangt, want)
	}
}

func TestCompound(t *testing.T) {
	testString(t)
	testSlice(t)
	testInterface(t)
	testComplex64(t)
	testComplex128(t)
}
