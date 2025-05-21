// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build linux && (amd64 || arm64 || loong64 || ppc64le)

#include <errno.h>
#include <stddef.h>
#include <stdint.h>
#include <string.h>
#include <signal.h>

#include "libcgolang.h"

// golang_sigaction_t is a C version of the sigactiont struct from
// defs_${golangos}_${golangarch}.golang.  This definition — and its conversion
// to and from struct sigaction — are specific to ${golangos}/${golangarch}.
typedef struct {
	uintptr_t handler;
	uint64_t flags;
#ifdef __loongarch__
	uint64_t mask;
	uintptr_t restorer;
#else
	uintptr_t restorer;
	uint64_t mask;
#endif
} golang_sigaction_t;

// SA_RESTORER is part of the kernel interface.
// This is Linux i386/amd64 specific.
#ifndef SA_RESTORER
#define SA_RESTORER 0x4000000
#endif

int32_t
x_cgolang_sigaction(intptr_t signum, const golang_sigaction_t *golangact, golang_sigaction_t *oldgolangact) {
	int32_t ret;
	struct sigaction act;
	struct sigaction oldact;
	size_t i;

	_cgolang_tsan_acquire();

	memset(&act, 0, sizeof act);
	memset(&oldact, 0, sizeof oldact);

	if (golangact) {
		if (golangact->flags & SA_SIGINFO) {
			act.sa_sigaction = (void(*)(int, siginfo_t*, void*))(golangact->handler);
		} else {
			act.sa_handler = (void(*)(int))(golangact->handler);
		}
		sigemptyset(&act.sa_mask);
		for (i = 0; i < 8 * sizeof(golangact->mask); i++) {
			if (golangact->mask & ((uint64_t)(1)<<i)) {
				sigaddset(&act.sa_mask, (int)(i+1));
			}
		}
		act.sa_flags = (int)(golangact->flags & ~(uint64_t)SA_RESTORER);
	}

	ret = sigaction((int)signum, golangact ? &act : NULL, oldgolangact ? &oldact : NULL);
	if (ret == -1) {
		// runtime.rt_sigaction expects _cgolang_sigaction to return errno on error.
		_cgolang_tsan_release();
		return errno;
	}

	if (oldgolangact) {
		if (oldact.sa_flags & SA_SIGINFO) {
			oldgolangact->handler = (uintptr_t)(oldact.sa_sigaction);
		} else {
			oldgolangact->handler = (uintptr_t)(oldact.sa_handler);
		}
		oldgolangact->mask = 0;
		for (i = 0; i < 8 * sizeof(oldgolangact->mask); i++) {
			if (sigismember(&oldact.sa_mask, (int)(i+1)) == 1) {
				oldgolangact->mask |= (uint64_t)(1)<<i;
			}
		}
		oldgolangact->flags = (uint64_t)oldact.sa_flags;
	}

	_cgolang_tsan_release();
	return ret;
}
