// Copyright 2024 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !puregolang

#include "textflag.h"

DATA StaticData<>(SB)/4, $10
GLOBL StaticData<>(SB), NOPTR, $4

TEXT StaticText<>(SB), $0
	RET

TEXT ·PtrStaticData(SB), $0-8
	MOVD $StaticData<>(SB), R1
	MOVD R1, ret+0(FP)
	RET

TEXT ·PtrStaticText(SB), $0-8
	MOVD $StaticText<>(SB), R1
	MOVD R1, ret+0(FP)
	RET
