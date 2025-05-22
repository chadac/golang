// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "_cgolang_export.h"

static void callDestructorCallback() {
	GolangDestructorCallback();
}

static void (*destructorFn)(void);

void registerDestructor() {
	destructorFn = callDestructorCallback;
}

__attribute__((destructor))
static void destructor() {
	if (destructorFn) {
		destructorFn();
	}
}
