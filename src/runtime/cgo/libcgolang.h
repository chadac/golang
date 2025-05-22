// Copyright 2009 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

#undef nil
#define nil ((void*)0)
#define nelem(x) (sizeof(x)/sizeof((x)[0]))

typedef uint32_t uint32;
typedef uint64_t uint64;
typedef uintptr_t uintptr;

/*
 * The beginning of the per-golangroutine structure,
 * as defined in ../pkg/runtime/runtime.h.
 * Just enough to edit these two fields.
 */
typedef struct G G;
struct G
{
	uintptr stacklo;
	uintptr stackhi;
};

/*
 * Arguments to the _cgolang_thread_start call.
 * Also known to ../pkg/runtime/runtime.h.
 */
typedef struct ThreadStart ThreadStart;
struct ThreadStart
{
	G *g;
	uintptr *tls;
	void (*fn)(void);
};

/*
 * Called by 5c/6c/8c world.
 * Makes a local copy of the ThreadStart and
 * calls _cgolang_sys_thread_start(ts).
 */
extern void (*_cgolang_thread_start)(ThreadStart *ts);

/*
 * Creates a new operating system thread without updating any Golang state
 * (OS dependent).
 */
extern void (*_cgolang_sys_thread_create)(void* (*func)(void*), void* arg);

/*
 * Indicates whether a dummy pthread per-thread variable is allocated.
 */
extern uintptr_t *_cgolang_pthread_key_created;

/*
 * Creates the new operating system thread (OS, arch dependent).
 */
void _cgolang_sys_thread_start(ThreadStart *ts);

/*
 * Waits for the Golang runtime to be initialized (OS dependent).
 * If runtime.SetCgolangTraceback is used to set a context function,
 * calls the context function and returns the context value.
 */
uintptr_t _cgolang_wait_runtime_init_done(void);

/*
 * Get the low and high boundaries of the stack.
 */
void x_cgolang_getstackbound(uintptr bounds[2]);

/*
 * Prints error then calls abort. For linux and android.
 */
void fatalf(const char* format, ...) __attribute__ ((noreturn));

/*
 * Registers the current mach thread port for EXC_BAD_ACCESS processing.
 */
void darwin_arm_init_thread_exception_port(void);

/*
 * Starts a mach message server processing EXC_BAD_ACCESS.
 */
void darwin_arm_init_mach_exception_handler(void);

/*
 * The cgolang context function. See runtime.SetCgolangTraceback.
 */
struct context_arg {
	uintptr_t Context;
};
extern void (*(_cgolang_get_context_function(void)))(struct context_arg*);

/*
 * The argument for the cgolang traceback callback. See runtime.SetCgolangTraceback.
 */
struct cgolangTracebackArg {
	uintptr_t  Context;
	uintptr_t  SigContext;
	uintptr_t* Buf;
	uintptr_t  Max;
};

/*
 * TSAN support.  This is only useful when building with
 *   CGO_CFLAGS="-fsanitize=thread" CGO_LDFLAGS="-fsanitize=thread" golang install
 */
#undef CGO_TSAN
#if defined(__has_feature)
# if __has_feature(thread_sanitizer)
#  define CGO_TSAN
# endif
#elif defined(__SANITIZE_THREAD__)
# define CGO_TSAN
#endif

#ifdef CGO_TSAN

// These must match the definitions in yesTsanProlog in cmd/cgolang/out.golang.
// In general we should call _cgolang_tsan_acquire when we enter C code,
// and call _cgolang_tsan_release when we return to Golang code.
// This is only necessary when calling code that might be instrumented
// by TSAN, which mostly means system library calls that TSAN intercepts.
// See the comment in cmd/cgolang/out.golang for more details.

long long _cgolang_sync __attribute__ ((common));

extern void __tsan_acquire(void*);
extern void __tsan_release(void*);

__attribute__ ((unused))
static void _cgolang_tsan_acquire() {
	__tsan_acquire(&_cgolang_sync);
}

__attribute__ ((unused))
static void _cgolang_tsan_release() {
	__tsan_release(&_cgolang_sync);
}

#else // !defined(CGO_TSAN)

#define _cgolang_tsan_acquire()
#define _cgolang_tsan_release()

#endif // !defined(CGO_TSAN)
