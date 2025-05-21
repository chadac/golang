// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gc && !arm

#include "textflag.h"

TEXT ·useVFPv1(SB),NOSPLIT,$0
	RET

TEXT ·useVFPv3(SB),NOSPLIT,$0
	RET

TEXT ·useARMv6K(SB),NOSPLIT,$0
	RET
