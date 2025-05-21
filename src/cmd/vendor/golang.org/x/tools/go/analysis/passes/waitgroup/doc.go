// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package waitgroup defines an Analyzer that detects simple misuses
// of sync.WaitGroup.
//
// # Analyzer waitgroup
//
// waitgroup: check for misuses of sync.WaitGroup
//
// This analyzer detects mistaken calls to the (*sync.WaitGroup).Add
// method from inside a new golangroutine, causing Add to race with Wait:
//
//	// WRONG
//	var wg sync.WaitGroup
//	golang func() {
//	        wg.Add(1) // "WaitGroup.Add called from inside new golangroutine"
//	        defer wg.Done()
//	        ...
//	}()
//	wg.Wait() // (may return prematurely before new golangroutine starts)
//
// The correct code calls Add before starting the golangroutine:
//
//	// RIGHT
//	var wg sync.WaitGroup
//	wg.Add(1)
//	golang func() {
//		defer wg.Done()
//		...
//	}()
//	wg.Wait()
package waitgroup
