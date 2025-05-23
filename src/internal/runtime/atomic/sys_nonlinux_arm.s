// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !linux

#include "textflag.h"

// TODO(minux): this is only valid for ARMv6+
// func armcas(ptr *int32, old int32, new int32) bool
// Atomically:
//	if *ptr == old {
//		*ptr = new
//		return true
//	} else {
//		return false
//	}
TEXT	·Cas(SB),NOSPLIT,$0
	JMP	·armcas(SB)

// Non-linux OSes support only single processor machines before ARMv7.
// So we don't need memory barriers if golangarm < 7. And we fail loud at
// startup (runtime.checkgolangarm) if it is a multi-processor but golangarm < 7.

TEXT	·Load(SB),NOSPLIT|NOFRAME,$0-8
	MOVW	addr+0(FP), R0
	MOVW	(R0), R1

	MOVB	runtime·golangarm(SB), R11
	CMP	$7, R11
	BLT	2(PC)
	DMB	MB_ISH

	MOVW	R1, ret+4(FP)
	RET

TEXT	·Store(SB),NOSPLIT,$0-8
	MOVW	addr+0(FP), R1
	MOVW	v+4(FP), R2

	MOVB	runtime·golangarm(SB), R8
	CMP	$7, R8
	BLT	2(PC)
	DMB	MB_ISH

	MOVW	R2, (R1)

	CMP	$7, R8
	BLT	2(PC)
	DMB	MB_ISH
	RET

TEXT	·Load8(SB),NOSPLIT|NOFRAME,$0-5
	MOVW	addr+0(FP), R0
	MOVB	(R0), R1

	MOVB	runtime·golangarm(SB), R11
	CMP	$7, R11
	BLT	2(PC)
	DMB	MB_ISH

	MOVB	R1, ret+4(FP)
	RET

TEXT	·Store8(SB),NOSPLIT,$0-5
	MOVW	addr+0(FP), R1
	MOVB	v+4(FP), R2

	MOVB	runtime·golangarm(SB), R8
	CMP	$7, R8
	BLT	2(PC)
	DMB	MB_ISH

	MOVB	R2, (R1)

	CMP	$7, R8
	BLT	2(PC)
	DMB	MB_ISH
	RET

