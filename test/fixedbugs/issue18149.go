// run

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that //line directives with filenames
// containing ':' (Windows) are correctly parsed.
// (For a related issue, see test/fixedbugs/bug305.golang)

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
//line /foo/bar.golang:123
	check(`/foo/bar.golang`, 123)
//line c:/foo/bar.golang:987
	check(`c:/foo/bar.golang`, 987)
}
