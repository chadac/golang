// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"internal/testenv"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestInst tests that only one instantiation of Sort is created, even though generic
// Sort is used for multiple pointer types across two packages.
func TestInst(t *testing.T) {
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveGoRun(t)

	// Build ptrsort.golang, which uses package mysort.
	var output []byte
	var err error
	filename := "ptrsort.golang"
	exename := "ptrsort"
	outname := "ptrsort.out"
	golangtool := testenv.GoToolPath(t)
	dest := filepath.Join(t.TempDir(), exename)
	cmd := testenv.Command(t, golangtool, "build", "-o", dest, filepath.Join("testdata", filename))
	if output, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed: %v:\nOutput: %s\n", err, output)
	}

	// Test that there is exactly one shape-based instantiation of Sort in
	// the executable.
	cmd = testenv.Command(t, golangtool, "tool", "nm", dest)
	if output, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\n", err, output)
	}
	// Look for shape-based instantiation of Sort, but ignore any extra wrapper
	// ending in "-tramp" (which are created on riscv).
	re := regexp.MustCompile(`\bSort\[.*shape.*\][^-]`)
	r := re.FindAllIndex(output, -1)
	if len(r) != 1 {
		t.Fatalf("Wanted 1 instantiations of Sort function, golangt %d\n", len(r))
	}

	// Actually run the test and make sure output is correct.
	cmd = testenv.Command(t, golangtool, "run", filepath.Join("testdata", filename))
	if output, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\n", err, output)
	}
	out, err := os.ReadFile(filepath.Join("testdata", outname))
	if err != nil {
		t.Fatalf("Could not find %s\n", outname)
	}
	if string(out) != string(output) {
		t.Fatalf("Wanted output %v, golangt %v\n", string(out), string(output))
	}
}
