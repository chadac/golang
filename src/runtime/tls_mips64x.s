// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build mips64 || mips64le

#include "golang_asm.h"
#include "golang_tls.h"
#include "funcdata.h"
#include "textflag.h"

// If !iscgolang, this is a no-op.
//
// NOTE: mcall() assumes this clobbers only R23 (REGTMP).
TEXT runtime·save_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVB	runtime·iscgolang(SB), R23
	BEQ	R23, nocgolang

	MOVV	R3, R23	// save R3
	MOVV	g, runtime·tls_g(SB) // TLS relocation clobbers R3
	MOVV	R23, R3	// restore R3

nocgolang:
	RET

TEXT runtime·load_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVV	runtime·tls_g(SB), g // TLS relocation clobbers R3
	RET

GLOBL runtime·tls_g(SB), TLSBSS, $8
