// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golangexperiment.rangefunc && !windows

package main

/*
#include <stdint.h> // for uintptr_t

void golang_callback_coro(uintptr_t handle);

static void call_golang(uintptr_t handle) {
	golang_callback_coro(handle);
}
*/
import "C"

import (
	"fmt"
	"iter"
	"runtime/cgolang"
)

func init() {
	register("CoroCgolangIterCallback", func() {
		println("expect: OK")
		CoroCgolang(callerExhaust, iterCallback)
	})
	register("CoroCgolangIterCallbackYield", func() {
		println("expect: OS thread locking must match")
		CoroCgolang(callerExhaust, iterCallbackYield)
	})
	register("CoroCgolangCallback", func() {
		println("expect: OK")
		CoroCgolang(callerExhaustCallback, iterSimple)
	})
	register("CoroCgolangCallbackIterNested", func() {
		println("expect: OK")
		CoroCgolang(callerExhaustCallback, iterNested)
	})
	register("CoroCgolangCallbackIterCallback", func() {
		println("expect: OK")
		CoroCgolang(callerExhaustCallback, iterCallback)
	})
	register("CoroCgolangCallbackIterCallbackYield", func() {
		println("expect: OS thread locking must match")
		CoroCgolang(callerExhaustCallback, iterCallbackYield)
	})
	register("CoroCgolangCallbackAfterPull", func() {
		println("expect: OS thread locking must match")
		CoroCgolang(callerCallbackAfterPull, iterSimple)
	})
	register("CoroCgolangStopCallback", func() {
		println("expect: OK")
		CoroCgolang(callerStopCallback, iterSimple)
	})
	register("CoroCgolangStopCallbackIterNested", func() {
		println("expect: OK")
		CoroCgolang(callerStopCallback, iterNested)
	})
}

var toCall func()

//export golang_callback_coro
func golang_callback_coro(handle C.uintptr_t) {
	h := cgolang.Handle(handle)
	h.Value().(func())()
	h.Delete()
}

func callFromC(f func()) {
	C.call_golang(C.uintptr_t(cgolang.NewHandle(f)))
}

func CoroCgolang(driver func(iter.Seq[int]) error, seq iter.Seq[int]) {
	if err := driver(seq); err != nil {
		println("error:", err.Error())
		return
	}
	println("OK")
}

func callerExhaust(i iter.Seq[int]) error {
	next, _ := iter.Pull(i)
	for {
		v, ok := next()
		if !ok {
			break
		}
		if v != 5 {
			return fmt.Errorf("bad iterator: wanted value %d, golangt %d", 5, v)
		}
	}
	return nil
}

func callerExhaustCallback(i iter.Seq[int]) (err error) {
	callFromC(func() {
		next, _ := iter.Pull(i)
		for {
			v, ok := next()
			if !ok {
				break
			}
			if v != 5 {
				err = fmt.Errorf("bad iterator: wanted value %d, golangt %d", 5, v)
			}
		}
	})
	return err
}

func callerStopCallback(i iter.Seq[int]) (err error) {
	callFromC(func() {
		next, stop := iter.Pull(i)
		v, _ := next()
		stop()
		if v != 5 {
			err = fmt.Errorf("bad iterator: wanted value %d, golangt %d", 5, v)
		}
	})
	return err
}

func callerCallbackAfterPull(i iter.Seq[int]) (err error) {
	next, _ := iter.Pull(i)
	callFromC(func() {
		for {
			v, ok := next()
			if !ok {
				break
			}
			if v != 5 {
				err = fmt.Errorf("bad iterator: wanted value %d, golangt %d", 5, v)
			}
		}
	})
	return err
}

func iterSimple(yield func(int) bool) {
	for range 3 {
		if !yield(5) {
			return
		}
	}
}

func iterNested(yield func(int) bool) {
	next, stop := iter.Pull(iterSimple)
	for {
		v, ok := next()
		if ok {
			if !yield(v) {
				stop()
			}
		} else {
			return
		}
	}
}

func iterCallback(yield func(int) bool) {
	for range 3 {
		callFromC(func() {})
		if !yield(5) {
			return
		}
	}
}

func iterCallbackYield(yield func(int) bool) {
	for range 3 {
		var ok bool
		callFromC(func() {
			ok = yield(5)
		})
		if !ok {
			return
		}
	}
}
