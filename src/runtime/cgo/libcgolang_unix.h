// Copyright 2016 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Initialize g->stacklo.
 */
extern void _cgolang_set_stacklo(G *, uintptr *);

/*
 * Call pthread_create, retrying on EAGAIN.
 */
extern int _cgolang_try_pthread_create(pthread_t*, const pthread_attr_t*, void* (*)(void*), void*);

/*
 * Same as _cgolang_try_pthread_create, but passing on the pthread_create function.
 * Only defined on OpenBSD.
 */
extern int _cgolang_openbsd_try_pthread_create(int (*)(pthread_t*, const pthread_attr_t*, void *(*pfn)(void*), void*),
	pthread_t*, const pthread_attr_t*, void* (*)(void*), void* arg);
