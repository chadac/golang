// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "golang_asm.h"
#include "golang_tls.h"
#include "textflag.h"

TEXT _rt0_amd64_windows(SB),NOSPLIT|NOFRAME,$-8
	JMP	_rt0_amd64(SB)

// When building with -buildmode=(c-shared or c-archive), this
// symbol is called. For dynamic libraries it is called when the
// library is loaded. For static libraries it is called when the
// final executable starts, during the C runtime initialization
// phase.
// Leave space for four pointers on the stack as required
// by the Windows amd64 calling convention.
TEXT _rt0_amd64_windows_lib(SB),NOSPLIT|NOFRAME,$40
	// Create a new thread to do the runtime initialization and return.
	MOVQ	BX, 32(SP) // callee-saved, preserved across the CALL
	MOVQ	SP, BX
	ANDQ	$~15, SP // alignment as per Windows requirement
	MOVQ	_cgolang_sys_thread_create(SB), AX
	MOVQ	$_rt0_amd64_windows_lib_golang(SB), CX
	MOVQ	$0, DX
	CALL	AX
	MOVQ	BX, SP
	MOVQ	32(SP), BX
	RET

TEXT _rt0_amd64_windows_lib_golang(SB),NOSPLIT|NOFRAME,$0
	MOVQ  $0, DI
	MOVQ	$0, SI
	MOVQ	$runtime·rt0_golang(SB), AX
	JMP	AX
