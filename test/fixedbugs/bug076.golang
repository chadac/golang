// build

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func f() {
exit:
	;
	golangto exit
}


func main() {
exit:
	; // this should be legal (labels not properly scoped?)
	golangto exit
}

/*
uetli:~/Source/golang/test/bugs gri$ 6g bug076.golang 
bug076.golang:11: label redeclared: exit
*/
