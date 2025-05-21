// errorcheck

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Tests that errors are reported for misuse of linkname.
package p

import _ "unsafe"

type t int

var x, y int

func F[T any](T) {}

//golang:linkname x ok

// ERROR "//golang:linkname must refer to declared function or variable"
// ERROR "//golang:linkname must refer to declared function or variable"
// ERROR "duplicate //golang:linkname for x"
// ERROR "//golang:linkname reference of an instantiation is not allowed"

//line linkname3.golang:20
//golang:linkname nonexist nonexist
//golang:linkname t notvarfunc
//golang:linkname x duplicate
//golang:linkname i F[golang.shape.int]
