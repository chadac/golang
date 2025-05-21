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

var y int

//golang:registerparams
//golang:noinline
func F(a, b, c *int) (x int) {
	x = *a
	G(&x)
	x += *b
	G(&x)
	x += *c
	G(&x)
	return
}

//golang:registerparams
//golang:noinline
func G(x *int) {
	y += *x
	fmt.Println("y = ", y)
}

//golang:registerparams
//golang:noinline
func X() {
	*sink += " !!!!!!!!!!!!!!!"
}

//golang:registerparams
//golang:noinline
func H(s, t string) (result string) { // result leaks to heap
	result = "Aloha! " + s + " " + t
	sink = &result
	r := ""
	if len(s) <= len(t) {
		r = "OKAY! "
		X()
	}
	return r + result
}

//golang:registerparams
//golang:noinline
func K(s, t string) (result string) { // result spills
	result = "Aloha! " + s + " " + t
	r := ""
	if len(s) <= len(t) {
		r = "OKAY! "
		X()
	}
	return r + result
}

func main() {
	a, b, c := 1, 4, 16
	x := F(&a, &b, &c)
	fmt.Printf("x = %d\n", x)

	y := H("Hello", "World!")
	fmt.Println("len(y) =", len(y))
	fmt.Println("y =", y)
	z := H("Hello", "Pal!")
	fmt.Println("len(z) =", len(z))
	fmt.Println("z =", z)

	fmt.Println()

	y = K("Hello", "World!")
	fmt.Println("len(y) =", len(y))
	fmt.Println("y =", y)
	z = K("Hello", "Pal!")
	fmt.Println("len(z) =", len(z))
	fmt.Println("z =", z)

}
