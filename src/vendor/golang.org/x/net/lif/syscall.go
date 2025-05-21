// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build solaris

package lif

import (
	"syscall"
	"unsafe"
)

//golang:cgolang_import_dynamic libc_ioctl ioctl "libc.so"

//golang:linkname procIoctl libc_ioctl

var procIoctl uintptr

func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno)

func ioctl(s, ioc uintptr, arg unsafe.Pointer) error {
	_, _, errno := sysvicall6(uintptr(unsafe.Pointer(&procIoctl)), 3, s, ioc, uintptr(arg), 0, 0, 0)
	if errno != 0 {
		return error(errno)
	}
	return nil
}
