// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include <process.h>
#include <stdlib.h>
#include <stdio.h>

void golangpanic(void);

static unsigned int __attribute__((__stdcall__))
die(void* x)
{
	golangpanic();
	return 0;
}

void
start(void)
{
	if(_beginthreadex(0, 0, die, 0, 0, 0) != 0)
		printf("_beginthreadex failed\n");
}
