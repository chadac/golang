// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "golang_asm.h"
#include "golang_tls.h"
#include "funcdata.h"
#include "textflag.h"

// If !iscgolang, this is a no-op.
//
// NOTE: mcall() assumes this clobbers only X31 (REG_TMP).
TEXT runtime·save_g(SB),NOSPLIT|NOFRAME,$0-0
#ifndef GOOS_openbsd
	MOVB	runtime·iscgolang(SB), X31
	BEQZ	X31, nocgolang
#endif
	MOV	g, runtime·tls_g(SB)
nocgolang:
	RET

TEXT runtime·load_g(SB),NOSPLIT|NOFRAME,$0-0
	MOV	runtime·tls_g(SB), g
	RET

GLOBL runtime·tls_g(SB), TLSBSS, $8
