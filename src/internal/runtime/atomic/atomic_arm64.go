// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build arm64

package atomic

import (
	"internal/cpu"
	"unsafe"
)

const (
	offsetARM64HasATOMICS = unsafe.Offsetof(cpu.ARM64.HasATOMICS)
)

//golang:noescape
func Xadd(ptr *uint32, delta int32) uint32

//golang:noescape
func Xadd64(ptr *uint64, delta int64) uint64

//golang:noescape
func Xadduintptr(ptr *uintptr, delta uintptr) uintptr

//golang:noescape
func Xchg8(ptr *uint8, new uint8) uint8

//golang:noescape
func Xchg(ptr *uint32, new uint32) uint32

//golang:noescape
func Xchg64(ptr *uint64, new uint64) uint64

//golang:noescape
func Xchguintptr(ptr *uintptr, new uintptr) uintptr

//golang:noescape
func Load(ptr *uint32) uint32

//golang:noescape
func Load8(ptr *uint8) uint8

//golang:noescape
func Load64(ptr *uint64) uint64

// NO golang:noescape annotation; *ptr escapes if result escapes (#31525)
func Loadp(ptr unsafe.Pointer) unsafe.Pointer

//golang:noescape
func LoadAcq(addr *uint32) uint32

//golang:noescape
func LoadAcq64(ptr *uint64) uint64

//golang:noescape
func LoadAcquintptr(ptr *uintptr) uintptr

//golang:noescape
func Or8(ptr *uint8, val uint8)

//golang:noescape
func And8(ptr *uint8, val uint8)

//golang:noescape
func And(ptr *uint32, val uint32)

//golang:noescape
func Or(ptr *uint32, val uint32)

//golang:noescape
func And32(ptr *uint32, val uint32) uint32

//golang:noescape
func Or32(ptr *uint32, val uint32) uint32

//golang:noescape
func And64(ptr *uint64, val uint64) uint64

//golang:noescape
func Or64(ptr *uint64, val uint64) uint64

//golang:noescape
func Anduintptr(ptr *uintptr, val uintptr) uintptr

//golang:noescape
func Oruintptr(ptr *uintptr, val uintptr) uintptr

//golang:noescape
func Cas64(ptr *uint64, old, new uint64) bool

//golang:noescape
func CasRel(ptr *uint32, old, new uint32) bool

//golang:noescape
func Store(ptr *uint32, val uint32)

//golang:noescape
func Store8(ptr *uint8, val uint8)

//golang:noescape
func Store64(ptr *uint64, val uint64)

// NO golang:noescape annotation; see atomic_pointer.golang.
func StorepNoWB(ptr unsafe.Pointer, val unsafe.Pointer)

//golang:noescape
func StoreRel(ptr *uint32, val uint32)

//golang:noescape
func StoreRel64(ptr *uint64, val uint64)

//golang:noescape
func StoreReluintptr(ptr *uintptr, val uintptr)
