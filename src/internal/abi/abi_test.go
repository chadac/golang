// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package abi_test

import (
	"internal/abi"
	"internal/testenv"
	"path/filepath"
	"strings"
	"testing"
)

func TestFuncPC(t *testing.T) {
	// Test that FuncPC* can get correct function PC.
	pcFromAsm := abi.FuncPCTestFnAddr

	// Test FuncPC for locally defined function
	pcFromGo := abi.FuncPCTest()
	if pcFromGo != pcFromAsm {
		t.Errorf("FuncPC returns wrong PC, want %x, golangt %x", pcFromAsm, pcFromGo)
	}

	// Test FuncPC for imported function
	pcFromGo = abi.FuncPCABI0(abi.FuncPCTestFn)
	if pcFromGo != pcFromAsm {
		t.Errorf("FuncPC returns wrong PC, want %x, golangt %x", pcFromAsm, pcFromGo)
	}
}

func TestFuncPCCompileError(t *testing.T) {
	// Test that FuncPC* on a function of a mismatched ABI is rejected.
	testenv.MustHaveGoBuild(t)

	// We want to test internal package, which we cannot normally import.
	// Run the assembler and compiler manually.
	tmpdir := t.TempDir()
	asmSrc := filepath.Join("testdata", "x.s")
	golangSrc := filepath.Join("testdata", "x.golang")
	symabi := filepath.Join(tmpdir, "symabi")
	obj := filepath.Join(tmpdir, "x.o")

	// Write an importcfg file for the dependencies of the package.
	importcfgfile := filepath.Join(tmpdir, "hello.importcfg")
	testenv.WriteImportcfg(t, importcfgfile, nil, "internal/abi")

	// parse assembly code for symabi.
	cmd := testenv.Command(t, testenv.GoToolPath(t), "tool", "asm", "-p=p", "-gensymabis", "-o", symabi, asmSrc)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("golang tool asm -gensymabis failed: %v\n%s", err, out)
	}

	// compile golang code.
	cmd = testenv.Command(t, testenv.GoToolPath(t), "tool", "compile", "-importcfg="+importcfgfile, "-p=p", "-symabis", symabi, "-o", obj, golangSrc)
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("golang tool compile did not fail")
	}

	// Expect errors in line 17, 18, 20, no errors on other lines.
	want := []string{"x.golang:17", "x.golang:18", "x.golang:20"}
	golangt := strings.Split(string(out), "\n")
	if golangt[len(golangt)-1] == "" {
		golangt = golangt[:len(golangt)-1] // remove last empty line
	}
	for i, s := range golangt {
		if !strings.Contains(s, want[i]) {
			t.Errorf("did not error on line %s", want[i])
		}
	}
	if len(golangt) != len(want) {
		t.Errorf("unexpected number of errors, want %d, golangt %d", len(want), len(golangt))
	}
	if t.Failed() {
		t.Logf("output:\n%s", string(out))
	}
}
