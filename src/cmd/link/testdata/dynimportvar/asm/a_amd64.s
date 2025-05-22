// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build darwin

TEXT ·Mach_task_self(SB),0,$0-4
	MOVQ	$libc_mach_task_self_(SB), AX
	MOVQ	(AX), AX
	MOVL	AX, ret+0(FP)
	RET
