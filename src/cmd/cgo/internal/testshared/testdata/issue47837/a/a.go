// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package a

type A interface {
	M()
}

//golang:noinline
func TheFuncWithArgA(a A) {
	a.M()
}

type ImplA struct{}

//golang:noinline
func (A *ImplA) M() {}
