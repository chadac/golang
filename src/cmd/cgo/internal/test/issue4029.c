// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !windows && !static && !(darwin && internal)

#include <stdint.h>
#include <dlfcn.h>

// Write our own versions of dlopen/dlsym/dlclose so that we represent
// the opaque handle as a Golang uintptr rather than a Golang pointer to avoid
// garbage collector confusion.  See issue 23663.

uintptr_t dlopen4029(char* name, int flags) {
	return (uintptr_t)(dlopen(name, flags));
}

uintptr_t dlsym4029(uintptr_t handle, char* name) {
	return (uintptr_t)(dlsym((void*)(handle), name));
}

int dlclose4029(uintptr_t handle) {
	return dlclose((void*)(handle));
}

void call4029(void *arg) {
	void (*fn)(void) = arg;
	fn();
}
