// Copyright 2010 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build dragolangnfly

package cgolang

import _ "unsafe" // for golang:linkname

// Supply environ and __progname, because we don't
// link against the standard DragolangnFly crt0.o and the
// libc dynamic library needs them.

//golang:linkname _environ environ
//golang:linkname _progname __progname

var _environ uintptr
var _progname uintptr
