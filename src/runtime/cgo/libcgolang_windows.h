// Copyright 2020 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Call _beginthread, aborting on failure.
void _cgolang_beginthread(unsigned long (__stdcall *func)(void*), void* arg);
