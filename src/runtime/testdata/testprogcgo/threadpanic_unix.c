// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!windows

#include <stdlib.h>
#include <stdio.h>
#include <pthread.h>

void golangpanic(void);

static void*
die(void* x)
{
	golangpanic();
	return 0;
}

void
start(void)
{
	pthread_t t;
	if(pthread_create(&t, 0, die, 0) != 0)
		printf("pthread_create failed\n");
}
