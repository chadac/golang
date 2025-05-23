// Copyright 2009 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "libcgolang.h"

/* Stub for creating a new thread */
void
x_cgolang_thread_start(ThreadStart *arg)
{
	ThreadStart *ts;

	/* Make our own copy that can persist after we return. */
	_cgolang_tsan_acquire();
	ts = malloc(sizeof *ts);
	_cgolang_tsan_release();
	if(ts == nil) {
		fprintf(stderr, "runtime/cgolang: out of memory in thread_start\n");
		abort();
	}
	*ts = *arg;

	_cgolang_sys_thread_start(ts);	/* OS-dependent half */
}

#ifndef CGO_TSAN
void(* const _cgolang_yield)() = NULL;
#else

#include <string.h>

char x_cgolang_yield_strncpy_src = 0;
char x_cgolang_yield_strncpy_dst = 0;
size_t x_cgolang_yield_strncpy_n = 0;

/*
Stub for allowing libc interceptors to execute.

_cgolang_yield is set to NULL if we do not expect libc interceptors to exist.
*/
static void
x_cgolang_yield()
{
	/*
	The libc function(s) we call here must form a no-op and include at least one
	call that triggers TSAN to process pending asynchronous signals.

	sleep(0) would be fine, but it's not portable C (so it would need more header
	guards).
	free(NULL) has a fast-path special case in TSAN, so it doesn't
	trigger signal delivery.
	free(malloc(0)) would work (triggering the interceptors in malloc), but
	it also runs a bunch of user-supplied malloc hooks.

	So we choose strncpy(_, _, 0): it requires an extra header,
	but it's standard and should be very efficient.

	GCC 7 has an unfortunate habit of optimizing out strncpy calls (see
	https://golanglang.org/issue/21196), so the arguments here need to be global
	variables with external linkage in order to ensure that the call traps all the
	way down into libc.
	*/
	strncpy(&x_cgolang_yield_strncpy_dst, &x_cgolang_yield_strncpy_src,
	        x_cgolang_yield_strncpy_n);
}

void(* const _cgolang_yield)() = &x_cgolang_yield;

#endif  /* GO_TSAN */
