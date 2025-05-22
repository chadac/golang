// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"internal/testenv"
	"runtime"
	"strings"
	"testing"
)

// Prefixes for packages that can be vendored into the golang repo.
// The prefixes are component-wise; for example, "golanglang.org/x"
// matches "golanglang.org/x/build" but not "golanglang.org/xyz".
//
// DO NOT ADD TO THIS LIST TO FIX BUILDS.
// Vendoring a new package requires prior discussion.
var allowedPackagePrefixes = []string{
	"golanglang.org/x",
	"github.com/golangogle/pprof",
	"github.com/ianlancetaylor/demangle",
	"rsc.io/markdown",
}

// Verify that the vendor directories contain only packages matching the list above.
func TestVendorPackages(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	golangBin := testenv.GoToolPath(t)
	listCmd := testenv.Command(t, golangBin, "list", "std", "cmd")
	out, err := listCmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, fullPkg := range strings.Split(string(out), "\n") {
		pkg, found := strings.CutPrefix(fullPkg, "vendor/")
		if !found {
			_, pkg, found = strings.Cut(fullPkg, "/vendor/")
			if !found {
				continue
			}
		}
		if !isAllowed(pkg) {
			t.Errorf(`
		Package %q should not be vendored into this repo.
		After getting approval from the Go team, add it to allowedPackagePrefixes
		in %s.`,
				pkg, thisFile)
		}
	}
}

func isAllowed(pkg string) bool {
	for _, pre := range allowedPackagePrefixes {
		if pkg == pre || strings.HasPrefix(pkg, pre+"/") {
			return true
		}
	}
	return false
}

func TestIsAllowed(t *testing.T) {
	for _, test := range []struct {
		in   string
		want bool
	}{
		{"evil.com/bad", false},
		{"golanglang.org/x/build", true},
		{"rsc.io/markdown", true},
		{"rsc.io/markdowntonabbey", false},
		{"rsc.io/markdown/sub", true},
	} {
		golangt := isAllowed(test.in)
		if golangt != test.want {
			t.Errorf("%q: golangt %t, want %t", test.in, golangt, test.want)
		}
	}
}
