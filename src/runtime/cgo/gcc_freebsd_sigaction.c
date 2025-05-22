// Copyright 2018 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build freebsd && amd64

#include <errno.h>
#include <stddef.h>
#include <stdint.h>
#include <string.h>
#include <signal.h>

#include "libcgolang.h"

// golang_sigaction_t is a C version of the sigactiont struct from
// os_freebsd.golang.  This definition — and its conversion to and from struct
// sigaction — are specific to freebsd/amd64.
typedef struct {
        uint32_t __bits[_SIG_WORDS];
} golang_sigset_t;
typedef struct {
	uintptr_t handler;
	int32_t flags;
	golang_sigset_t mask;
} golang_sigaction_t;

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
			if (golangact->mask.__bits[i/32] & ((uint32_t)(1)<<(i&31))) {
				sigaddset(&act.sa_mask, i+1);
			}
		}
		act.sa_flags = golangact->flags;
	}

	ret = sigaction(signum, golangact ? &act : NULL, oldgolangact ? &oldact : NULL);
	if (ret == -1) {
		// runtime.sigaction expects _cgolang_sigaction to return errno on error.
		_cgolang_tsan_release();
		return errno;
	}

	if (oldgolangact) {
		if (oldact.sa_flags & SA_SIGINFO) {
			oldgolangact->handler = (uintptr_t)(oldact.sa_sigaction);
		} else {
			oldgolangact->handler = (uintptr_t)(oldact.sa_handler);
		}
		for (i = 0 ; i < _SIG_WORDS; i++) {
			oldgolangact->mask.__bits[i] = 0;
		}
		for (i = 0; i < 8 * sizeof(oldgolangact->mask); i++) {
			if (sigismember(&oldact.sa_mask, i+1) == 1) {
				oldgolangact->mask.__bits[i/32] |= (uint32_t)(1)<<(i&31);
			}
		}
		oldgolangact->flags = oldact.sa_flags;
	}

	_cgolang_tsan_release();
	return ret;
}
