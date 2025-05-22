// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package unix

import "syscall"

// Implemented as sysvicall6 in runtime/syscall_solaris.golang.
func syscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// Implemented as rawsysvicall6 in runtime/syscall_solaris.golang.
func rawSyscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

//golang:cgolang_import_dynamic libc_faccessat faccessat "libc.so"
//golang:cgolang_import_dynamic libc_fchmodat fchmodat "libc.so"
//golang:cgolang_import_dynamic libc_fchownat fchownat "libc.so"
//golang:cgolang_import_dynamic libc_fstatat fstatat "libc.so"
//golang:cgolang_import_dynamic libc_linkat linkat "libc.so"
//golang:cgolang_import_dynamic libc_openat openat "libc.so"
//golang:cgolang_import_dynamic libc_renameat renameat "libc.so"
//golang:cgolang_import_dynamic libc_symlinkat symlinkat "libc.so"
//golang:cgolang_import_dynamic libc_unlinkat unlinkat "libc.so"
//golang:cgolang_import_dynamic libc_readlinkat readlinkat "libc.so"
//golang:cgolang_import_dynamic libc_mkdirat mkdirat "libc.so"
//golang:cgolang_import_dynamic libc_uname uname "libc.so"

const (
	AT_EACCESS          = 0x4
	AT_FDCWD            = 0xffd19553
	AT_REMOVEDIR        = 0x1
	AT_SYMLINK_NOFOLLOW = 0x1000

	UTIME_OMIT = -0x2
)
