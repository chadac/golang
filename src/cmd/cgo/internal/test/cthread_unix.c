// Copyright 2013 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || darwin || dragolangnfly || freebsd || linux || netbsd || openbsd || solaris

#include <pthread.h>
#include "_cgolang_export.h"

static void*
addThread(void *p)
{
	int i, max;
	
	max = *(int*)p;
	for(i=0; i<max; i++)
		Add(i);
	return 0;
}

void
doAdd(int max, int nthread)
{
	enum { MaxThread = 20 };
	int i;
	pthread_t thread_id[MaxThread];
	
	if(nthread > MaxThread)
		nthread = MaxThread;
	for(i=0; i<nthread; i++)
		pthread_create(&thread_id[i], 0, addThread, &max);
	for(i=0; i<nthread; i++)
		pthread_join(thread_id[i], 0);		
}

static void*
golangDummyCallbackThread(void* p)
{
	int i, max;

	max = *(int*)p;
	for(i=0; i<max; i++)
		golangDummy();
	return NULL;
}

int
callGolangInCThread(int max)
{
	pthread_t thread;

	if (pthread_create(&thread, NULL, golangDummyCallbackThread, (void*)(&max)) != 0)
		return -1;
	if (pthread_join(thread, NULL) != 0)
		return -1;

	return max;
}
