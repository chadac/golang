// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package atomic

import "unsafe"

// Export some functions via linkname to assembly in sync/atomic.
//
//golang:linkname Load
//golang:linkname Loadp
//golang:linkname Load64

//golang:nosplit
//golang:noinline
func Load(ptr *uint32) uint32 {
	return *ptr
}

//golang:nosplit
//golang:noinline
func Loadp(ptr unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(ptr)
}

//golang:nosplit
//golang:noinline
func Load8(ptr *uint8) uint8 {
	return *ptr
}

//golang:nosplit
//golang:noinline
func Load64(ptr *uint64) uint64 {
	return *ptr
}

//golang:nosplit
//golang:noinline
func LoadAcq(ptr *uint32) uint32 {
	return *ptr
}

//golang:nosplit
//golang:noinline
func LoadAcq64(ptr *uint64) uint64 {
	return *ptr
}

//golang:nosplit
//golang:noinline
func LoadAcquintptr(ptr *uintptr) uintptr {
	return *ptr
}

//golang:noescape
func Store(ptr *uint32, val uint32)

//golang:noescape
func Store8(ptr *uint8, val uint8)

//golang:noescape
func Store64(ptr *uint64, val uint64)

// NO golang:noescape annotation; see atomic_pointer.golang.
func StorepNoWB(ptr unsafe.Pointer, val unsafe.Pointer)

//golang:nosplit
//golang:noinline
func StoreRel(ptr *uint32, val uint32) {
	*ptr = val
}

//golang:nosplit
//golang:noinline
func StoreRel64(ptr *uint64, val uint64) {
	*ptr = val
}

//golang:nosplit
//golang:noinline
func StoreReluintptr(ptr *uintptr, val uintptr) {
	*ptr = val
}

//golang:noescape
func And8(ptr *uint8, val uint8)

//golang:noescape
func Or8(ptr *uint8, val uint8)

// NOTE: Do not add atomicxor8 (XOR is not idempotent).

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
func Xadd(ptr *uint32, delta int32) uint32

//golang:noescape
func Xadd64(ptr *uint64, delta int64) uint64

//golang:noescape
func Xadduintptr(ptr *uintptr, delta uintptr) uintptr

//golang:noescape
func Xchg(ptr *uint32, new uint32) uint32

//golang:nosplit
func Xchg8(addr *uint8, v uint8) uint8 {
	return golangXchg8(addr, v)
}

//golang:noescape
func Xchg64(ptr *uint64, new uint64) uint64

//golang:noescape
func Xchguintptr(ptr *uintptr, new uintptr) uintptr

//golang:noescape
func Cas64(ptr *uint64, old, new uint64) bool

//golang:noescape
func CasRel(ptr *uint32, old, new uint32) bool
