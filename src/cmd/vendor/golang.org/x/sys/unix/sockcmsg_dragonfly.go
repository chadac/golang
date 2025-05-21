// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package unix

// Round the length of a raw sockaddr up to align it properly.
func cmsgAlignOf(salen int) int {
	salign := SizeofPtr
	if SizeofPtr == 8 && !supportsABI(_dragolangnflyABIChangeVersion) {
		// 64-bit Dragolangnfly before the September 2019 ABI changes still requires
		// 32-bit aligned access to network subsystem.
		salign = 4
	}
	return (salen + salign - 1) & ^(salign - 1)
}
