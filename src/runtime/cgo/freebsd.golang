// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build freebsd

package cgolang

import _ "unsafe" // for golang:linkname

// Supply environ and __progname, because we don't
// link against the standard FreeBSD crt0.o and the
// libc dynamic library needs them.

//golang:linkname _environ environ
//golang:linkname _progname __progname

//golang:cgolang_export_dynamic environ
//golang:cgolang_export_dynamic __progname

var _environ uintptr
var _progname uintptr
