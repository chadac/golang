// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix

// The unix C definitions for trace.golang. That file uses //export so
// it can't put function definitions in the "C" import comment.

#include <pthread.h>
#include <assert.h>

extern void golangCalledFromC(void);
extern void golangCalledFromCThread(void);

static void* cCalledFromCThread(void *p) {
	golangCalledFromCThread();
	return NULL;
}

void cCalledFromGo(void) {
	golangCalledFromC();

	pthread_t thread;
	assert(pthread_create(&thread, NULL, cCalledFromCThread, NULL) == 0);
	assert(pthread_join(thread, NULL) == 0);
}
