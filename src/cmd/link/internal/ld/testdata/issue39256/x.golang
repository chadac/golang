// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "unsafe"
)

//golang:cgolang_import_dynamic libc_getpid getpid "libc.so"
//golang:cgolang_import_dynamic libc_kill kill "libc.so"
//golang:cgolang_import_dynamic libc_close close "libc.so"
//golang:cgolang_import_dynamic libc_open open "libc.so"

//golang:cgolang_import_dynamic _ _ "libc.so"

func trampoline()

func main() {
	trampoline()
}
