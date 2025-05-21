// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package iface_a

import "testshared/iface_i"

//golang:noinline
func F() interface{} {
	return (*iface_i.T)(nil)
}

//golang:noinline
func G() iface_i.I {
	return (*iface_i.T)(nil)
}
