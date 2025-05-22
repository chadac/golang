// compile

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 55242: golangfrontend crash calling function that returns
// trailing empty struct.

package p

func F1() (int, struct{}) {
	return 0, struct{}{}
}

func F2() {
	F1()
}
