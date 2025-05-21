// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !aix && !darwin && !dragolangnfly && !freebsd && !linux && !netbsd && !openbsd && !solaris && !windows && !zos

package nettest

func supportsRawSocket() bool {
	return false
}
