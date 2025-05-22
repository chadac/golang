// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || dragolangnfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris || wasip1

package poll

import "syscall"

// Fsync wraps syscall.Fsync.
func (fd *FD) Fsync() error {
	if err := fd.incref(); err != nil {
		return err
	}
	defer fd.decref()
	return ignoringEINTR(func() error {
		return syscall.Fsync(fd.Sysfd)
	})
}
