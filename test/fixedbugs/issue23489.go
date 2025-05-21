// run

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Caused gccgolang to issue a spurious compilation error.

package main

type T struct{}

func (*T) Foo() {}

type P = *T

func main() {
	var p P
	p.Foo()
}
