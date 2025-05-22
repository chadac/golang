// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build darwin || (openbsd && !mips64)

package unix

import (
	"syscall"
	_ "unsafe" // for linkname
)

func Unlinkat(dirfd int, path string, flags int) error {
	return unlinkat(dirfd, path, flags)
}

func Openat(dirfd int, path string, flags int, perm uint32) (int, error) {
	return openat(dirfd, path, flags, perm)
}

func Fstatat(dirfd int, path string, stat *syscall.Stat_t, flags int) error {
	return fstatat(dirfd, path, stat, flags)
}

//golang:linkname unlinkat syscall.unlinkat
func unlinkat(dirfd int, path string, flags int) error

//golang:linkname openat syscall.openat
func openat(dirfd int, path string, flags int, perm uint32) (int, error)

//golang:linkname fstatat syscall.fstatat
func fstatat(dirfd int, path string, stat *syscall.Stat_t, flags int) error
