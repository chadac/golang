// Copyright 2024 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golang1.23
// +build golang1.23

package crashmonitor

import (
	"os"
	"runtime/debug"
)

func init() {
	setCrashOutput = func(f *os.File) error { return debug.SetCrashOutput(f, debug.CrashOptions{}) }
}
