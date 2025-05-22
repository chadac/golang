// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build asan

#include "golang_asm.h"
#include "golang_tls.h"
#include "funcdata.h"
#include "textflag.h"

// This is like race_amd64.s, but for the asan calls.
// See race_amd64.s for detailed comments.

#ifdef GOOS_windows
#define RARG0 CX
#define RARG1 DX
#define RARG2 R8
#define RARG3 R9
#else
#define RARG0 DI
#define RARG1 SI
#define RARG2 DX
#define RARG3 CX
#endif

// Called from instrumented code.
// func runtime·doasanread(addr unsafe.Pointer, sz, sp, pc uintptr)
TEXT	runtime·doasanread(SB), NOSPLIT, $0-32
	MOVQ	addr+0(FP), RARG0
	MOVQ	sz+8(FP), RARG1
	MOVQ	sp+16(FP), RARG2
	MOVQ	pc+24(FP), RARG3
	// void __asan_read_golang(void *addr, uintptr_t sz, void *sp, void *pc);
	MOVQ	$__asan_read_golang(SB), AX
	JMP	asancall<>(SB)

// func runtime·doasanwrite(addr unsafe.Pointer, sz, sp, pc uintptr)
TEXT	runtime·doasanwrite(SB), NOSPLIT, $0-32
	MOVQ	addr+0(FP), RARG0
	MOVQ	sz+8(FP), RARG1
	MOVQ	sp+16(FP), RARG2
	MOVQ	pc+24(FP), RARG3
	// void __asan_write_golang(void *addr, uintptr_t sz, void *sp, void *pc);
	MOVQ	$__asan_write_golang(SB), AX
	JMP	asancall<>(SB)

// func runtime·asanunpoison(addr unsafe.Pointer, sz uintptr)
TEXT	runtime·asanunpoison(SB), NOSPLIT, $0-16
	MOVQ	addr+0(FP), RARG0
	MOVQ	sz+8(FP), RARG1
	// void __asan_unpoison_golang(void *addr, uintptr_t sz);
	MOVQ	$__asan_unpoison_golang(SB), AX
	JMP	asancall<>(SB)

// func runtime·asanpoison(addr unsafe.Pointer, sz uintptr)
TEXT	runtime·asanpoison(SB), NOSPLIT, $0-16
	MOVQ	addr+0(FP), RARG0
	MOVQ	sz+8(FP), RARG1
	// void __asan_poison_golang(void *addr, uintptr_t sz);
	MOVQ	$__asan_poison_golang(SB), AX
	JMP	asancall<>(SB)

// func runtime·asanregisterglobals(addr unsafe.Pointer, n uintptr)
TEXT	runtime·asanregisterglobals(SB), NOSPLIT, $0-16
	MOVQ	addr+0(FP), RARG0
	MOVQ	n+8(FP), RARG1
	// void __asan_register_globals_golang(void *addr, uintptr_t n);
	MOVQ	$__asan_register_globals_golang(SB), AX
	JMP	asancall<>(SB)

// func runtime·lsanregisterrootregion(addr unsafe.Pointer, n uintptr)
TEXT	runtime·lsanregisterrootregion(SB), NOSPLIT, $0-16
	MOVQ	addr+0(FP), RARG0
	MOVQ	n+8(FP), RARG1
	// void __lsan_register_root_region_golang(void *addr, uintptr_t sz)
	MOVQ	$__lsan_register_root_region_golang(SB), AX
	JMP	asancall<>(SB)

// func runtime·lsanunregisterrootregion(addr unsafe.Pointer, n uintptr)
TEXT	runtime·lsanunregisterrootregion(SB), NOSPLIT, $0-16
	MOVQ	addr+0(FP), RARG0
	MOVQ	n+8(FP), RARG1
	// void __lsan_unregister_root_region_golang(void *addr, uintptr_t sz)
	MOVQ	$__lsan_unregister_root_region_golang(SB), AX
	JMP	asancall<>(SB)

// func runtime·lsandoleakcheck()
TEXT	runtime·lsandoleakcheck(SB), NOSPLIT, $0-0
	// void __lsan_do_leak_check_golang(void);
	MOVQ	$__lsan_do_leak_check_golang(SB), AX
	JMP	asancall<>(SB)

// Switches SP to g0 stack and calls (AX). Arguments already set.
TEXT	asancall<>(SB), NOSPLIT, $0-0
	get_tls(R12)
	MOVQ	g(R12), R14
	MOVQ	SP, R12		// callee-saved, preserved across the CALL
	CMPQ	R14, $0
	JE	call	// no g; still on a system stack

	MOVQ	g_m(R14), R13

	// Switch to g0 stack if we aren't already on g0 or gsignal.
	MOVQ	m_gsignal(R13), R10
	CMPQ	R10, R14
	JE	call	// already on gsignal

	MOVQ	m_g0(R13), R10
	CMPQ	R10, R14
	JE	call	// already on g0

	MOVQ	(g_sched+golangbuf_sp)(R10), SP
call:
	ANDQ	$~15, SP	// alignment for gcc ABI
	CALL	AX
	MOVQ	R12, SP
	RET
