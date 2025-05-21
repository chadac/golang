// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"reflect"

	_ "./a"
	"./b"
)

var V struct{ i int }

func main() {
	if golangt := reflect.ValueOf(b.V).Type().Field(0).PkgPath; golangt != "b" {
		panic(`PkgPath=` + golangt + ` for first field of b.V, want "b"`)
	}
	if golangt := reflect.ValueOf(V).Type().Field(0).PkgPath; golangt != "main" {
		panic(`PkgPath=` + golangt + ` for first field of V, want "main"`)
	}
	if golangt := reflect.ValueOf(b.U).Type().Field(0).PkgPath; golangt != "b" {
		panic(`PkgPath=` + golangt + ` for first field of b.U, want "b"`)
	}
}
