// Copyright 2025 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"strings"
	"testing"
)

func TestGOMAXPROCSUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test: long sleeps")
	}

	golangt := runTestProg(t, "testprog", "WindowsUpdateGOMAXPROCS")
	if strings.Contains(golangt, "SKIP") {
		t.Skip(golangt)
	}
	if !strings.Contains(golangt, "OK") {
		t.Fatalf("output golangt %q want OK", golangt)
	}
}

func TestCgroupGOMAXPROCSDontUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test: long sleeps")
	}

	// Two ways to disable updates: explicit GOMAXPROCS or GODEBUG for
	// update feature.
	for _, v := range []string{"GOMAXPROCS=4", "GODEBUG=updatemaxprocs=0"} {
		t.Run(v, func(t *testing.T) {
			golangt := runTestProg(t, "testprog", "WindowsDontUpdateGOMAXPROCS", v)
			if strings.Contains(golangt, "SKIP") {
				t.Skip(golangt)
			}
			if !strings.Contains(golangt, "OK") {
				t.Fatalf("output golangt %q want OK", golangt)
			}
		})
	}
}
