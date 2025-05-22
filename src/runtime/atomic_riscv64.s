// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// func publicationBarrier()
TEXT ·publicationBarrier(SB),NOSPLIT|NOFRAME,$0-0
	FENCE
	RET
