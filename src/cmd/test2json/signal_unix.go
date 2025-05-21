// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix || js || wasip1

package main

import (
	"os"
	"syscall"
)

var signalsToIgnore = []os.Signal{os.Interrupt, syscall.SIGQUIT}
