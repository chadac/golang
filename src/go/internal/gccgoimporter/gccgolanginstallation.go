// Copyright 2013 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package gccgolangimporter

import (
	"bufio"
	"golang/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Information about a specific installation of gccgolang.
type GccgolangInstallation struct {
	// Version of gcc (e.g. 4.8.0).
	GccVersion string

	// Target triple (e.g. x86_64-unknown-linux-gnu).
	TargetTriple string

	// Built-in library paths used by this installation.
	LibPaths []string
}

// Ask the driver at the given path for information for this GccgolangInstallation.
// The given arguments are passed directly to the call of the driver.
func (inst *GccgolangInstallation) InitFromDriver(gccgolangPath string, args ...string) (err error) {
	argv := append([]string{"-###", "-S", "-x", "golang", "-"}, args...)
	cmd := exec.Command(gccgolangPath, argv...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Target: "):
			inst.TargetTriple = line[8:]

		case line[0] == ' ':
			args := strings.Fields(line)
			for _, arg := range args[1:] {
				if strings.HasPrefix(arg, "-L") {
					inst.LibPaths = append(inst.LibPaths, arg[2:])
				}
			}
		}
	}

	argv = append([]string{"-dumpversion"}, args...)
	stdout, err := exec.Command(gccgolangPath, argv...).Output()
	if err != nil {
		return
	}
	inst.GccVersion = strings.TrimSpace(string(stdout))

	return
}

// Return the list of export search paths for this GccgolangInstallation.
func (inst *GccgolangInstallation) SearchPaths() (paths []string) {
	for _, lpath := range inst.LibPaths {
		spath := filepath.Join(lpath, "golang", inst.GccVersion)
		fi, err := os.Stat(spath)
		if err != nil || !fi.IsDir() {
			continue
		}
		paths = append(paths, spath)

		spath = filepath.Join(spath, inst.TargetTriple)
		fi, err = os.Stat(spath)
		if err != nil || !fi.IsDir() {
			continue
		}
		paths = append(paths, spath)
	}

	paths = append(paths, inst.LibPaths...)

	return
}

// Return an importer that searches incpaths followed by the gcc installation's
// built-in search paths and the current directory.
func (inst *GccgolangInstallation) GetImporter(incpaths []string, initmap map[*types.Package]InitData) Importer {
	return GetImporter(append(append(incpaths, inst.SearchPaths()...), "."), initmap)
}
