// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT	·getFP(SB), NOSPLIT|NOFRAME, $0-8
	MOVQ	BP, ret+0(FP)
	RET
