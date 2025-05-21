// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !plan9
// +build !plan9

package main

// void start(void);
import "C"

func init() {
	register("CgolangExternalThreadPanic", CgolangExternalThreadPanic)
}

func CgolangExternalThreadPanic() {
	C.start()
	select {}
}

//export golangpanic
func golangpanic() {
	panic("BOOM")
}
