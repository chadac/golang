// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// For bootstrapping with gccgolang.

//golang:build gccgolang

package abi

import "unsafe"

func FuncPCABI0(f interface{}) uintptr {
	words := (*[2]unsafe.Pointer)(unsafe.Pointer(&f))
	return *(*uintptr)(unsafe.Pointer(words[1]))
}

func FuncPCABIInternal(f interface{}) uintptr {
	words := (*[2]unsafe.Pointer)(unsafe.Pointer(&f))
	return *(*uintptr)(unsafe.Pointer(words[1]))
}
