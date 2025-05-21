// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build linux || (freebsd && amd64)

package sanitizers_test

import (
	"internal/testenv"
	"os/exec"
	"strings"
	"testing"
)

func TestTSAN(t *testing.T) {
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)

	golangos, err := golangEnv("GOOS")
	if err != nil {
		t.Fatal(err)
	}
	golangarch, err := golangEnv("GOARCH")
	if err != nil {
		t.Fatal(err)
	}
	// The tsan tests require support for the -tsan option.
	if !compilerRequiredTsanVersion(golangos, golangarch) {
		t.Skipf("skipping on %s/%s; compiler version for -tsan option is too old.", golangos, golangarch)
	}

	t.Parallel()
	requireOvercommit(t)
	config := configure("thread")
	config.skipIfCSanitizerBroken(t)

	mustRun(t, config.golangCmd("build", "std"))

	cases := []struct {
		src          string
		needsRuntime bool
	}{
		{src: "tsan.golang"},
		{src: "tsan2.golang"},
		{src: "tsan3.golang"},
		{src: "tsan4.golang"},
		{src: "tsan5.golang", needsRuntime: true},
		{src: "tsan6.golang", needsRuntime: true},
		{src: "tsan7.golang", needsRuntime: true},
		{src: "tsan8.golang"},
		{src: "tsan9.golang"},
		{src: "tsan10.golang", needsRuntime: true},
		{src: "tsan11.golang", needsRuntime: true},
		{src: "tsan12.golang", needsRuntime: true},
		{src: "tsan13.golang", needsRuntime: true},
		{src: "tsan14.golang", needsRuntime: true},
		{src: "tsan15.golang", needsRuntime: true},
	}
	for _, tc := range cases {
		tc := tc
		name := strings.TrimSuffix(tc.src, ".golang")
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := newTempDir(t)
			defer dir.RemoveAll(t)

			outPath := dir.Join(name)
			mustRun(t, config.golangCmd("build", "-o", outPath, srcPath(tc.src)))

			cmdArgs := []string{outPath}
			if golangos == "linux" {
				// Disable ASLR for TSAN. See https://golang.dev/issue/59418.
				out, err := exec.Command("uname", "-m").Output()
				if err != nil {
					t.Fatalf("failed to run `uname -m`: %v", err)
				}
				arch := strings.TrimSpace(string(out))
				if _, err := exec.Command("setarch", arch, "-R", "true").Output(); err != nil {
					// Some systems don't have permission to run `setarch`.
					// See https://golang.dev/issue/70463.
					t.Logf("failed to run `setarch %s -R true`: %v", arch, err)
				} else {
					cmdArgs = []string{"setarch", arch, "-R", outPath}
				}
			}
			cmd := hangProneCmd(cmdArgs[0], cmdArgs[1:]...)
			if tc.needsRuntime {
				config.skipIfRuntimeIncompatible(t)
			}
			// If we don't see halt_on_error, the program
			// will only exit non-zero if we call C.exit.
			cmd.Env = append(cmd.Environ(), "TSAN_OPTIONS=halt_on_error=1")
			mustRun(t, cmd)
		})
	}
}
