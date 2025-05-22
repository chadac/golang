// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build js

package runtime

// resetMemoryDataView signals the JS front-end that WebAssembly's memory.grow instruction has been used.
// This allows the front-end to replace the old DataView object with a new one.
//
//golang:wasmimport golangjs runtime.resetMemoryDataView
func resetMemoryDataView()
