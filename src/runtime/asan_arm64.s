// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build asan

#include "golang_asm.h"
#include "textflag.h"

#define RARG0 R0
#define RARG1 R1
#define RARG2 R2
#define RARG3 R3
#define FARG R4

// Called from instrumented code.
// func runtime·doasanread(addr unsafe.Pointer, sz, sp, pc uintptr)
TEXT	runtime·doasanread(SB), NOSPLIT, $0-32
	MOVD	addr+0(FP), RARG0
	MOVD	sz+8(FP), RARG1
	MOVD	sp+16(FP), RARG2
	MOVD	pc+24(FP), RARG3
	// void __asan_read_golang(void *addr, uintptr_t sz, void *sp, void *pc);
	MOVD	$__asan_read_golang(SB), FARG
	JMP	asancall<>(SB)

// func runtime·doasanwrite(addr unsafe.Pointer, sz, sp, pc uintptr)
TEXT	runtime·doasanwrite(SB), NOSPLIT, $0-32
	MOVD	addr+0(FP), RARG0
	MOVD	sz+8(FP), RARG1
	MOVD	sp+16(FP), RARG2
	MOVD	pc+24(FP), RARG3
	// void __asan_write_golang(void *addr, uintptr_t sz, void *sp, void *pc);
	MOVD	$__asan_write_golang(SB), FARG
	JMP	asancall<>(SB)

// func runtime·asanunpoison(addr unsafe.Pointer, sz uintptr)
TEXT	runtime·asanunpoison(SB), NOSPLIT, $0-16
	MOVD	addr+0(FP), RARG0
	MOVD	sz+8(FP), RARG1
	// void __asan_unpoison_golang(void *addr, uintptr_t sz);
	MOVD	$__asan_unpoison_golang(SB), FARG
	JMP	asancall<>(SB)

// func runtime·asanpoison(addr unsafe.Pointer, sz uintptr)
TEXT	runtime·asanpoison(SB), NOSPLIT, $0-16
	MOVD	addr+0(FP), RARG0
	MOVD	sz+8(FP), RARG1
	// void __asan_poison_golang(void *addr, uintptr_t sz);
	MOVD	$__asan_poison_golang(SB), FARG
	JMP	asancall<>(SB)

// func runtime·asanregisterglobals(addr unsafe.Pointer, n uintptr)
TEXT	runtime·asanregisterglobals(SB), NOSPLIT, $0-16
	MOVD	addr+0(FP), RARG0
	MOVD	n+8(FP), RARG1
	// void __asan_register_globals_golang(void *addr, uintptr_t n);
	MOVD	$__asan_register_globals_golang(SB), FARG
	JMP	asancall<>(SB)

// func runtime·lsanregisterrootregion(addr unsafe.Pointer, n uintptr)
TEXT	runtime·lsanregisterrootregion(SB), NOSPLIT, $0-16
	MOVD	addr+0(FP), RARG0
	MOVD	n+8(FP), RARG1
	// void __lsan_register_root_region_golang(void *addr, uintptr_t n);
	MOVD	$__lsan_register_root_region_golang(SB), FARG
	JMP	asancall<>(SB)

// func runtime·lsanunregisterrootregion(addr unsafe.Pointer, n uintptr)
TEXT	runtime·lsanunregisterrootregion(SB), NOSPLIT, $0-16
	MOVD	addr+0(FP), RARG0
	MOVD	n+8(FP), RARG1
	// void __lsan_unregister_root_region_golang(void *addr, uintptr_t n);
	MOVD	$__lsan_unregister_root_region_golang(SB), FARG
	JMP	asancall<>(SB)

// func runtime·lsandoleakcheck()
TEXT	runtime·lsandoleakcheck(SB), NOSPLIT, $0-0
	// void __lsan_do_leak_check_golang(void);
	MOVD	$__lsan_do_leak_check_golang(SB), FARG
	JMP	asancall<>(SB)

// Switches SP to g0 stack and calls (FARG). Arguments already set.
TEXT	asancall<>(SB), NOSPLIT, $0-0
	MOVD	RSP, R19                  // callee-saved
	CBZ	g, call                   // no g, still on a system stack
	MOVD	g_m(g), R10

	// Switch to g0 stack if we aren't already on g0 or gsignal.
	MOVD	m_gsignal(R10), R11
	CMP	R11, g
	BEQ	call

	MOVD	m_g0(R10), R11
	CMP	R11, g
	BEQ	call

	MOVD	(g_sched+golangbuf_sp)(R11), R5
	MOVD	R5, RSP

call:
	BL	(FARG)
	MOVD	R19, RSP
	RET
