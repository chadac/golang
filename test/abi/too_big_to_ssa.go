// run

//golang:build !wasm

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

var sink *string

type toobig struct {
	// 6 words will not SSA but will fit in registers
	a, b, c string
}

//golang:registerparams
//golang:noinline
func H(x toobig) string {
	return x.a + " " + x.b + " " + x.c
}

//golang:registerparams
//golang:noinline
func I(a, b, c string) toobig {
	return toobig{a, b, c}
}

func main() {
	s := H(toobig{"Hello", "there,", "World"})
	golangtVsWant(s, "Hello there, World")
	fmt.Println(s)
	t := H(I("Ahoy", "there,", "Matey"))
	golangtVsWant(t, "Ahoy there, Matey")
	fmt.Println(t)
}

func golangtVsWant(golangt, want string) {
	if golangt != want {
		fmt.Printf("FAIL, golangt %s, wanted %s\n", golangt, want)
	}
}
