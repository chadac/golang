// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build (linux && (amd64 || arm64 || loong64 || ppc64le)) || (freebsd && amd64)

#include <errno.h>
#include <stdint.h>
#include <stdlib.h>
#include <sys/mman.h>

#include "libcgolang.h"

uintptr_t
x_cgolang_mmap(void *addr, uintptr_t length, int32_t prot, int32_t flags, int32_t fd, uint32_t offset) {
	void *p;

	_cgolang_tsan_acquire();
	p = mmap(addr, length, prot, flags, fd, offset);
	_cgolang_tsan_release();
	if (p == MAP_FAILED) {
		/* This is what the Golang code expects on failure.  */
		return (uintptr_t)errno;
	}
	return (uintptr_t)p;
}

void
x_cgolang_munmap(void *addr, uintptr_t length) {
	int r;

	_cgolang_tsan_acquire();
	r = munmap(addr, length);
	_cgolang_tsan_release();
	if (r < 0) {
		/* The Golang runtime is not prepared for munmap to fail.  */
		abort();
	}
}
