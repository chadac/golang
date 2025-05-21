// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build lldb

// Used by gcc_signal_darwin_arm64.c when doing the test build during cgolang.
// We hope that for real binaries the definition provided by Go will take precedence
// and the linker will drop this .o file altogether, which is why this definition
// is all by itself in its own file.
void __attribute__((weak)) xx_cgolang_panicmem(void) {}
