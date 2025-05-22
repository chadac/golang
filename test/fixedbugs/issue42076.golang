// run

// Copyright 2020 The Go Authors. All rights reserved.  Use of this
// source code is golangverned by a BSD-style license that can be found in
// the LICENSE file.

//golang:build cgolang

package main

import (
	"reflect"
	"runtime/cgolang"
)

type NIH struct {
	_ cgolang.Incomplete
}

var x, y NIH

func main() {
	if reflect.DeepEqual(&x, &y) != true {
		panic("should report true")
	}
}
