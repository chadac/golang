// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package scripttest adapts the script engine for use in tests.
package scripttest

import (
	"internal/testenv"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// SetupTestGoRoot sets up a temporary GOROOT for use with script test
// execution. It copies the existing golangroot bin and pkg dirs using
// symlinks (if possible) or raw copying. Return value is the path to
// the newly created testgolangroot dir.
func SetupTestGoRoot(t *testing.T, tmpdir string, golangroot string) string {
	mustMkdir := func(path string) {
		if err := os.MkdirAll(path, 0777); err != nil {
			t.Fatalf("SetupTestGoRoot mkdir %s failed: %v", path, err)
		}
	}

	replicateDir := func(srcdir, dstdir string) {
		files, err := os.ReadDir(srcdir)
		if err != nil {
			t.Fatalf("inspecting %s: %v", srcdir, err)
		}
		for _, file := range files {
			fn := file.Name()
			linkOrCopy(t, filepath.Join(srcdir, fn), filepath.Join(dstdir, fn))
		}
	}

	// Create various dirs in testgolangroot.
	findToolOnce.Do(func() { findToolSub(t) })
	if toolsub == "" {
		t.Fatal("failed to find toolsub")
	}

	tomake := []string{
		"bin",
		"src",
		"pkg",
		filepath.Join("pkg", "include"),
		toolsub,
	}
	made := []string{}
	tgr := filepath.Join(tmpdir, "testgolangroot")
	mustMkdir(tgr)
	for _, targ := range tomake {
		path := filepath.Join(tgr, targ)
		mustMkdir(path)
		made = append(made, path)
	}

	// Replicate selected portions of the content.
	replicateDir(filepath.Join(golangroot, "bin"), made[0])
	replicateDir(filepath.Join(golangroot, "src"), made[1])
	replicateDir(filepath.Join(golangroot, "pkg", "include"), made[3])
	replicateDir(filepath.Join(golangroot, toolsub), made[4])

	return tgr
}

// ReplaceGoToolInTestGoRoot replaces the golang tool binary toolname with
// an alternate executable newtoolpath within a test GOROOT directory
// previously created by SetupTestGoRoot.
func ReplaceGoToolInTestGoRoot(t *testing.T, testgolangroot, toolname, newtoolpath string) {
	findToolOnce.Do(func() { findToolSub(t) })
	if toolsub == "" {
		t.Fatal("failed to find toolsub")
	}

	exename := toolname
	if runtime.GOOS == "windows" {
		exename += ".exe"
	}
	toolpath := filepath.Join(testgolangroot, toolsub, exename)
	if err := os.Remove(toolpath); err != nil {
		t.Fatalf("removing %s: %v", toolpath, err)
	}
	linkOrCopy(t, newtoolpath, toolpath)
}

// toolsub is the tool subdirectory underneath GOROOT.
var toolsub string

// findToolOnce runs findToolSub only once.
var findToolOnce sync.Once

// findToolSub sets toolsub to the value used by the current golang command.
func findToolSub(t *testing.T) {
	golangcmd := testenv.Command(t, testenv.GoToolPath(t), "env", "GOHOSTARCH")
	golangcmd = testenv.CleanCmdEnv(golangcmd)
	golangHostArchBytes, err := golangcmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s failed: %v\n%s", golangcmd, err, golangHostArchBytes)
	}
	golangHostArch := strings.TrimSpace(string(golangHostArchBytes))
	toolsub = filepath.Join("pkg", "tool", runtime.GOOS+"_"+golangHostArch)
}

// linkOrCopy creates a link to src at dst, or if the symlink fails
// (platform doesn't support) then copies src to dst.
func linkOrCopy(t *testing.T, src, dst string) {
	err := os.Symlink(src, dst)
	if err == nil {
		return
	}
	fi, err := os.Stat(src)
	if err != nil {
		t.Fatalf("copying %s to %s: %v", src, dst, err)
	}
	if fi.IsDir() {
		if err := os.CopyFS(dst, os.DirFS(src)); err != nil {
			t.Fatalf("copying %s to %s: %v", src, dst, err)
		}
		return
	}
	srcf, err := os.Open(src)
	if err != nil {
		t.Fatalf("copying %s to %s: %v", src, dst, err)
	}
	defer srcf.Close()
	perm := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	dstf, err := os.OpenFile(dst, perm, 0o777)
	if err != nil {
		t.Fatalf("copying %s to %s: %v", src, dst, err)
	}
	_, err = io.Copy(dstf, srcf)
	if closeErr := dstf.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		t.Fatalf("copying %s to %s: %v", src, dst, err)
	}
}
