// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build openbsd && !mips64

package runtime

import (
	"internal/abi"
	"unsafe"
)

//golang:nosplit
//golang:cgolang_unsafe_args
func thrsleep(ident uintptr, clock_id int32, tsp *timespec, lock uintptr, abort *uint32) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(thrsleep_trampoline)), unsafe.Pointer(&ident))
	KeepAlive(tsp)
	KeepAlive(abort)
	return ret
}
func thrsleep_trampoline()

//golang:nosplit
//golang:cgolang_unsafe_args
func thrwakeup(ident uintptr, n int32) int32 {
	return libcCall(unsafe.Pointer(abi.FuncPCABI0(thrwakeup_trampoline)), unsafe.Pointer(&ident))
}
func thrwakeup_trampoline()

//golang:nosplit
func osyield() {
	libcCall(unsafe.Pointer(abi.FuncPCABI0(sched_yield_trampoline)), unsafe.Pointer(nil))
}
func sched_yield_trampoline()

//golang:nosplit
func osyield_no_g() {
	asmcgolangcall_no_g(unsafe.Pointer(abi.FuncPCABI0(sched_yield_trampoline)), unsafe.Pointer(nil))
}

//golang:cgolang_import_dynamic libc_thrsleep __thrsleep "libc.so"
//golang:cgolang_import_dynamic libc_thrwakeup __thrwakeup "libc.so"
//golang:cgolang_import_dynamic libc_sched_yield sched_yield "libc.so"

//golang:cgolang_import_dynamic _ _ "libc.so"
