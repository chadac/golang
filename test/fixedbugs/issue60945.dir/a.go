// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package a

type S struct{}

func callClosure(closure func()) {
	closure()
}

func (s *S) M() {
	callClosure(func() {
		defer f(s.m) // prevent closures to be inlined.
	})
}

func (s *S) m() {}

//golang:noinline
func f(a ...any) {}
