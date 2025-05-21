// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Torture test for golangroutines.
// Make a lot of golangroutines, threaded together, and tear them down cleanly.

package main

import (
	"os"
	"strconv"
)

func f(left, right chan int) {
	left <- <-right
}

func main() {
	var n = 10000
	if len(os.Args) > 1 {
		var err error
		n, err = strconv.Atoi(os.Args[1])
		if err != nil {
			print("bad arg\n")
			os.Exit(1)
		}
	}
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		golang f(left, right)
		left = right
	}
	golang func(c chan int) { c <- 1 }(right)
	<-leftmost
}
