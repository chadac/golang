// errorcheck

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that misplaced directives are diagnosed.

// ok
//golang:build !ignore

package main

//golang:build bad // ERROR "misplaced compiler directive"

//golang:noinline // ERROR "misplaced compiler directive"
type (
	T2  int //golang:noinline // ERROR "misplaced compiler directive"
	T2b int
	T2c int
	T3  int
)

//golang:noinline // ERROR "misplaced compiler directive"
type (
	T4 int
)

//golang:noinline // ERROR "misplaced compiler directive"
type ()

type T5 int

func g() {} //golang:noinline // ERROR "misplaced compiler directive"

// ok: attached to f (duplicated yes, but ok)
//golang:noinline

//golang:noinline
func f() {
	//golang:noinline // ERROR "misplaced compiler directive"
	x := 1

	//golang:noinline // ERROR "misplaced compiler directive"
	{
		_ = x //golang:noinline // ERROR "misplaced compiler directive"
	}
	var y int //golang:noinline // ERROR "misplaced compiler directive"
	//golang:noinline // ERROR "misplaced compiler directive"
	_ = y

	const c = 1

	_ = func() {}
}

// EOF
//golang:noinline // ERROR "misplaced compiler directive"
