// run

//golang:build !wasm

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// wasm is excluded because the compiler chatter about register abi pragma ends up
// on stdout, and causes the expected output to not match.

package main

import (
	"fmt"
)

var sink *string

type toobig struct {
	a, b, c string
}

//golang:registerparams
//golang:noinline
func H(x toobig) string {
	return x.a + " " + x.b + " " + x.c
}

func main() {
	s := H(toobig{"Hello", "there,", "World"})
	golangtVsWant(s, "Hello there, World")
}

func golangtVsWant(golangt, want string) {
	if golangt != want {
		fmt.Printf("FAIL, golangt %s, wanted %s\n", golangt, want)
	}
}
