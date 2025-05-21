// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || linux

// Supporting definitions for os_uname.golang on AIX and Linux.

package osinfo

import "syscall"

type utsname = syscall.Utsname

func uname(buf *utsname) error {
	return syscall.Uname(buf)
}
