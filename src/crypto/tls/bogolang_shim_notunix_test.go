// Copyright 2025 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !unix || wasm

package tls

func pauseProcess() {
	panic("-wait-for-debugger not supported on this OS")
}
