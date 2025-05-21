// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// The C definitions for gprof.golang. That file uses //export so
// it can't put function definitions in the "C" import comment.

#include <stdint.h>
#include <stdlib.h>

// Functions exported from Go.
extern void GoSleep();

struct cgolangContextArg {
	uintptr_t context;
};

void gprofCgolangContext(void *arg) {
	((struct cgolangContextArg*)arg)->context = 1;
}

void gprofCgolangTraceback(void *arg) {
	// spend some time here so the P is more likely to be retaken.
	volatile int i;
	for (i = 0; i < 123456789; i++);
}

void CallGoSleep() {
	GoSleep();
}
