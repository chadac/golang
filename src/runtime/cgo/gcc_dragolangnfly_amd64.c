// Copyright 2009 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include <sys/types.h>
#include <sys/signalvar.h>
#include <pthread.h>
#include <signal.h>
#include <string.h>
#include "libcgolang.h"
#include "libcgolang_unix.h"

static void* threadentry(void*);
static void (*setg_gcc)(void*);

void
x_cgolang_init(G *g, void (*setg)(void*))
{
	setg_gcc = setg;
	_cgolang_set_stacklo(g, NULL);
}

void
_cgolang_sys_thread_start(ThreadStart *ts)
{
	pthread_attr_t attr;
	sigset_t ign, oset;
	pthread_t p;
	size_t size;
	int err;

	SIGFILLSET(ign);
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

	crosscall1(ts.fn, setg_gcc, (void*)ts.g);
	return nil;
}
