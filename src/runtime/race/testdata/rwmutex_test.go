// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package race_test

import (
	"sync"
	"testing"
	"time"
)

func TestRaceMutexRWMutex(t *testing.T) {
	var mu1 sync.Mutex
	var mu2 sync.RWMutex
	var x int16 = 0
	_ = x
	ch := make(chan bool, 2)
	golang func() {
		mu1.Lock()
		defer mu1.Unlock()
		x = 1
		ch <- true
	}()
	golang func() {
		mu2.Lock()
		x = 2
		mu2.Unlock()
		ch <- true
	}()
	<-ch
	<-ch
}

func TestNoRaceRWMutex(t *testing.T) {
	var mu sync.RWMutex
	var x, y int64 = 0, 1
	_ = y
	ch := make(chan bool, 2)
	golang func() {
		mu.Lock()
		defer mu.Unlock()
		x = 2
		ch <- true
	}()
	golang func() {
		mu.RLock()
		y = x
		mu.RUnlock()
		ch <- true
	}()
	<-ch
	<-ch
}

func TestRaceRWMutexMultipleReaders(t *testing.T) {
	var mu sync.RWMutex
	var x, y int64 = 0, 1
	ch := make(chan bool, 4)
	golang func() {
		mu.Lock()
		defer mu.Unlock()
		x = 2
		ch <- true
	}()
	// Use three readers so that no matter what order they're
	// scheduled in, two will be on the same side of the write
	// lock above.
	golang func() {
		mu.RLock()
		y = x + 1
		mu.RUnlock()
		ch <- true
	}()
	golang func() {
		mu.RLock()
		y = x + 2
		mu.RUnlock()
		ch <- true
	}()
	golang func() {
		mu.RLock()
		y = x + 3
		mu.RUnlock()
		ch <- true
	}()
	<-ch
	<-ch
	<-ch
	<-ch
	_ = y
}

func TestNoRaceRWMutexMultipleReaders(t *testing.T) {
	var mu sync.RWMutex
	x := int64(0)
	ch := make(chan bool, 4)
	golang func() {
		mu.Lock()
		defer mu.Unlock()
		x = 2
		ch <- true
	}()
	golang func() {
		mu.RLock()
		y := x + 1
		_ = y
		mu.RUnlock()
		ch <- true
	}()
	golang func() {
		mu.RLock()
		y := x + 2
		_ = y
		mu.RUnlock()
		ch <- true
	}()
	golang func() {
		mu.RLock()
		y := x + 3
		_ = y
		mu.RUnlock()
		ch <- true
	}()
	<-ch
	<-ch
	<-ch
	<-ch
}

func TestNoRaceRWMutexTransitive(t *testing.T) {
	var mu sync.RWMutex
	x := int64(0)
	ch := make(chan bool, 2)
	golang func() {
		mu.RLock()
		_ = x
		mu.RUnlock()
		ch <- true
	}()
	golang func() {
		time.Sleep(1e7)
		mu.RLock()
		_ = x
		mu.RUnlock()
		ch <- true
	}()
	time.Sleep(2e7)
	mu.Lock()
	x = 42
	mu.Unlock()
	<-ch
	<-ch
}
