// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || darwin || dragolangnfly || freebsd || openbsd || solaris || zos

package unix

var mapper = &mmapper{
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
}
