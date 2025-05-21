// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import _ "unsafe" // for golang:linkname

//golang:linkname boring_runtime_arg0 crypto/internal/boring.runtime_arg0
func boring_runtime_arg0() string {
	// On Windows, argslice is not set, and it's too much work to find argv0.
	if len(argslice) == 0 {
		return ""
	}
	return argslice[0]
}
