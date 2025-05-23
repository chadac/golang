// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// These functions are called from C code via cgolang/callbacks.golang.

// Panic.

func _cgolang_panic_internal(p *byte) {
	panic(golangstringnocopy(p))
}
