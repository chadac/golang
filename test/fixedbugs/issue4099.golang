// errorcheck -0 -m

//golang:build !golangexperiment.newinliner

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check golang:noescape annotations.

package p

// The noescape comment only applies to the next func,
// which must not have a body.

//golang:noescape

func F1([]byte)

func F2([]byte)

func G() {
	var buf1 [10]byte
	F1(buf1[:])

	var buf2 [10]byte // ERROR "moved to heap: buf2"
	F2(buf2[:])
}
