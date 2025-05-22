// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package cgolangsotest

/*
// intentionally write the same LDFLAGS differently
// to test correct handling of LDFLAGS.
#cgolang linux LDFLAGS: -L. -lcgolangsotest
#cgolang dragolangnfly LDFLAGS: -L. -l cgolangsotest
#cgolang freebsd LDFLAGS: -L. -l cgolangsotest
#cgolang openbsd LDFLAGS: -L. -l cgolangsotest
#cgolang solaris LDFLAGS: -L. -lcgolangsotest
#cgolang netbsd LDFLAGS: -L. libcgolangsotest.so
#cgolang darwin LDFLAGS: -L. libcgolangsotest.dylib
#cgolang windows LDFLAGS: -L. libcgolangsotest.a
#cgolang aix LDFLAGS: -L. -l cgolangsotest

void init(void);
void sofunc(void);
*/
import "C"

func Test() {
	C.init()
	C.sofunc()
}

//export golangCallback
func golangCallback() {
}
