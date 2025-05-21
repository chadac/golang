// compile

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// PR61264: IncDec statements involving composite literals caused in ICE in gccgolang.

package main

func main() {
        map[int]int{}[0]++
}
