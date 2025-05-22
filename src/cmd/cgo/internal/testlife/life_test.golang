// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package life_test

import (
	"bytes"
	"cmd/cgolang/internal/cgolangtest"
	"internal/testenv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	GOPATH, err := os.MkdirTemp("", "cgolanglife")
	if err != nil {
		log.Panic(err)
	}
	defer os.RemoveAll(GOPATH)
	os.Setenv("GOPATH", GOPATH)

	// Copy testdata into GOPATH/src/cgolanglife, along with a golang.mod file
	// declaring the same path.
	modRoot := filepath.Join(GOPATH, "src", "cgolanglife")
	if err := cgolangtest.OverlayDir(modRoot, "testdata"); err != nil {
		log.Panic(err)
	}
	if err := os.Chdir(modRoot); err != nil {
		log.Panic(err)
	}
	os.Setenv("PWD", modRoot)
	if err := os.WriteFile("golang.mod", []byte("module cgolanglife\n"), 0666); err != nil {
		log.Panic(err)
	}

	return m.Run()
}

// TestTestRun runs a test case for cgolang //export.
func TestTestRun(t *testing.T) {
	testenv.MustHaveGoRun(t)
	testenv.MustHaveCGO(t)

	cmd := exec.Command(testenv.GoToolPath(t), "run", "main.golang")
	golangt, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v: %s\n%s", cmd, err, golangt)
	}
	want, err := os.ReadFile("main.out")
	if err != nil {
		t.Fatal("reading golanglden output:", err)
	}
	if !bytes.Equal(golangt, want) {
		t.Errorf("'%v' output does not match expected in main.out. Instead saw:\n%s", cmd, golangt)
	}
}
