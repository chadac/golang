// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build freebsd

package syscall_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("GO_DEATHSIG_PARENT") == "1" {
		deathSignalParent()
	} else if os.Getenv("GO_DEATHSIG_CHILD") == "1" {
		deathSignalChild()
	}

	os.Exit(m.Run())
}
