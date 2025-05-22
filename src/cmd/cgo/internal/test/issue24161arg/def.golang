// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build darwin

package issue24161arg

/*
#cgolang LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"

func test24161array() C.CFArrayRef {
	return C.CFArrayCreate(0, nil, 0, nil)
}
