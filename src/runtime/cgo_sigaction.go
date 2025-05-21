// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Support for sanitizers. See runtime/cgolang/sigaction.golang.

//golang:build (linux && (amd64 || arm64 || loong64 || ppc64le)) || (freebsd && amd64)

package runtime

import "unsafe"

// _cgolang_sigaction is filled in by runtime/cgolang when it is linked into the
// program, so it is only non-nil when using cgolang.
//
//golang:linkname _cgolang_sigaction _cgolang_sigaction
var _cgolang_sigaction unsafe.Pointer

//golang:nosplit
//golang:nowritebarrierrec
func sigaction(sig uint32, new, old *sigactiont) {
	// racewalk.golang avoids adding sanitizing instrumentation to package runtime,
	// but we might be calling into instrumented C functions here,
	// so we need the pointer parameters to be properly marked.
	//
	// Mark the input as having been written before the call
	// and the output as read after.
	if msanenabled && new != nil {
		msanwrite(unsafe.Pointer(new), unsafe.Sizeof(*new))
	}
	if asanenabled && new != nil {
		asanwrite(unsafe.Pointer(new), unsafe.Sizeof(*new))
	}
	if _cgolang_sigaction == nil || inForkedChild {
		sysSigaction(sig, new, old)
	} else {
		// We need to call _cgolang_sigaction, which means we need a big enough stack
		// for C.  To complicate matters, we may be in libpreinit (before the
		// runtime has been initialized) or in an asynchronous signal handler (with
		// the current thread in transition between golangroutines, or with the g0
		// system stack already in use).

		var ret int32

		var g *g
		if mainStarted {
			g = getg()
		}
		sp := uintptr(unsafe.Pointer(&sig))
		switch {
		case g == nil:
			// No g: we're on a C stack or a signal stack.
			ret = callCgolangSigaction(uintptr(sig), new, old)
		case sp < g.stack.lo || sp >= g.stack.hi:
			// We're no longer on g's stack, so we must be handling a signal.  It's
			// possible that we interrupted the thread during a transition between g
			// and g0, so we should stay on the current stack to avoid corrupting g0.
			ret = callCgolangSigaction(uintptr(sig), new, old)
		default:
			// We're running on g's stack, so either we're not in a signal handler or
			// the signal handler has set the correct g.  If we're on gsignal or g0,
			// systemstack will make the call directly; otherwise, it will switch to
			// g0 to ensure we have enough room to call a libc function.
			//
			// The function literal that we pass to systemstack is not nosplit, but
			// that's ok: we'll be running on a fresh, clean system stack so the stack
			// check will always succeed anyway.
			systemstack(func() {
				ret = callCgolangSigaction(uintptr(sig), new, old)
			})
		}

		const EINVAL = 22
		if ret == EINVAL {
			// libc reserves certain signals — normally 32-33 — for pthreads, and
			// returns EINVAL for sigaction calls on those signals.  If we get EINVAL,
			// fall back to making the syscall directly.
			sysSigaction(sig, new, old)
		}
	}

	if msanenabled && old != nil {
		msanread(unsafe.Pointer(old), unsafe.Sizeof(*old))
	}
	if asanenabled && old != nil {
		asanread(unsafe.Pointer(old), unsafe.Sizeof(*old))
	}
}

// callCgolangSigaction calls the sigaction function in the runtime/cgolang package
// using the GCC calling convention. It is implemented in assembly.
//
//golang:noescape
func callCgolangSigaction(sig uintptr, new, old *sigactiont) int32
