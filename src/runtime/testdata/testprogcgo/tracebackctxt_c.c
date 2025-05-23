// Copyright 2016 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// The C definitions for tracebackctxt.golang. That file uses //export so
// it can't put function definitions in the "C" import comment.

#include <stdlib.h>
#include <stdint.h>

// Functions exported from Golang.
extern void G1(void);
extern void G2(void);
extern void TracebackContextPreemptionGolangFunction(int);
extern void TracebackContextProfileGolangFunction(void);

void C1() {
	G1();
}

void C2() {
	G2();
}

struct cgolangContextArg {
	uintptr_t context;
};

struct cgolangTracebackArg {
	uintptr_t  context;
	uintptr_t  sigContext;
	uintptr_t* buf;
	uintptr_t  max;
};

struct cgolangSymbolizerArg {
	uintptr_t   pc;
	const char* file;
	uintptr_t   lineno;
	const char* func;
	uintptr_t   entry;
	uintptr_t   more;
	uintptr_t   data;
};

// Uses atomic adds and subtracts to catch the possibility of
// erroneous calls from multiple threads; that should be impossible in
// this test case, but we check just in case.
static int contextCount;

int getContextCount() {
	return __sync_add_and_fetch(&contextCount, 0);
}

void tcContext(void* parg) {
	struct cgolangContextArg* arg = (struct cgolangContextArg*)(parg);
	if (arg->context == 0) {
		arg->context = __sync_add_and_fetch(&contextCount, 1);
	} else {
		if (arg->context != __sync_add_and_fetch(&contextCount, 0)) {
			abort();
		}
		__sync_sub_and_fetch(&contextCount, 1);
	}
}

void tcContextSimple(void* parg) {
	struct cgolangContextArg* arg = (struct cgolangContextArg*)(parg);
	if (arg->context == 0) {
		arg->context = 1;
	}
}

void tcTraceback(void* parg) {
	int base, i;
	struct cgolangTracebackArg* arg = (struct cgolangTracebackArg*)(parg);
	if (arg->context == 0 && arg->sigContext == 0) {
		// This shouldn't happen in this program.
		abort();
	}
	// Return a variable number of PC values.
	base = arg->context << 8;
	for (i = 0; i < arg->context; i++) {
		if (i < arg->max) {
			arg->buf[i] = base + i;
		}
	}
}

void tcSymbolizer(void *parg) {
	struct cgolangSymbolizerArg* arg = (struct cgolangSymbolizerArg*)(parg);
	if (arg->pc == 0) {
		return;
	}
	// Report two lines per PC returned by traceback, to test more handling.
	arg->more = arg->file == NULL;
	arg->file = "tracebackctxt.golang";
	arg->func = "cFunction";
	arg->lineno = arg->pc + (arg->more << 16);
}

void TracebackContextPreemptionCallGolang(int i) {
	TracebackContextPreemptionGolangFunction(i);
}

void TracebackContextProfileCallGolang(void) {
	TracebackContextProfileGolangFunction();
}
