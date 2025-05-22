// Copyright 2020 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// The C definitions for traceback.golang. That file uses //export so
// it can't put function definitions in the "C" import comment.

#include <stdint.h>

char *p;

int crashInGolang;
extern void h1(void);

int tracebackF3(void) {
	if (crashInGolang)
		h1();
	else
		*p = 0;
	return 0;
}

int tracebackF2(void) {
	return tracebackF3();
}

int tracebackF1(void) {
	return tracebackF2();
}

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

void cgolangTraceback(void* parg) {
	struct cgolangTracebackArg* arg = (struct cgolangTracebackArg*)(parg);
	arg->buf[0] = 1;
	arg->buf[1] = 2;
	arg->buf[2] = 3;
	arg->buf[3] = 0;
}

void cgolangSymbolizer(void* parg) {
	struct cgolangSymbolizerArg* arg = (struct cgolangSymbolizerArg*)(parg);
	if (arg->pc != arg->data + 1) {
		arg->file = "unexpected data";
	} else {
		arg->file = "cgolang symbolizer";
	}
	arg->lineno = arg->data + 1;
	arg->data++;
}
