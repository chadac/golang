// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	_ "unsafe"
)

//golang:linkname runtime_ignoreHangup internal/poll.runtime_ignoreHangup
func runtime_ignoreHangup() {
	getg().m.ignoreHangup = true
}

//golang:linkname runtime_unignoreHangup internal/poll.runtime_unignoreHangup
func runtime_unignoreHangup(sig string) {
	getg().m.ignoreHangup = false
}

func ignoredNote(note *byte) bool {
	if note == nil {
		return false
	}
	if golangstringnocopy(note) != "hangup" {
		return false
	}
	return getg().m.ignoreHangup
}
