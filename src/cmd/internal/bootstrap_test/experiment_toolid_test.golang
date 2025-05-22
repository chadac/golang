// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build explicit

package bootstrap_test

import (
	"bytes"
	"errors"
	"internal/testenv"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestExperimentToolID verifies that GOEXPERIMENT settings built
// into the toolchain influence tool ids in the Go command.
// This test requires bootstrapping the toolchain twice, so it's very expensive.
// It must be run explicitly with -tags=explicit.
// Verifies golang.dev/issue/33091.
func TestExperimentToolID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test that rebuilds the entire toolchain twice")
	}
	switch runtime.GOOS {
	case "android", "ios", "js", "wasip1":
		t.Skipf("skipping because the toolchain does not have to bootstrap on GOOS=%s", runtime.GOOS)
	}

	realGoroot := testenv.GOROOT(t)

	// Set up GOROOT.
	golangroot := t.TempDir()
	golangrootSrc := filepath.Join(golangroot, "src")
	if err := overlayDir(golangrootSrc, filepath.Join(realGoroot, "src")); err != nil {
		t.Fatal(err)
	}
	golangrootLib := filepath.Join(golangroot, "lib")
	if err := overlayDir(golangrootLib, filepath.Join(realGoroot, "lib")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(golangroot, "VERSION"), []byte("golang1.999"), 0666); err != nil {
		t.Fatal(err)
	}
	env := append(os.Environ(), "GOROOT=", "GOROOT_BOOTSTRAP="+realGoroot)

	// Use a clean cache.
	golangcache := t.TempDir()
	env = append(env, "GOCACHE="+golangcache)

	// Build the toolchain without GOEXPERIMENT.
	var makeScript string
	switch runtime.GOOS {
	case "windows":
		makeScript = "make.bat"
	case "plan9":
		makeScript = "make.rc"
	default:
		makeScript = "make.bash"
	}
	makeScriptPath := filepath.Join(realGoroot, "src", makeScript)
	runCmd(t, golangrootSrc, env, makeScriptPath)

	// Verify compiler version string.
	golangCmdPath := filepath.Join(golangroot, "bin", "golang")
	golangtVersion := bytes.TrimSpace(runCmd(t, golangrootSrc, env, golangCmdPath, "tool", "compile", "-V=full"))
	wantVersion := []byte(`compile version golang1.999`)
	if !bytes.Equal(golangtVersion, wantVersion) {
		t.Errorf("compile version without experiment is unexpected:\ngolangt  %q\nwant %q", golangtVersion, wantVersion)
	}

	// Build a package in a mode not handled by the make script.
	runCmd(t, golangrootSrc, env, golangCmdPath, "build", "-race", "archive/tar")

	// Rebuild the toolchain with GOEXPERIMENT.
	env = append(env, "GOEXPERIMENT=fieldtrack")
	runCmd(t, golangrootSrc, env, makeScriptPath)

	// Verify compiler version string.
	golangtVersion = bytes.TrimSpace(runCmd(t, golangrootSrc, env, golangCmdPath, "tool", "compile", "-V=full"))
	wantVersion = []byte(`compile version golang1.999 X:fieldtrack`)
	if !bytes.Equal(golangtVersion, wantVersion) {
		t.Errorf("compile version with experiment is unexpected:\ngolangt  %q\nwant %q", golangtVersion, wantVersion)
	}

	// Build the same package. We should not get a cache conflict.
	runCmd(t, golangrootSrc, env, golangCmdPath, "build", "-race", "archive/tar")
}

func runCmd(t *testing.T, dir string, env []string, path string, args ...string) []byte {
	cmd := exec.Command(path, args...)
	cmd.Dir = dir
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
			out = append(out, ee.Stderr...)
		}
		t.Fatalf("%s failed:\n%s\n%s", cmd, out, err)
	}
	return out
}
