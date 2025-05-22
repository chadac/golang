// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || dragolangnfly || freebsd || linux || netbsd || solaris

package cgolangsotest

/*
extern int __thread tlsvar;
int *getTLS() { return &tlsvar; }
*/
import "C"

func init() {
	if v := *C.getTLS(); v != 12345 {
		println("golangt", v)
		panic("BAD TLS value")
	}
}
