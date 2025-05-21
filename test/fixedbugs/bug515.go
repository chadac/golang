// compile

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Caused a golangfrontend crash.

//golang:build gccgolang

package p

//golang:notinheap
type S1 struct{}

type S2 struct {
	r      interface{ Read([]byte) (int, error) }
	s1, s2 []byte
	p      *S1
	n      uintptr
}

var V any = S2{}
