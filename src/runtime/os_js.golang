// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build js && wasm

package runtime

import (
	"unsafe"
)

func exit(code int32)

func write1(fd uintptr, p unsafe.Pointer, n int32) int32 {
	if fd > 2 {
		throw("runtime.write to fd > 2 is unsupported")
	}
	wasmWrite(fd, p, n)
	return n
}

//golang:wasmimport golangjs runtime.wasmWrite
//golang:noescape
func wasmWrite(fd uintptr, p unsafe.Pointer, n int32)

func usleep(usec uint32) {
	// TODO(neelance): implement usleep
}

//golang:wasmimport golangjs runtime.getRandomData
//golang:noescape
func getRandomData(r []byte)

func readRandom(r []byte) int {
	getRandomData(r)
	return len(r)
}

func golangenvs() {
	golangenvs_unix()
}
