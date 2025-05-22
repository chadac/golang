// Copyright 2016 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix || windows

#include "libcgolang.h"

// Releases the cgolang traceback context.
void _cgolang_release_context(uintptr_t ctxt) {
	void (*pfn)(struct context_arg*);

	pfn = _cgolang_get_context_function();
	if (ctxt != 0 && pfn != nil) {
		struct context_arg arg;

		arg.Context = ctxt;
		(*pfn)(&arg);
	}
}
