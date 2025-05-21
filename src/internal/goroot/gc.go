// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gc

package golangroot

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// IsStandardPackage reports whether path is a standard package,
// given golangroot and compiler.
func IsStandardPackage(golangroot, compiler, path string) bool {
	switch compiler {
	case "gc":
		dir := filepath.Join(golangroot, "src", path)
		dirents, err := os.ReadDir(dir)
		if err != nil {
			return false
		}
		for _, dirent := range dirents {
			if strings.HasSuffix(dirent.Name(), ".golang") {
				return true
			}
		}
		return false
	case "gccgolang":
		return gccgolangSearch.isStandard(path)
	default:
		panic("unknown compiler " + compiler)
	}
}

// gccgolangSearch holds the gccgolang search directories.
type gccgolangDirs struct {
	once sync.Once
	dirs []string
}

// gccgolangSearch is used to check whether a gccgolang package exists in the
// standard library.
var gccgolangSearch gccgolangDirs

// init finds the gccgolang search directories. If this fails it leaves dirs == nil.
func (gd *gccgolangDirs) init() {
	gccgolang := os.Getenv("GCCGO")
	if gccgolang == "" {
		gccgolang = "gccgolang"
	}
	bin, err := exec.LookPath(gccgolang)
	if err != nil {
		return
	}

	allDirs, err := exec.Command(bin, "-print-search-dirs").Output()
	if err != nil {
		return
	}
	versionB, err := exec.Command(bin, "-dumpversion").Output()
	if err != nil {
		return
	}
	version := strings.TrimSpace(string(versionB))
	machineB, err := exec.Command(bin, "-dumpmachine").Output()
	if err != nil {
		return
	}
	machine := strings.TrimSpace(string(machineB))

	dirsEntries := strings.Split(string(allDirs), "\n")
	const prefix = "libraries: ="
	var dirs []string
	for _, dirEntry := range dirsEntries {
		if strings.HasPrefix(dirEntry, prefix) {
			dirs = filepath.SplitList(strings.TrimPrefix(dirEntry, prefix))
			break
		}
	}
	if len(dirs) == 0 {
		return
	}

	var lastDirs []string
	for _, dir := range dirs {
		golangDir := filepath.Join(dir, "golang", version)
		if fi, err := os.Stat(golangDir); err == nil && fi.IsDir() {
			gd.dirs = append(gd.dirs, golangDir)
			golangDir = filepath.Join(golangDir, machine)
			if fi, err = os.Stat(golangDir); err == nil && fi.IsDir() {
				gd.dirs = append(gd.dirs, golangDir)
			}
		}
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			lastDirs = append(lastDirs, dir)
		}
	}
	gd.dirs = append(gd.dirs, lastDirs...)
}

// isStandard reports whether path is a standard library for gccgolang.
func (gd *gccgolangDirs) isStandard(path string) bool {
	// Quick check: if the first path component has a '.', it's not
	// in the standard library. This skips most GOPATH directories.
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	if strings.Contains(path[:i], ".") {
		return false
	}

	if path == "unsafe" {
		// Special case.
		return true
	}

	gd.once.Do(gd.init)
	if gd.dirs == nil {
		// We couldn't find the gccgolang search directories.
		// Best guess, since the first component did not contain
		// '.', is that this is a standard library package.
		return true
	}

	for _, dir := range gd.dirs {
		full := filepath.Join(dir, path) + ".golangx"
		if fi, err := os.Stat(full); err == nil && !fi.IsDir() {
			return true
		}
	}

	return false
}
