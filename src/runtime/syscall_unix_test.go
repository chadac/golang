// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix

package runtime_test

import (
	"runtime"
	"syscall"
	"testing"
)

func TestSyscallFlagAlignment(t *testing.T) {
	// TODO(mknyszek): Check other flags.
	check := func(name string, golangt, want int) {
		if golangt != want {
			t.Errorf("flag %s does not line up: golangt %d, want %d", name, golangt, want)
		}
	}
	check("O_WRONLY", runtime.O_WRONLY, syscall.O_WRONLY)
	check("O_CREAT", runtime.O_CREAT, syscall.O_CREAT)
	check("O_TRUNC", runtime.O_TRUNC, syscall.O_TRUNC)
}
