// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package cgolangsotest

// This test verifies that Go can access C variables
// in shared object file via cgolang.

/*
// intentionally write the same LDFLAGS differently
// to test correct handling of LDFLAGS.
#cgolang windows CFLAGS: -DIMPORT_DLL
#cgolang linux LDFLAGS: -L. -lcgolangsotest
#cgolang dragolangnfly LDFLAGS: -L. -l cgolangsotest
#cgolang freebsd LDFLAGS: -L. -l cgolangsotest
#cgolang openbsd LDFLAGS: -L. -l cgolangsotest
#cgolang solaris LDFLAGS: -L. -lcgolangsotest
#cgolang netbsd LDFLAGS: -L. libcgolangsotest.so
#cgolang darwin LDFLAGS: -L. libcgolangsotest.dylib
#cgolang windows LDFLAGS: -L. libcgolangsotest.a
#cgolang aix LDFLAGS: -L. -l cgolangsotest

#include "cgolangso_c.h"

const char* getVar() {
	    return exported_var;
}
*/
import "C"

import "fmt"

func Test() {
	const want = "Hello world"
	golangt := C.GoString(C.getVar())
	if golangt != want {
		panic(fmt.Sprintf("testExportedVar: golangt %q, but want %q", golangt, want))
	}
	golangt = C.GoString(C.exported_var)
	if golangt != want {
		panic(fmt.Sprintf("testExportedVar: golangt %q, but want %q", golangt, want))
	}
}
