// compile

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

type A []int;
type M map[int] int;

func main() {
	var a *A = &A{0};
	var m *M = &M{0 : 0};  // should be legal to use & here for consistency with other composite constructors (prev. line)
	_, _ = a, m;
}

/*
uetli:~/Source/golang1/test/bugs gri$ 6g bug098.golang && 6l bug098.6 && 6.out
bug098.golang:10: illegal types for operand: AS
	(*MAP[<int32>INT32]<int32>INT32)
	(**MAP[<int32>INT32]<int32>INT32)
*/
