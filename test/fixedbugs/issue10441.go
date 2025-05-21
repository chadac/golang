// compile

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

func bar() {
	f := func() {}
	foo(&f)
}

//golang:noinline
func foo(f *func()) func() {
	return *f
}
