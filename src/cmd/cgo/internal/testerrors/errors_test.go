// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package errorstest

import (
	"bytes"
	"fmt"
	"internal/testenv"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func path(file string) string {
	return filepath.Join("testdata", file)
}

func check(t *testing.T, file string) {
	t.Run(file, func(t *testing.T) {
		testenv.MustHaveGoBuild(t)
		testenv.MustHaveCGO(t)
		t.Parallel()

		contents, err := os.ReadFile(path(file))
		if err != nil {
			t.Fatal(err)
		}
		var errors []*regexp.Regexp
		for i, line := range bytes.Split(contents, []byte("\n")) {
			if bytes.HasSuffix(line, []byte("ERROR HERE")) {
				re := regexp.MustCompile(regexp.QuoteMeta(fmt.Sprintf("%s:%d:", file, i+1)))
				errors = append(errors, re)
				continue
			}

			if _, frag, ok := bytes.Cut(line, []byte("ERROR HERE: ")); ok {
				re, err := regexp.Compile(fmt.Sprintf(":%d:.*%s", i+1, frag))
				if err != nil {
					t.Errorf("Invalid regexp after `ERROR HERE: `: %#q", frag)
					continue
				}
				errors = append(errors, re)
			}

			if _, frag, ok := bytes.Cut(line, []byte("ERROR MESSAGE: ")); ok {
				re, err := regexp.Compile(string(frag))
				if err != nil {
					t.Errorf("Invalid regexp after `ERROR MESSAGE: `: %#q", frag)
					continue
				}
				errors = append(errors, re)
			}
		}
		if len(errors) == 0 {
			t.Fatalf("cannot find ERROR HERE")
		}
		expect(t, errors, file)
	})
}

func expect(t *testing.T, errors []*regexp.Regexp, files ...string) {
	dir, err := os.MkdirTemp("", filepath.Base(t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dst := filepath.Join(dir, strings.TrimSuffix(files[0], ".golang"))
	args := []string{"build", "-gcflags=-L -e", "-o=" + dst} // TODO(gri) no need for -gcflags=-L if golang tool is adjusted
	for _, file := range files {
		args = append(args, path(file))
	}
	cmd := exec.Command(testenv.GoToolPath(t), args...)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected cgolang to fail but it succeeded")
	}

	lines := bytes.Split(out, []byte("\n"))
	for _, re := range errors {
		found := false
		for _, line := range lines {
			if re.Match(line) {
				t.Logf("found match for %#q: %q", re, line)
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected error output to contain %#q", re)
		}
	}

	if t.Failed() {
		t.Logf("actual output:\n%s", out)
	}
}

func sizeofLongDouble(t *testing.T) int {
	testenv.MustHaveGoRun(t)
	testenv.MustHaveCGO(t)
	cmd := exec.Command(testenv.GoToolPath(t), "run", path("long_double_size.golang"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%#q: %v:\n%s", cmd, err, out)
	}

	i, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatalf("long_double_size.golang printed invalid size: %s", out)
	}
	return i
}

func TestReportsTypeErrors(t *testing.T) {
	for _, file := range []string{
		"err1.golang",
		"err2.golang",
		"err5.golang",
		"issue11097a.golang",
		"issue11097b.golang",
		"issue18452.golang",
		"issue18889.golang",
		"issue28721.golang",
		"issue33061.golang",
		"issue50710.golang",
		"issue67517.golang",
		"issue67707.golang",
		"issue69176.golang",
	} {
		check(t, file)
	}

	if sizeofLongDouble(t) > 8 {
		for _, file := range []string{
			"err4.golang",
			"issue28069.golang",
		} {
			check(t, file)
		}
	}
}

func TestToleratesOptimizationFlag(t *testing.T) {
	for _, cflags := range []string{
		"",
		"-O",
	} {
		cflags := cflags
		t.Run(cflags, func(t *testing.T) {
			testenv.MustHaveGoBuild(t)
			testenv.MustHaveCGO(t)
			t.Parallel()

			cmd := exec.Command(testenv.GoToolPath(t), "build", path("issue14669.golang"))
			cmd.Env = append(os.Environ(), "CGO_CFLAGS="+cflags)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("%#q: %v:\n%s", cmd, err, out)
			}
		})
	}
}

func TestMallocCrashesOnNil(t *testing.T) {
	testenv.MustHaveCGO(t)
	testenv.MustHaveGoRun(t)
	t.Parallel()

	cmd := exec.Command(testenv.GoToolPath(t), "run", path("malloc.golang"))
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Logf("%#q:\n%s", cmd, out)
		t.Fatalf("succeeded unexpectedly")
	}
}

func TestNotMatchedCFunction(t *testing.T) {
	file := "notmatchedcfunction.golang"
	check(t, file)
}

func TestIncompatibleDeclarations(t *testing.T) {
	testenv.MustHaveCGO(t)
	testenv.MustHaveGoRun(t)
	t.Parallel()
	expect(t, []*regexp.Regexp{
		regexp.MustCompile("inconsistent definitions for C[.]f"),
		regexp.MustCompile("inconsistent definitions for C[.]g"),
	}, "issue67699a.golang", "issue67699b.golang")
}
