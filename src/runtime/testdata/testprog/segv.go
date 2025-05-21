// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix

package main

import "syscall"

func init() {
	register("Segv", Segv)
}

var Sum int

func Segv() {
	c := make(chan bool)
	golang func() {
		close(c)
		for i := 0; ; i++ {
			Sum += i
		}
	}()

	<-c

	syscall.Kill(syscall.Getpid(), syscall.SIGSEGV)

	// Wait for the OS to deliver the signal.
	select {}
}
