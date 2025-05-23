// Copyright 2009 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build linux && (386 || arm || loong64 || mips || mipsle || mips64 || mips64le || riscv64)

#include <pthread.h>
#include <string.h>
#include <signal.h>
#include "libcgolang.h"
#include "libcgolang_unix.h"

static void *threadentry(void*);

void (*x_cgolang_inittls)(void **tlsg, void **tlsbase) __attribute__((common));
static void (*setg_gcc)(void*);

void
x_cgolang_init(G *g, void (*setg)(void*), void **tlsg, void **tlsbase)
{
	setg_gcc = setg;

	_cgolang_set_stacklo(g, NULL);

	if (x_cgolang_inittls) {
		x_cgolang_inittls(tlsg, tlsbase);
	}
}

void
_cgolang_sys_thread_start(ThreadStart *ts)
{
	pthread_attr_t attr;
	sigset_t ign, oset;
	pthread_t p;
	size_t size;
	int err;

	sigfillset(&ign);
	pthread_sigmask(SIG_SETMASK, &ign, &oset);

	pthread_attr_init(&attr);
	pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
	pthread_attr_getstacksize(&attr, &size);
	// Leave stacklo=0 and set stackhi=size; mstart will do the rest.
	ts->g->stackhi = size;
	err = _cgolang_try_pthread_create(&p, &attr, threadentry, ts);

	pthread_sigmask(SIG_SETMASK, &oset, nil);

	if (err != 0) {
		fatalf("pthread_create failed: %s", strerror(err));
	}
}

extern void crosscall1(void (*fn)(void), void (*setg_gcc)(void*), void *g);
static void*
threadentry(void *v)
{
	ThreadStart ts;

	ts = *(ThreadStart*)v;
	free(v);

	crosscall1(ts.fn, setg_gcc, ts.g);
	return nil;
}
