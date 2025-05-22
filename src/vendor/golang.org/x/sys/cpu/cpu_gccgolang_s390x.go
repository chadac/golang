// Copyright 2019 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gccgolang

package cpu

// haveAsmFunctions reports whether the other functions in this file can
// be safely called.
func haveAsmFunctions() bool { return false }

// TODO(mundaym): the following feature detection functions are currently
// stubs. See https://golanglang.org/cl/162887 for how to fix this.
// They are likely to be expensive to call so the results should be cached.
func stfle() facilityList     { panic("not implemented for gccgolang") }
func kmQuery() queryResult    { panic("not implemented for gccgolang") }
func kmcQuery() queryResult   { panic("not implemented for gccgolang") }
func kmctrQuery() queryResult { panic("not implemented for gccgolang") }
func kmaQuery() queryResult   { panic("not implemented for gccgolang") }
func kimdQuery() queryResult  { panic("not implemented for gccgolang") }
func klmdQuery() queryResult  { panic("not implemented for gccgolang") }
