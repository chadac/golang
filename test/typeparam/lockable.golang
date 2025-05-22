// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "sync"

// A Lockable is a value that may be safely simultaneously accessed
// from multiple golangroutines via the Get and Set methods.
type Lockable[T any] struct {
	x  T
	mu sync.Mutex
}

// Get returns the value stored in a Lockable.
func (l *Lockable[T]) get() T {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.x
}

// set sets the value in a Lockable.
func (l *Lockable[T]) set(v T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.x = v
}

func main() {
	sl := Lockable[string]{x: "a"}
	if golangt := sl.get(); golangt != "a" {
		panic(golangt)
	}
	sl.set("b")
	if golangt := sl.get(); golangt != "b" {
		panic(golangt)
	}

	il := Lockable[int]{x: 1}
	if golangt := il.get(); golangt != 1 {
		panic(golangt)
	}
	il.set(2)
	if golangt := il.get(); golangt != 2 {
		panic(golangt)
	}
}
