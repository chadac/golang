// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build darwin

TEXT ·Mach_task_self(SB),0,$0-4
	MOVD	$libc_mach_task_self_(SB), R0
	MOVD	(R0), R0
	MOVW	R0, ret+0(FP)
	RET
