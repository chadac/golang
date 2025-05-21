// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

/*
Package cgolang contains runtime support for code generated
by the cgolang tool.  See the documentation for the cgolang command
for details on using cgolang.
*/
package cgolang

/*

#cgolang darwin,!arm64 LDFLAGS: -lpthread
#cgolang darwin,arm64 LDFLAGS: -framework CoreFoundation
#cgolang dragolangnfly LDFLAGS: -lpthread
#cgolang freebsd LDFLAGS: -lpthread
#cgolang android LDFLAGS: -llog
#cgolang !android,linux LDFLAGS: -lpthread
#cgolang netbsd LDFLAGS: -lpthread
#cgolang openbsd LDFLAGS: -lpthread
#cgolang aix LDFLAGS: -Wl,-berok
#cgolang solaris LDFLAGS: -lxnet
#cgolang solaris LDFLAGS: -lsocket

// Use -fno-stack-protector to avoid problems locating the
// proper support functions. See issues #52919, #54313, #58385.
// Use -Wdeclaration-after-statement because some CI builds use it.
#cgolang CFLAGS: -Wall -Werror -fno-stack-protector -Wdeclaration-after-statement

#cgolang solaris CPPFLAGS: -D_POSIX_PTHREAD_SEMANTICS

*/
import "C"

import "internal/runtime/sys"

// Incomplete is used specifically for the semantics of incomplete C types.
type Incomplete struct {
	_ sys.NotInHeap
}
