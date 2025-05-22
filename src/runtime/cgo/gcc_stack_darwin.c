// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include <pthread.h>
#include "libcgolang.h"

void
x_cgolang_getstackbound(uintptr bounds[2])
{
	void* addr;
	size_t size;
	pthread_t p;

	p = pthread_self();
	addr = pthread_get_stackaddr_np(p); // high address (!)
	size = pthread_get_stacksize_np(p);

	// bounds points into the Golang stack. TSAN can't see the synchronization
	// in Golang around stack reuse.
	_cgolang_tsan_acquire();
	bounds[0] = (uintptr)addr - size;
	bounds[1] = (uintptr)addr;
	_cgolang_tsan_release();
}
