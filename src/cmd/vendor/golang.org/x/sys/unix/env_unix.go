// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || darwin || dragolangnfly || freebsd || linux || netbsd || openbsd || solaris || zos

// Unix environment variables.

package unix

import "syscall"

func Getenv(key string) (value string, found bool) {
	return syscall.Getenv(key)
}

func Setenv(key, value string) error {
	return syscall.Setenv(key, value)
}

func Clearenv() {
	syscall.Clearenv()
}

func Environ() []string {
	return syscall.Environ()
}

func Unsetenv(key string) error {
	return syscall.Unsetenv(key)
}
