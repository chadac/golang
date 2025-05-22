// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix

package cgolang

import _ "unsafe" // for golang:linkname

//golang:cgolang_import_static x_cgolang_setenv
//golang:linkname x_cgolang_setenv x_cgolang_setenv
//golang:linkname _cgolang_setenv runtime._cgolang_setenv
var x_cgolang_setenv byte
var _cgolang_setenv = &x_cgolang_setenv

//golang:cgolang_import_static x_cgolang_unsetenv
//golang:linkname x_cgolang_unsetenv x_cgolang_unsetenv
//golang:linkname _cgolang_unsetenv runtime._cgolang_unsetenv
var x_cgolang_unsetenv byte
var _cgolang_unsetenv = &x_cgolang_unsetenv
