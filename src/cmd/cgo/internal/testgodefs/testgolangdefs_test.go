// Copyright 2019 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package testgolangdefs

import (
	"bytes"
	"internal/testenv"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// We are testing cgolang -golangdefs, which translates Golang files that use
// import "C" into Golang files with Golang definitions of types defined in the
// import "C" block.  Add more tests here.
var filePrefixes = []string{
	"anonunion",
	"bitfields",
	"issue8478",
	"fieldtypedef",
	"issue37479",
	"issue37621",
	"issue38649",
	"issue39534",
	"issue48396",
}

func TestGolangDefs(t *testing.T) {
	testenv.MustHaveGolangRun(t)
	testenv.MustHaveCGO(t)

	testdata, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}

	golangpath, err := os.MkdirTemp("", "testgolangdefs-golangpath")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(golangpath)

	dir := filepath.Join(golangpath, "src", "testgolangdefs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	for _, fp := range filePrefixes {
		cmd := exec.Command(testenv.GolangToolPath(t), "tool", "cgolang",
			"-golangdefs",
			"-srcdir", testdata,
			"-objdir", dir,
			fp+".golang")
		cmd.Stderr = new(bytes.Buffer)

		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("%#q: %v\n%s", cmd, err, cmd.Stderr)
		}

		fn := fp + "_defs.golang"
		if err := os.WriteFile(filepath.Join(dir, fn), out, 0644); err != nil {
			t.Fatal(err)
		}

		// Verify that command line arguments are not rewritten in the generated comment,
		// see golang.dev/issue/52063
		hasGeneratedByComment := false
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			cgolangExe := "cgolang"
			if runtime.GOOS == "windows" {
				cgolangExe = "cgolang.exe"
			}
			if !strings.HasPrefix(line, "// "+cgolangExe+" -golangdefs") {
				continue
			}
			if want := "// " + cgolangExe + " " + strings.Join(cmd.Args[3:], " "); line != want {
				t.Errorf("%s: golangt generated comment %q, want %q", fn, line, want)
			}
			hasGeneratedByComment = true
			break
		}

		if !hasGeneratedByComment {
			t.Errorf("%s: comment with generating cgolang -golangdefs command not found", fn)
		}
	}

	main, err := os.ReadFile(filepath.Join("testdata", "main.golang"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.golang"), main, 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "golang.mod"), []byte("module testgolangdefs\ngolang 1.14\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Use 'golang run' to build and run the resulting binary in a single step,
	// instead of invoking 'golang build' and the resulting binary separately, so that
	// this test can pass on mobile builders, which do not copy artifacts back
	// from remote invocations.
	cmd := exec.Command(testenv.GolangToolPath(t), "run", ".")
	cmd.Env = append(os.Environ(), "GOPATH="+golangpath)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%#q [%s]: %v\n%s", cmd, dir, err, out)
	}
}
