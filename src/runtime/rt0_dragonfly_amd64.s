// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// On Dragolangnfly argc/argv are passed in DI, not SP, so we can't use _rt0_amd64.
TEXT _rt0_amd64_dragolangnfly(SB),NOSPLIT,$-8
	LEAQ	8(DI), SI // argv
	MOVQ	0(DI), DI // argc
	JMP	runtime·rt0_golang(SB)

TEXT _rt0_amd64_dragolangnfly_lib(SB),NOSPLIT,$0
	JMP	_rt0_amd64_lib(SB)
