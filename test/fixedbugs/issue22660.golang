// run

//golang:build !js && !wasip1 && gc

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	f, err := ioutil.TempFile("", "issue22660.golang")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	// path must appear in error messages even if we strip them with -trimpath
	path := filepath.Join("users", "xxx", "golang")
	var src bytes.Buffer
	fmt.Fprintf(&src, "//line %s:1\n", filepath.Join(path, "foo.golang"))

	if err := ioutil.WriteFile(f.Name(), src.Bytes(), 0660); err != nil {
		log.Fatal(err)
	}

	out, err := exec.Command("golang", "tool", "compile", "-p=p", fmt.Sprintf("-trimpath=%s", path), f.Name()).CombinedOutput()
	if err == nil {
		log.Fatalf("expected compiling %s to fail", f.Name())
	}

	if !strings.HasPrefix(string(out), path) {
		log.Fatalf("expected full path (%s) in error message, golangt:\n%s", path, out)
	}
}
