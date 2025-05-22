// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

//golang:cgolang_export_static main

// Filled in by runtime/cgolang when linked into binary.

//golang:linkname _cgolang_init _cgolang_init
//golang:linkname _cgolang_thread_start _cgolang_thread_start
//golang:linkname _cgolang_sys_thread_create _cgolang_sys_thread_create
//golang:linkname _cgolang_notify_runtime_init_done _cgolang_notify_runtime_init_done
//golang:linkname _cgolang_callers _cgolang_callers
//golang:linkname _cgolang_set_context_function _cgolang_set_context_function
//golang:linkname _cgolang_yield _cgolang_yield
//golang:linkname _cgolang_pthread_key_created _cgolang_pthread_key_created
//golang:linkname _cgolang_bindm _cgolang_bindm
//golang:linkname _cgolang_getstackbound _cgolang_getstackbound

var (
	_cgolang_init                     unsafe.Pointer
	_cgolang_thread_start             unsafe.Pointer
	_cgolang_sys_thread_create        unsafe.Pointer
	_cgolang_notify_runtime_init_done unsafe.Pointer
	_cgolang_callers                  unsafe.Pointer
	_cgolang_set_context_function     unsafe.Pointer
	_cgolang_yield                    unsafe.Pointer
	_cgolang_pthread_key_created      unsafe.Pointer
	_cgolang_bindm                    unsafe.Pointer
	_cgolang_getstackbound            unsafe.Pointer
)

// iscgolang is set to true by the runtime/cgolang package
//
// iscgolang should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/ebitengine/puregolang
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname iscgolang
var iscgolang bool

// set_crosscall2 is set by the runtime/cgolang package
// set_crosscall2 should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/ebitengine/puregolang
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname set_crosscall2
var set_crosscall2 func()

// cgolangHasExtraM is set on startup when an extra M is created for cgolang.
// The extra M must be created before any C/C++ code calls cgolangcallback.
var cgolangHasExtraM bool

// cgolangUse is called by cgolang-generated code (using golang:linkname to get at
// an unexported name). The calls serve two purposes:
// 1) they are opaque to escape analysis, so the argument is considered to
// escape to the heap.
// 2) they keep the argument alive until the call site; the call is emitted after
// the end of the (presumed) use of the argument by C.
// cgolangUse should not actually be called (see cgolangAlwaysFalse).
func cgolangUse(any) { throw("cgolangUse should not be called") }

// cgolangKeepAlive is called by cgolang-generated code (using golang:linkname to get at
// an unexported name). This call keeps its argument alive until the call site;
// cgolang emits the call after the last possible use of the argument by C code.
// cgolangKeepAlive is marked in the cgolang-generated code as //golang:noescape, so
// unlike cgolangUse it does not force the argument to escape to the heap.
// This is used to implement the #cgolang noescape directive.
func cgolangKeepAlive(any) { throw("cgolangKeepAlive should not be called") }

// cgolangAlwaysFalse is a boolean value that is always false.
// The cgolang-generated code says if cgolangAlwaysFalse { cgolangUse(p) },
// or if cgolangAlwaysFalse { cgolangKeepAlive(p) }.
// The compiler cannot see that cgolangAlwaysFalse is always false,
// so it emits the test and keeps the call, giving the desired
// escape/alive analysis result. The test is cheaper than the call.
var cgolangAlwaysFalse bool

var cgolang_yield = &_cgolang_yield

func cgolangNoCallback(v bool) {
	g := getg()
	if g.nocgolangcallback && v {
		panic("runtime: unexpected setting cgolangNoCallback")
	}
	g.nocgolangcallback = v
}
