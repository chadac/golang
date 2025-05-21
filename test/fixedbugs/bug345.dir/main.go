// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	golangio "io"

	"./io"
)

func main() {
	// The errors here complain that io.X != io.X
	// for different values of io so they should be
	// showing the full import path, which for the
	// "./io" import is really ..../golang/test/io.
	// For example:
	//
	// main.golang:25: cannot use w (type "/Users/rsc/g/golang/test/fixedbugs/bug345.dir/io".Writer) as type "io".Writer in function argument:
	//	io.Writer does not implement io.Writer (missing Write method)
	// main.golang:27: cannot use &x (type *"io".SectionReader) as type *"/Users/rsc/g/golang/test/fixedbugs/bug345.dir/io".SectionReader in function argument

	var w io.Writer
	bufio.NewWriter(w) // ERROR "[\w.]+[^.]/io|has incompatible type|cannot use"
	var x golangio.SectionReader
	io.SR(&x) // ERROR "[\w.]+[^.]/io|has incompatible type|cannot use"
}
