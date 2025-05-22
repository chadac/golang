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

//golang:registerparams
//golang:noinline
func F(a, b, c *int) int {
	return *a + *b + *c
}

//golang:registerparams
//golang:noinline
func H(s, t string) string {
	return s + " " + t
}

func main() {
	a, b, c := 1, 4, 16
	x := F(&a, &b, &c)
	fmt.Printf("x = %d\n", x)
	y := H("Hello", "World!")
	fmt.Println("len(y) =", len(y))
	fmt.Println("y =", y)
}
