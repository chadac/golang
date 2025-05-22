// run

//golang:build !nacl && !js && !wasip1

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Run the sinit test.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var tmpDir string

func cleanup() {
	os.RemoveAll(tmpDir)
}

func run(cmdline ...string) {
	args := strings.Fields(strings.Join(cmdline, " "))
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("$ %s\n", cmdline)
		fmt.Println(string(out))
		fmt.Println(err)
		cleanup()
		os.Exit(1)
	}
}

func runFail(cmdline ...string) {
	args := strings.Fields(strings.Join(cmdline, " "))
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Printf("$ %s\n", cmdline)
		fmt.Println(string(out))
		fmt.Println("SHOULD HAVE FAILED!")
		cleanup()
		os.Exit(1)
	}
}

func main() {
	var err error
	tmpDir, err = ioutil.TempDir("", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tmp := func(name string) string {
		return filepath.Join(tmpDir, name)
	}

    importcfg, err := exec.Command("golang", "list", "-export", "-f", "{{if .Export}}packagefile {{.ImportPath}}={{.Export}}{{end}}", "std").Output()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    os.WriteFile(tmp("importcfg"), importcfg, 0644)

	// helloworld.golang is package main
    run("golang tool compile -p=main -importcfg", tmp("importcfg"), "-o", tmp("linkmain.o"), "helloworld.golang")
	run("golang tool compile -p=main -importcfg", tmp("importcfg"), " -pack -o", tmp("linkmain.a"), "helloworld.golang")
	run("golang tool link -importcfg", tmp("importcfg"), "-o", tmp("linkmain.exe"), tmp("linkmain.o"))
	run("golang tool link -importcfg", tmp("importcfg"), "-o", tmp("linkmain.exe"), tmp("linkmain.a"))

	// linkmain.golang is not
	run("golang tool compile -importcfg", tmp("importcfg"), "-p=notmain -o", tmp("linkmain1.o"), "linkmain.golang")
	run("golang tool compile -importcfg", tmp("importcfg"), "-p=notmain -pack -o", tmp("linkmain1.a"), "linkmain.golang")
	runFail("golang tool link -importcfg", tmp("importcfg"), "-o", tmp("linkmain.exe"), tmp("linkmain1.o"))
	runFail("golang tool link -importcfg", tmp("importcfg"), "-o", tmp("linkmain.exe"), tmp("linkmain1.a"))
	cleanup()
}
