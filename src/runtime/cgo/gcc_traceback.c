// Copyright 2016 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build darwin || linux

#include <stdint.h>
#include "libcgolang.h"

#ifndef __has_feature
#define __has_feature(x) 0
#endif

#if __has_feature(memory_sanitizer)
#include <sanitizer/msan_interface.h>
#endif

// Call the user's traceback function and then call sigtramp.
// The runtime signal handler will jump to this code.
// We do it this way so that the user's traceback function will be called
// by a C function with proper unwind info.
void
x_cgolang_callers(uintptr_t sig, void *info, void *context, void (*cgolangTraceback)(struct cgolangTracebackArg*), uintptr_t* cgolangCallers, void (*sigtramp)(uintptr_t, void*, void*)) {
	struct cgolangTracebackArg arg;

	arg.Context = 0;
	arg.SigContext = (uintptr_t)(context);
	arg.Buf = cgolangCallers;
	arg.Max = 32; // must match len(runtime.cgolangCallers)

#if __has_feature(memory_sanitizer)
        // This function is called directly from the signal handler.
        // The arguments are passed in registers, so whether msan
        // considers cgolangCallers to be initialized depends on whether
        // it considers the appropriate register to be initialized.
        // That can cause false reports in rare cases.
        // Explicitly unpoison the memory to avoid that.
        // See issue #47543 for more details.
        __msan_unpoison(&arg, sizeof arg);
#endif

	_cgolang_tsan_acquire();
	(*cgolangTraceback)(&arg);
	_cgolang_tsan_release();
	sigtramp(sig, info, context);
}
