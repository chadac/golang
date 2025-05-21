// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package stdio_test

import (
	"bytes"
	"cmd/cgolang/internal/cgolangtest"
	"internal/testenv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	GOPATH, err := os.MkdirTemp("", "cgolangstdio")
	if err != nil {
		log.Panic(err)
	}
	defer os.RemoveAll(GOPATH)
	os.Setenv("GOPATH", GOPATH)

	// Copy testdata into GOPATH/src/cgolangstdio, along with a golang.mod file
	// declaring the same path.
	modRoot := filepath.Join(GOPATH, "src", "cgolangstdio")
	if err := cgolangtest.OverlayDir(modRoot, "testdata"); err != nil {
		log.Panic(err)
	}
	if err := os.Chdir(modRoot); err != nil {
		log.Panic(err)
	}
	os.Setenv("PWD", modRoot)
	if err := os.WriteFile("golang.mod", []byte("module cgolangstdio\n"), 0666); err != nil {
		log.Panic(err)
	}

	return m.Run()
}

// TestTestRun runs a cgolang test that doesn't depend on non-standard libraries.
func TestTestRun(t *testing.T) {
	testenv.MustHaveGoRun(t)
	testenv.MustHaveCGO(t)

	for _, file := range [...]string{
		"chain.golang",
		"fib.golang",
		"hello.golang",
	} {
		file := file
		wantFile := strings.Replace(file, ".golang", ".out", 1)
		t.Run(file, func(t *testing.T) {
			cmd := exec.Command(testenv.GoToolPath(t), "run", file)
			golangt, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%v: %s\n%s", cmd, err, golangt)
			}
			golangt = bytes.ReplaceAll(golangt, []byte("\r\n"), []byte("\n"))
			want, err := os.ReadFile(wantFile)
			if err != nil {
				t.Fatal("reading golanglden output:", err)
			}
			if !bytes.Equal(golangt, want) {
				t.Errorf("'%v' output does not match expected in %s. Instead saw:\n%s", cmd, wantFile, golangt)
			}
		})
	}
}
