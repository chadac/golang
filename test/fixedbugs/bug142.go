// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func panic1(s string) bool {
	panic(s);
}

func main() {
	x := false && panic1("first") && panic1("second");
	x = x == true && panic1("first") && panic1("second");
}

/*
; 6.out
second
panic PC=0x250f98
main·panic1+0x36 /Users/rsc/golangX/test/bugs/bug142.golang:6
	main·panic1(0xae30, 0x0)
main·main+0x23 /Users/rsc/golangX/test/bugs/bug142.golang:10
	main·main()
mainstart+0xf /Users/rsc/golangX/src/runtime/amd64/asm.s:53
	mainstart()
sys·Goexit /Users/rsc/golangX/src/runtime/proc.c:124
	sys·Goexit()
; 
*/
