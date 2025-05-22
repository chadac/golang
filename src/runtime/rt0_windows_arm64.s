// Copyright 2018 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "golang_asm.h"
#include "golang_tls.h"
#include "textflag.h"

// This is the entry point for the program from the
// kernel for an ordinary -buildmode=exe program.
TEXT _rt0_arm64_windows(SB),NOSPLIT|NOFRAME,$0
	B	·rt0_golang(SB)

TEXT _rt0_arm64_windows_lib(SB),NOSPLIT|NOFRAME,$0
	MOVD	$_rt0_arm64_windows_lib_golang(SB), R0
	MOVD	$0, R1
	MOVD	_cgolang_sys_thread_create(SB), R2
	B	(R2)

TEXT _rt0_arm64_windows_lib_golang(SB),NOSPLIT|NOFRAME,$0
	MOVD	$0, R0
	MOVD	$0, R1
	MOVD	$runtime·rt0_golang(SB), R2
	B	(R2)

TEXT main(SB),NOSPLIT|NOFRAME,$0
	MOVD	$runtime·rt0_golang(SB), R2
	B	(R2)

