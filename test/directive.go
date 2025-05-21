// errorcheck

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that misplaced directives are diagnosed.

//golang:noinline // ERROR "misplaced compiler directive"

//golang:noinline // ERROR "misplaced compiler directive"
package main

//golang:nosplit
func f1() {}

//golang:nosplit
//golang:noinline
func f2() {}

//golang:noinline // ERROR "misplaced compiler directive"

//golang:noinline // ERROR "misplaced compiler directive"
var x int

//golang:noinline // ERROR "misplaced compiler directive"
const c = 1

//golang:noinline // ERROR "misplaced compiler directive"
type T int

type (
	//golang:noinline // ERROR "misplaced compiler directive"
	T2 int
	//golang:noinline // ERROR "misplaced compiler directive"
	T3 int
)

//golang:noinline
func f() {
	x := 1

	{
		_ = x
	}
	//golang:noinline // ERROR "misplaced compiler directive"
	var y int
	_ = y

	//golang:noinline // ERROR "misplaced compiler directive"
	const c = 1

	_ = func() {}

	//golang:noinline // ERROR "misplaced compiler directive"
	type T int
}
