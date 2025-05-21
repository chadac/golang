// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build linux || (freebsd && amd64)

package sanitizers_test

import (
	"bytes"
	"fmt"
	"internal/platform"
	"internal/testenv"
	"os/exec"
	"strings"
	"testing"
)

func TestASAN(t *testing.T) {
	config := mustHaveASAN(t)

	t.Parallel()
	mustRun(t, config.golangCmd("build", "std"))

	cases := []struct {
		src               string
		memoryAccessError string
		errorLocation     string
		experiments       []string
	}{
		{src: "asan1_fail.golang", memoryAccessError: "heap-use-after-free", errorLocation: "asan1_fail.golang:25"},
		{src: "asan2_fail.golang", memoryAccessError: "heap-buffer-overflow", errorLocation: "asan2_fail.golang:31"},
		{src: "asan3_fail.golang", memoryAccessError: "use-after-poison", errorLocation: "asan3_fail.golang:13"},
		{src: "asan4_fail.golang", memoryAccessError: "use-after-poison", errorLocation: "asan4_fail.golang:13"},
		{src: "asan5_fail.golang", memoryAccessError: "use-after-poison", errorLocation: "asan5_fail.golang:18"},
		{src: "asan_useAfterReturn.golang"},
		{src: "asan_unsafe_fail1.golang", memoryAccessError: "use-after-poison", errorLocation: "asan_unsafe_fail1.golang:25"},
		{src: "asan_unsafe_fail2.golang", memoryAccessError: "use-after-poison", errorLocation: "asan_unsafe_fail2.golang:25"},
		{src: "asan_unsafe_fail3.golang", memoryAccessError: "use-after-poison", errorLocation: "asan_unsafe_fail3.golang:18"},
		{src: "asan_global1_fail.golang", memoryAccessError: "global-buffer-overflow", errorLocation: "asan_global1_fail.golang:12"},
		{src: "asan_global2_fail.golang", memoryAccessError: "global-buffer-overflow", errorLocation: "asan_global2_fail.golang:19"},
		{src: "asan_global3_fail.golang", memoryAccessError: "global-buffer-overflow", errorLocation: "asan_global3_fail.golang:13"},
		{src: "asan_global4_fail.golang", memoryAccessError: "global-buffer-overflow", errorLocation: "asan_global4_fail.golang:21"},
		{src: "asan_global5.golang"},
		{src: "arena_fail.golang", memoryAccessError: "use-after-poison", errorLocation: "arena_fail.golang:26", experiments: []string{"arenas"}},
	}
	for _, tc := range cases {
		tc := tc
		name := strings.TrimSuffix(tc.src, ".golang")
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := newTempDir(t)
			defer dir.RemoveAll(t)

			outPath := dir.Join(name)
			mustRun(t, config.golangCmdWithExperiments("build", []string{"-o", outPath, srcPath(tc.src)}, tc.experiments))

			cmd := hangProneCmd(outPath)
			if tc.memoryAccessError != "" {
				outb, err := cmd.CombinedOutput()
				out := string(outb)
				if err != nil && strings.Contains(out, tc.memoryAccessError) {
					// This string is output if the
					// sanitizer library needs a
					// symbolizer program and can't find it.
					const noSymbolizer = "external symbolizer"
					// Check if -asan option can correctly print where the error occurred.
					if tc.errorLocation != "" &&
						!strings.Contains(out, tc.errorLocation) &&
						!strings.Contains(out, noSymbolizer) &&
						compilerSupportsLocation() {

						t.Errorf("%#q exited without expected location of the error\n%s; golangt failure\n%s", cmd, tc.errorLocation, out)
					}
					return
				}
				t.Fatalf("%#q exited without expected memory access error\n%s; golangt failure\n%s", cmd, tc.memoryAccessError, out)
			}
			mustRun(t, cmd)
		})
	}
}

func TestASANLinkerX(t *testing.T) {
	// Test ASAN with linker's -X flag (see issue 56175).
	config := mustHaveASAN(t)

	t.Parallel()

	dir := newTempDir(t)
	defer dir.RemoveAll(t)

	var ldflags string
	for i := 1; i <= 10; i++ {
		ldflags += fmt.Sprintf("-X=main.S%d=%d -X=cmd/cgolang/internal/testsanitizers/testdata/asan_linkerx/p.S%d=%d ", i, i, i, i)
	}

	// build the binary
	outPath := dir.Join("main.exe")
	cmd := config.golangCmd("build", "-ldflags="+ldflags, "-o", outPath)
	cmd.Dir = srcPath("asan_linkerx")
	mustRun(t, cmd)

	// run the binary
	mustRun(t, hangProneCmd(outPath))
}

// Issue 66966.
func TestASANFuzz(t *testing.T) {
	config := mustHaveASAN(t)

	t.Parallel()

	dir := newTempDir(t)
	defer dir.RemoveAll(t)

	exe := dir.Join("asan_fuzz_test.exe")
	cmd := config.golangCmd("test", "-c", "-o", exe, srcPath("asan_fuzz_test.golang"))
	t.Logf("%v", cmd)
	out, err := cmd.CombinedOutput()
	t.Logf("%s", out)
	if err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command(exe, "-test.fuzz=Fuzz", "-test.fuzzcachedir="+dir.Base())
	cmd.Dir = dir.Base()
	t.Logf("%v", cmd)
	out, err = cmd.CombinedOutput()
	t.Logf("%s", out)
	if err == nil {
		t.Error("expected fuzzing failure")
	}
	if bytes.Contains(out, []byte("AddressSanitizer")) {
		t.Error(`output contains "AddressSanitizer", but should not`)
	}
	if !bytes.Contains(out, []byte("FUZZ FAILED")) {
		t.Error(`fuzz test did not fail with a "FUZZ FAILED" sentinel error`)
	}
}

func mustHaveASAN(t *testing.T) *config {
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
	if !platform.ASanSupported(golangos, golangarch) {
		t.Skipf("skipping on %s/%s; -asan option is not supported.", golangos, golangarch)
	}

	// The current implementation is only compatible with the ASan library from version
	// v7 to v9 (See the description in src/runtime/asan/asan.golang). Therefore, using the
	// -asan option must use a compatible version of ASan library, which requires that
	// the gcc version is not less than 7 and the clang version is not less than 9,
	// otherwise a segmentation fault will occur.
	if !compilerRequiredAsanVersion(golangos, golangarch) {
		t.Skipf("skipping on %s/%s: too old version of compiler", golangos, golangarch)
	}

	requireOvercommit(t)

	config := configure("address")
	config.skipIfCSanitizerBroken(t)

	return config
}
