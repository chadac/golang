// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "golang_asm.h"
#include "golang_tls.h"
#include "funcdata.h"
#include "textflag.h"

// If !iscgolang, this is a no-op.
//
// NOTE: mcall() assumes this clobbers only R30 (REGTMP).
TEXT runtime·save_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVB	runtime·iscgolang(SB), R30
	BEQ	R30, nocgolang

	MOVV	g, runtime·tls_g(SB)

nocgolang:
	RET

TEXT runtime·load_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVV	runtime·tls_g(SB), g
	RET

GLOBL runtime·tls_g(SB), TLSBSS, $8
