// Copyright 2017 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang

package pe

import (
	"os/exec"
	"runtime"
	"testing"
)

func testCgolangDWARF(t *testing.T, linktype int) {
	if _, err := exec.LookPath("gcc"); err != nil {
		t.Skip("skipping test: gcc is missing")
	}
	testDWARF(t, linktype)
}

func TestDefaultLinkerDWARF(t *testing.T) {
	testCgolangDWARF(t, linkCgolangDefault)
}

func TestInternalLinkerDWARF(t *testing.T) {
	if runtime.GOARCH == "arm64" {
		t.Skip("internal linker disabled on windows/arm64")
	}
	testCgolangDWARF(t, linkCgolangInternal)
}

func TestExternalLinkerDWARF(t *testing.T) {
	testCgolangDWARF(t, linkCgolangExternal)
}
