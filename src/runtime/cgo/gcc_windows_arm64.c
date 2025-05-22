// Copyright 2009 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <process.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include "libcgolang.h"
#include "libcgolang_windows.h"

static unsigned long __stdcall threadentry(void*);
static void (*setg_gcc)(void*);

void
x_cgolang_init(G *g, void (*setg)(void*))
{
	setg_gcc = setg;
}

void
_cgolang_sys_thread_start(ThreadStart *ts)
{
	_cgolang_beginthread(threadentry, ts);
}

extern void crosscall1(void (*fn)(void), void (*setg_gcc)(void*), void *g);

static unsigned long
__stdcall
threadentry(void *v)
{
	ThreadStart ts;

	ts = *(ThreadStart*)v;
	free(v);

	crosscall1(ts.fn, setg_gcc, (void *)ts.g);
	return 0;
}
