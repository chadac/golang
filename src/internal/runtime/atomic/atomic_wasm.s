// Copyright 2020 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·StorepNoWB(SB), NOSPLIT, $0-16
	MOVD ptr+0(FP), R0
	MOVD val+8(FP), 0(R0)
	RET
