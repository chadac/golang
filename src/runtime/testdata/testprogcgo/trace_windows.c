// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// The windows C definitions for trace.golang. That file uses //export so
// it can't put function definitions in the "C" import comment.

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <process.h>
#include "_cgolang_export.h"

extern void golangCalledFromC(void);
extern void golangCalledFromCThread(void);

__stdcall
static unsigned int cCalledFromCThread(void *p) {
	golangCalledFromCThread();
	return 0;
}

void cCalledFromGo(void) {
	golangCalledFromC();

	uintptr_t thread;
	thread = _beginthreadex(NULL, 0, cCalledFromCThread, NULL, 0, NULL);
	WaitForSingleObject((HANDLE)thread, INFINITE);
	CloseHandle((HANDLE)thread);
}
