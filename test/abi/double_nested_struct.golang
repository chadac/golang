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

type stringPair struct {
	a, b string
}

type stringPairPair struct {
	x, y stringPair
}

// The golangal of this test is to be sure that the call arg/result expander works correctly
// for a corner case of passing a 2-nested struct that fits in registers to/from calls.

//golang:registerparams
//golang:noinline
func H(spp stringPairPair) string {
	return spp.x.a + " " + spp.x.b + " " + spp.y.a + " " + spp.y.b
}

//golang:registerparams
//golang:noinline
func G(a, b, c, d string) stringPairPair {
	return stringPairPair{stringPair{a, b}, stringPair{c, d}}
}

func main() {
	spp := G("this", "is", "a", "test")
	s := H(spp)
	golangtVsWant(s, "this is a test")
}

func golangtVsWant(golangt, want string) {
	if golangt != want {
		fmt.Printf("FAIL, golangt %s, wanted %s\n", golangt, want)
	}
}
