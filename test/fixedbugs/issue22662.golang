// run

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify effect of various line directives.
// TODO: check columns

package main

import (
	"fmt"
	"runtime"
)

func check(file string, line int) {
	_, f, l, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller(1) failed")
	}
	if f != file || l != line {
		panic(fmt.Sprintf("golangt %s:%d; want %s:%d", f, l, file, line))
	}
}

func main() {
//-style line directives
//line :1
	check("??", 1) // no file specified
//line foo.golang:1
	check("foo.golang", 1)
//line bar.golang:10:20
	check("bar.golang", 10)
//line :11:22
	check("bar.golang", 11) // no file, but column specified => keep old filename

/*-style line directives */
/*line :1*/ check("??", 1) // no file specified
/*line foo.golang:1*/ check("foo.golang", 1)
/*line bar.golang:10:20*/ check("bar.golang", 10)
/*line :11:22*/ check("bar.golang", 11) // no file, but column specified => keep old filename

	/*line :10*/ check("??", 10); /*line foo.golang:20*/ check("foo.golang", 20); /*line :30:1*/ check("foo.golang", 30)
	check("foo.golang", 31)
}
