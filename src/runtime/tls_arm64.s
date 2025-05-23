// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "golang_asm.h"
#include "golang_tls.h"
#include "funcdata.h"
#include "textflag.h"
#include "tls_arm64.h"

TEXT runtime·load_g(SB),NOSPLIT,$0
#ifndef GOOS_darwin
#ifndef GOOS_openbsd
#ifndef GOOS_windows
	MOVB	runtime·iscgolang(SB), R0
	CBZ	R0, nocgolang
#endif
#endif
#endif

	MRS_TPIDR_R0
#ifdef TLS_darwin
	// Darwin sometimes returns unaligned pointers
	AND	$0xfffffffffffffff8, R0
#endif
	MOVD	runtime·tls_g(SB), R27
	MOVD	(R0)(R27), g

nocgolang:
	RET

TEXT runtime·save_g(SB),NOSPLIT,$0
#ifndef GOOS_darwin
#ifndef GOOS_openbsd
#ifndef GOOS_windows
	MOVB	runtime·iscgolang(SB), R0
	CBZ	R0, nocgolang
#endif
#endif
#endif

	MRS_TPIDR_R0
#ifdef TLS_darwin
	// Darwin sometimes returns unaligned pointers
	AND	$0xfffffffffffffff8, R0
#endif
	MOVD	runtime·tls_g(SB), R27
	MOVD	g, (R0)(R27)

nocgolang:
	RET

#ifdef TLSG_IS_VARIABLE
#ifdef GOOS_android
// Use the free TLS_SLOT_APP slot #2 on Android Q.
// Earlier androids are set up in gcc_android.c.
DATA runtime·tls_g+0(SB)/8, $16
#endif
GLOBL runtime·tls_g+0(SB), NOPTR, $8
#else
GLOBL runtime·tls_g+0(SB), TLSBSS, $8
#endif
