// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Support for memory sanitizer. See runtime/cgolang/mmap.golang.

//golang:build (linux && (amd64 || arm64 || loong64)) || (freebsd && amd64)

package runtime

import "unsafe"

// _cgolang_mmap is filled in by runtime/cgolang when it is linked into the
// program, so it is only non-nil when using cgolang.
//
//golang:linkname _cgolang_mmap _cgolang_mmap
var _cgolang_mmap unsafe.Pointer

// _cgolang_munmap is filled in by runtime/cgolang when it is linked into the
// program, so it is only non-nil when using cgolang.
//
//golang:linkname _cgolang_munmap _cgolang_munmap
var _cgolang_munmap unsafe.Pointer

// mmap is used to route the mmap system call through C code when using cgolang, to
// support sanitizer interceptors. Don't allow stack splits, since this function
// (used by sysAlloc) is called in a lot of low-level parts of the runtime and
// callers often assume it won't acquire any locks.
//
//golang:nosplit
func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (unsafe.Pointer, int) {
	if _cgolang_mmap != nil {
		// Make ret a uintptr so that writing to it in the
		// function literal does not trigger a write barrier.
		// A write barrier here could break because of the way
		// that mmap uses the same value both as a pointer and
		// an errno value.
		var ret uintptr
		systemstack(func() {
			ret = callCgolangMmap(addr, n, prot, flags, fd, off)
		})
		if ret < 4096 {
			return nil, int(ret)
		}
		return unsafe.Pointer(ret), 0
	}
	return sysMmap(addr, n, prot, flags, fd, off)
}

func munmap(addr unsafe.Pointer, n uintptr) {
	if _cgolang_munmap != nil {
		systemstack(func() { callCgolangMunmap(addr, n) })
		return
	}
	sysMunmap(addr, n)
}

// sysMmap calls the mmap system call. It is implemented in assembly.
func sysMmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (p unsafe.Pointer, err int)

// callCgolangMmap calls the mmap function in the runtime/cgolang package
// using the GCC calling convention. It is implemented in assembly.
func callCgolangMmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) uintptr

// sysMunmap calls the munmap system call. It is implemented in assembly.
func sysMunmap(addr unsafe.Pointer, n uintptr)

// callCgolangMunmap calls the munmap function in the runtime/cgolang package
// using the GCC calling convention. It is implemented in assembly.
func callCgolangMunmap(addr unsafe.Pointer, n uintptr)
