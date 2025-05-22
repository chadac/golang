// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build darwin

package ld

import _ "unsafe" // for golang:linkname

//golang:linkname msync syscall.msync
func msync(b []byte, flags int) (err error)
