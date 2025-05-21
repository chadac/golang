// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build mips || mipsle

#include "golang_asm.h"
#include "golang_tls.h"
#include "funcdata.h"
#include "textflag.h"

// If !iscgolang, this is a no-op.
// NOTE: golanggolang assumes load_g only clobers g (R30) and REGTMP (R23)
TEXT runtime·save_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVB	runtime·iscgolang(SB), R23
	BEQ	R23, nocgolang

	MOVW	R3, R23
	MOVW	g, runtime·tls_g(SB) // TLS relocation clobbers R3
	MOVW	R23, R3

nocgolang:
	RET

TEXT runtime·load_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	runtime·tls_g(SB), g // TLS relocation clobbers R3
	RET

GLOBL runtime·tls_g(SB), TLSBSS, $4
