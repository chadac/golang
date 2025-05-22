// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·publicationBarrier(SB),NOSPLIT|NOFRAME,$0-0
	DMB	$0xe	// DMB ST
	RET
