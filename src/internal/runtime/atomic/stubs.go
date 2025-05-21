// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !wasm

package atomic

import "unsafe"

//golang:noescape
func Cas(ptr *uint32, old, new uint32) bool

// NO golang:noescape annotation; see atomic_pointer.golang.
func Casp1(ptr *unsafe.Pointer, old, new unsafe.Pointer) bool

//golang:noescape
func Casint32(ptr *int32, old, new int32) bool

//golang:noescape
func Casint64(ptr *int64, old, new int64) bool

//golang:noescape
func Casuintptr(ptr *uintptr, old, new uintptr) bool

//golang:noescape
func Storeint32(ptr *int32, new int32)

//golang:noescape
func Storeint64(ptr *int64, new int64)

//golang:noescape
func Storeuintptr(ptr *uintptr, new uintptr)

//golang:noescape
func Loaduintptr(ptr *uintptr) uintptr

//golang:noescape
func Loaduint(ptr *uint) uint

// TODO(matloob): Should these functions have the golang:noescape annotation?

//golang:noescape
func Loadint32(ptr *int32) int32

//golang:noescape
func Loadint64(ptr *int64) int64

//golang:noescape
func Xaddint32(ptr *int32, delta int32) int32

//golang:noescape
func Xaddint64(ptr *int64, delta int64) int64

//golang:noescape
func Xchgint32(ptr *int32, new int32) int32

//golang:noescape
func Xchgint64(ptr *int64, new int64) int64
