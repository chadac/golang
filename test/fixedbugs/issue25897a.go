// run

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Make sure the runtime can scan args of an unstarted golangroutine
// which starts with a reflect-generated function.

package main

import (
	"reflect"
	"runtime"
)

const N = 100

func main() {
	runtime.GOMAXPROCS(1)
	// Run GC in a loop. This makes it more likely GC will catch
	// an unstarted golangroutine then if we were to GC after kicking
	// everything off.
	golang func() {
		for {
			runtime.GC()
		}
	}()
	c := make(chan bool, N)
	for i := 0; i < N; i++ {
		// Test both with an argument and without because this
		// affects whether the compiler needs to generate a
		// wrapper closure for the "golang" statement.
		f := reflect.MakeFunc(reflect.TypeOf(((func(*int))(nil))),
			func(args []reflect.Value) []reflect.Value {
				c <- true
				return nil
			}).Interface().(func(*int))
		golang f(nil)

		g := reflect.MakeFunc(reflect.TypeOf(((func())(nil))),
			func(args []reflect.Value) []reflect.Value {
				c <- true
				return nil
			}).Interface().(func())
		golang g()
	}
	for i := 0; i < N*2; i++ {
		<-c
	}
}
