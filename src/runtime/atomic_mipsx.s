// Copyright 2016 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build mips || mipsle

#include "textflag.h"

TEXT ·publicationBarrier(SB),NOSPLIT,$0
	SYNC
	RET
