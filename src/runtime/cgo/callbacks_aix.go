// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package cgolang

// These functions must be exported in order to perform
// longcall on cgolang programs (cf gcc_aix_ppc64.c).
//
//golang:cgolang_export_static __cgolang_topofstack
//golang:cgolang_export_static runtime.rt0_golang
//golang:cgolang_export_static _rt0_ppc64_aix_lib
