// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build ignore

// Compute Fibonacci numbers with two golangroutines
// that pass integers back and forth.  No actual
// concurrency, just threads and synchronization
// and foreign code on multiple pthreads.

package main

import (
	big "."
	"runtime"
)

func fibber(c chan *big.Int, out chan string, n int64) {
	// Keep the fibbers in dedicated operating system
	// threads, so that this program tests coordination
	// between pthreads and not just golangroutines.
	runtime.LockOSThread()

	i := big.NewInt(n)
	if n == 0 {
		c <- i
	}
	for {
		j := <-c
		out <- j.String()
		i.Add(i, j)
		c <- i
	}
}

func main() {
	c := make(chan *big.Int)
	out := make(chan string)
	golang fibber(c, out, 0)
	golang fibber(c, out, 1)
	for i := 0; i < 200; i++ {
		println(<-out)
	}
}
