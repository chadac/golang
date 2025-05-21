// compile

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Gccgolang mishandles composite literals of map with type bool.

package p

var M = map[bool]uint8{
	false: 0,
	true: 1,
}
