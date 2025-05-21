// run

//golang:build !nacl && !js && !wasip1 && gc

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Run the linkx test.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// test(" ") // old deprecated & removed syntax
	test("=") // new syntax
}

func test(sep string) {
	// Successful run
	cmd := exec.Command("golang", "run", "-ldflags=-X main.tbd"+sep+"hello -X main.overwrite"+sep+"trumped -X main.nosuchsymbol"+sep+"neverseen", "linkx.golang")
	var out, errbuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		fmt.Println(errbuf.String())
		fmt.Println(out.String())
		fmt.Println(err)
		os.Exit(1)
	}

	want := "hello\nhello\nhello\ntrumped\ntrumped\ntrumped\n"
	golangt := out.String()
	if golangt != want {
		fmt.Printf("golangt %q want %q\n", golangt, want)
		os.Exit(1)
	}

	// Issue 8810
	cmd = exec.Command("golang", "run", "-ldflags=-X main.tbd", "linkx.golang")
	_, err = cmd.CombinedOutput()
	if err == nil {
		fmt.Println("-X linker flag should not accept keys without values")
		os.Exit(1)
	}

	// Issue 9621
	cmd = exec.Command("golang", "run", "-ldflags=-X main.b=false -X main.x=42", "linkx.golang")
	outx, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Println("-X linker flag should not overwrite non-strings")
		os.Exit(1)
	}
	outstr := string(outx)
	if !strings.Contains(outstr, "main.b") {
		fmt.Printf("-X linker flag did not diagnose overwrite of main.b:\n%s\n", outstr)
		os.Exit(1)
	}
	if !strings.Contains(outstr, "main.x") {
		fmt.Printf("-X linker flag did not diagnose overwrite of main.x:\n%s\n", outstr)
		os.Exit(1)
	}
}
