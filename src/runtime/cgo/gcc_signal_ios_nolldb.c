// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !lldb && ios && arm64

#include <stdint.h>

void darwin_arm_init_thread_exception_port() {}
void darwin_arm_init_mach_exception_handler() {}
