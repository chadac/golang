// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#ifdef __CYGWIN__
#error "don't use the cygwin compiler to build native Windows programs; use MinGW instead"
#endif

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#include <stdio.h>
#include <stdlib.h>

#include "libcgolang.h"
#include "libcgolang_windows.h"

// Ensure there's one symbol marked __declspec(dllexport).
// If there are no exported symbols, the unfortunate behavior of
// the binutils linker is to also strip the relocations table,
// resulting in non-PIE binary. The other option is the
// --export-all-symbols flag, but we don't need to export all symbols
// and this may overflow the export table (#40795).
// See https://sourceware.org/bugzilla/show_bug.cgi?id=19011
__declspec(dllexport) int _cgolang_dummy_export;

static volatile LONG runtime_init_once_gate = 0;
static volatile LONG runtime_init_once_done = 0;

static CRITICAL_SECTION runtime_init_cs;

static HANDLE runtime_init_wait;
static int runtime_init_done;

uintptr_t x_cgolang_pthread_key_created;
void (*x_crosscall2_ptr)(void (*fn)(void *), void *, int, size_t);

// Pre-initialize the runtime synchronization objects
void
_cgolang_preinit_init() {
	 runtime_init_wait = CreateEvent(NULL, TRUE, FALSE, NULL);
	 if (runtime_init_wait == NULL) {
		fprintf(stderr, "runtime: failed to create runtime initialization wait event.\n");
		abort();
	 }

	 InitializeCriticalSection(&runtime_init_cs);
}

// Make sure that the preinit sequence has run.
void
_cgolang_maybe_run_preinit() {
	 if (!InterlockedExchangeAdd(&runtime_init_once_done, 0)) {
			if (InterlockedIncrement(&runtime_init_once_gate) == 1) {
				 _cgolang_preinit_init();
				 InterlockedIncrement(&runtime_init_once_done);
			} else {
				 // Decrement to avoid overflow.
				 InterlockedDecrement(&runtime_init_once_gate);
				 while(!InterlockedExchangeAdd(&runtime_init_once_done, 0)) {
						Sleep(0);
				 }
			}
	 }
}

void
x_cgolang_sys_thread_create(unsigned long (__stdcall *func)(void*), void* arg) {
	_cgolang_beginthread(func, arg);
}

int
_cgolang_is_runtime_initialized() {
	 int status;

	 EnterCriticalSection(&runtime_init_cs);
	 status = runtime_init_done;
	 LeaveCriticalSection(&runtime_init_cs);
	 return status;
}

uintptr_t
_cgolang_wait_runtime_init_done(void) {
	void (*pfn)(struct context_arg*);

	 _cgolang_maybe_run_preinit();
	while (!_cgolang_is_runtime_initialized()) {
			WaitForSingleObject(runtime_init_wait, INFINITE);
	}
	pfn = _cgolang_get_context_function();
	if (pfn != nil) {
		struct context_arg arg;

		arg.Context = 0;
		(*pfn)(&arg);
		return arg.Context;
	}
	return 0;
}

// Should not be used since x_cgolang_pthread_key_created will always be zero.
void x_cgolang_bindm(void* dummy) {
	fprintf(stderr, "unexpected cgolang_bindm on Windows\n");
	abort();
}

void
x_cgolang_notify_runtime_init_done(void* dummy) {
	 _cgolang_maybe_run_preinit();

	 EnterCriticalSection(&runtime_init_cs);
	runtime_init_done = 1;
	 LeaveCriticalSection(&runtime_init_cs);

	 if (!SetEvent(runtime_init_wait)) {
		fprintf(stderr, "runtime: failed to signal runtime initialization complete.\n");
		abort();
	}
}

// The context function, used when tracing back C calls into Golang.
static void (*cgolang_context_function)(struct context_arg*);

// Sets the context function to call to record the traceback context
// when calling a Golang function from C code. Called from runtime.SetCgolangTraceback.
void x_cgolang_set_context_function(void (*context)(struct context_arg*)) {
	EnterCriticalSection(&runtime_init_cs);
	cgolang_context_function = context;
	LeaveCriticalSection(&runtime_init_cs);
}

// Gets the context function.
void (*(_cgolang_get_context_function(void)))(struct context_arg*) {
	void (*ret)(struct context_arg*);

	EnterCriticalSection(&runtime_init_cs);
	ret = cgolang_context_function;
	LeaveCriticalSection(&runtime_init_cs);
	return ret;
}

void _cgolang_beginthread(unsigned long (__stdcall *func)(void*), void* arg) {
	int tries;
	HANDLE thandle;

	for (tries = 0; tries < 20; tries++) {
		thandle = CreateThread(NULL, 0, func, arg, 0, NULL);
		if (thandle == 0 && GetLastError() == ERROR_NOT_ENOUGH_MEMORY) {
			// "Insufficient resources", try again in a bit.
			//
			// Note that the first Sleep(0) is a yield.
			Sleep(tries); // milliseconds
			continue;
		} else if (thandle == 0) {
			break;
		}
		CloseHandle(thandle);
		return; // Success!
	}

	fprintf(stderr, "runtime: failed to create new OS thread (%lu)\n", GetLastError());
	abort();
}
