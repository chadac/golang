// Copyright 2019 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "cgolang/abi_arm64.h"

TEXT _rt0_arm64_netbsd(SB),NOSPLIT|NOFRAME,$0
	MOVD	0(RSP), R0	// argc
	ADD	$8, RSP, R1	// argv
	BL	main(SB)

// When building with -buildmode=c-shared, this symbol is called when the shared
// library is loaded.
TEXT _rt0_arm64_netbsd_lib(SB),NOSPLIT,$184
	// Preserve callee-save registers.
	SAVE_R19_TO_R28(24)
	SAVE_F8_TO_F15(104)

	// Initialize g as null in case of using g later e.g. sigaction in cgolang_sigaction.golang
	MOVD	ZR, g

	MOVD	R0, _rt0_arm64_netbsd_lib_argc<>(SB)
	MOVD	R1, _rt0_arm64_netbsd_lib_argv<>(SB)

	// Synchronous initialization.
	MOVD	$runtime·libpreinit(SB), R4
	BL	(R4)

	// Create a new thread to do the runtime initialization and return.
	MOVD	_cgolang_sys_thread_create(SB), R4
	CBZ	R4, nocgolang
	MOVD	$_rt0_arm64_netbsd_lib_golang(SB), R0
	MOVD	$0, R1
	SUB	$16, RSP		// reserve 16 bytes for sp-8 where fp may be saved.
	BL	(R4)
	ADD	$16, RSP
	B	restore

nocgolang:
	MOVD	$0x800000, R0                     // stacksize = 8192KB
	MOVD	$_rt0_arm64_netbsd_lib_golang(SB), R1
	MOVD	R0, 8(RSP)
	MOVD	R1, 16(RSP)
	MOVD	$runtime·newosproc0(SB),R4
	BL	(R4)

restore:
	// Restore callee-save registers.
	RESTORE_R19_TO_R28(24)
	RESTORE_F8_TO_F15(104)
	RET

TEXT _rt0_arm64_netbsd_lib_golang(SB),NOSPLIT,$0
	MOVD	_rt0_arm64_netbsd_lib_argc<>(SB), R0
	MOVD	_rt0_arm64_netbsd_lib_argv<>(SB), R1
	MOVD	$runtime·rt0_golang(SB),R4
	B       (R4)

DATA _rt0_arm64_netbsd_lib_argc<>(SB)/8, $0
GLOBL _rt0_arm64_netbsd_lib_argc<>(SB),NOPTR, $8
DATA _rt0_arm64_netbsd_lib_argv<>(SB)/8, $0
GLOBL _rt0_arm64_netbsd_lib_argv<>(SB),NOPTR, $8


TEXT main(SB),NOSPLIT|NOFRAME,$0
	MOVD	$runtime·rt0_golang(SB), R2
	BL	(R2)
exit:
	MOVD	$0, R0
	SVC	$1	// sys_exit
