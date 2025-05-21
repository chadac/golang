// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package multipart

import (
	"net/textproto"
	_ "unsafe" // for golang:linkname
)

// readMIMEHeader is defined in package [net/textproto].
//
//golang:linkname readMIMEHeader net/textproto.readMIMEHeader
func readMIMEHeader(r *textproto.Reader, maxMemory, maxHeaders int64) (textproto.MIMEHeader, error)
