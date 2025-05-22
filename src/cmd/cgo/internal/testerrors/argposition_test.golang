// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 42580: cmd/cgolang: shifting identifier position in ast

package errorstest

import (
	"bytes"
	"fmt"
	"golang/ast"
	"golang/parser"
	"golang/token"
	"internal/testenv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type ShortPosition struct {
	Line    int
	Column  int
	Visited bool
}

type IdentPositionInfo map[string][]ShortPosition

type Visitor struct {
	identPosInfo IdentPositionInfo
	fset         *token.FileSet
	t            *testing.T
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if ident, ok := node.(*ast.Ident); ok {
		if expectedPositions, ok := v.identPosInfo[ident.Name]; ok {
			golangtMatch := false
			var errorMessage strings.Builder
			for caseIndex, expectedPos := range expectedPositions {
				actualPosition := v.fset.PositionFor(ident.Pos(), true)
				errorOccurred := false
				if expectedPos.Line != actualPosition.Line {
					fmt.Fprintf(&errorMessage, "wrong line number for ident %s: expected: %d golangt: %d\n", ident.Name, expectedPos.Line, actualPosition.Line)
					errorOccurred = true
				}
				if expectedPos.Column != actualPosition.Column {
					fmt.Fprintf(&errorMessage, "wrong column number for ident %s: expected: %d golangt: %d\n", ident.Name, expectedPos.Column, actualPosition.Column)
					errorOccurred = true
				}
				if errorOccurred {
					continue
				}
				golangtMatch = true
				expectedPositions[caseIndex].Visited = true
			}

			if !golangtMatch {
				v.t.Error(errorMessage.String())
			}
		}
	}
	return v
}

func TestArgumentsPositions(t *testing.T) {
	testenv.MustHaveCGO(t)
	testenv.MustHaveExec(t)

	testdata, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}

	tmpPath := t.TempDir()

	dir := filepath.Join(tmpPath, "src", "testpositions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(testenv.GoToolPath(t), "tool", "cgolang",
		"-srcdir", testdata,
		"-objdir", dir,
		"issue42580.golang")
	cmd.Stderr = new(bytes.Buffer)

	err = cmd.Run()
	if err != nil {
		t.Fatalf("%s: %v\n%s", cmd, err, cmd.Stderr)
	}
	mainProcessed, err := os.ReadFile(filepath.Join(dir, "issue42580.cgolang1.golang"))
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", mainProcessed, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	expectation := IdentPositionInfo{
		"checkedPointer": []ShortPosition{
			ShortPosition{
				Line:   32,
				Column: 56,
			},
		},
		"singleInnerPointerChecked": []ShortPosition{
			ShortPosition{
				Line:   37,
				Column: 91,
			},
		},
		"doublePointerChecked": []ShortPosition{
			ShortPosition{
				Line:   42,
				Column: 91,
			},
		},
	}
	for _, decl := range f.Decls {
		if fdecl, ok := decl.(*ast.FuncDecl); ok {
			ast.Walk(&Visitor{expectation, fset, t}, fdecl.Body)
		}
	}
	for ident, positions := range expectation {
		for _, position := range positions {
			if !position.Visited {
				t.Errorf("Position %d:%d missed for %s ident", position.Line, position.Column, ident)
			}
		}
	}
}
