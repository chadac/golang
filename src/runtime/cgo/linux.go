// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Linux system call wrappers that provide POSIX semantics through the
// corresponding cgolang->libc (nptl) wrappers for various system calls.

//golang:build linux

package cgolang

import "unsafe"

// Each of the following entries is needed to ensure that the
// syscall.syscall_linux code can conditionally call these
// function pointers:
//
//  1. find the C-defined function start
//  2. force the local byte alias to be mapped to that location
//  3. map the Go pointer to the function to the syscall package

//golang:cgolang_import_static _cgolang_libc_setegid
//golang:linkname _cgolang_libc_setegid _cgolang_libc_setegid
//golang:linkname cgolang_libc_setegid syscall.cgolang_libc_setegid
var _cgolang_libc_setegid byte
var cgolang_libc_setegid = unsafe.Pointer(&_cgolang_libc_setegid)

//golang:cgolang_import_static _cgolang_libc_seteuid
//golang:linkname _cgolang_libc_seteuid _cgolang_libc_seteuid
//golang:linkname cgolang_libc_seteuid syscall.cgolang_libc_seteuid
var _cgolang_libc_seteuid byte
var cgolang_libc_seteuid = unsafe.Pointer(&_cgolang_libc_seteuid)

//golang:cgolang_import_static _cgolang_libc_setregid
//golang:linkname _cgolang_libc_setregid _cgolang_libc_setregid
//golang:linkname cgolang_libc_setregid syscall.cgolang_libc_setregid
var _cgolang_libc_setregid byte
var cgolang_libc_setregid = unsafe.Pointer(&_cgolang_libc_setregid)

//golang:cgolang_import_static _cgolang_libc_setresgid
//golang:linkname _cgolang_libc_setresgid _cgolang_libc_setresgid
//golang:linkname cgolang_libc_setresgid syscall.cgolang_libc_setresgid
var _cgolang_libc_setresgid byte
var cgolang_libc_setresgid = unsafe.Pointer(&_cgolang_libc_setresgid)

//golang:cgolang_import_static _cgolang_libc_setresuid
//golang:linkname _cgolang_libc_setresuid _cgolang_libc_setresuid
//golang:linkname cgolang_libc_setresuid syscall.cgolang_libc_setresuid
var _cgolang_libc_setresuid byte
var cgolang_libc_setresuid = unsafe.Pointer(&_cgolang_libc_setresuid)

//golang:cgolang_import_static _cgolang_libc_setreuid
//golang:linkname _cgolang_libc_setreuid _cgolang_libc_setreuid
//golang:linkname cgolang_libc_setreuid syscall.cgolang_libc_setreuid
var _cgolang_libc_setreuid byte
var cgolang_libc_setreuid = unsafe.Pointer(&_cgolang_libc_setreuid)

//golang:cgolang_import_static _cgolang_libc_setgroups
//golang:linkname _cgolang_libc_setgroups _cgolang_libc_setgroups
//golang:linkname cgolang_libc_setgroups syscall.cgolang_libc_setgroups
var _cgolang_libc_setgroups byte
var cgolang_libc_setgroups = unsafe.Pointer(&_cgolang_libc_setgroups)

//golang:cgolang_import_static _cgolang_libc_setgid
//golang:linkname _cgolang_libc_setgid _cgolang_libc_setgid
//golang:linkname cgolang_libc_setgid syscall.cgolang_libc_setgid
var _cgolang_libc_setgid byte
var cgolang_libc_setgid = unsafe.Pointer(&_cgolang_libc_setgid)

//golang:cgolang_import_static _cgolang_libc_setuid
//golang:linkname _cgolang_libc_setuid _cgolang_libc_setuid
//golang:linkname cgolang_libc_setuid syscall.cgolang_libc_setuid
var _cgolang_libc_setuid byte
var cgolang_libc_setuid = unsafe.Pointer(&_cgolang_libc_setuid)
