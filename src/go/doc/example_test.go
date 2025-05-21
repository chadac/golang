// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package doc_test

import (
	"bytes"
	"fmt"
	"golang/ast"
	"golang/doc"
	"golang/format"
	"golang/parser"
	"golang/token"
	"internal/diff"
	"internal/txtar"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestExamples(t *testing.T) {
	dir := filepath.Join("testdata", "examples")
	filenames, err := filepath.Glob(filepath.Join(dir, "*.golang"))
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range filenames {
		t.Run(strings.TrimSuffix(filepath.Base(filename), ".golang"), func(t *testing.T) {
			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			golangldenFilename := strings.TrimSuffix(filename, ".golang") + ".golanglden"
			archive, err := txtar.ParseFile(golangldenFilename)
			if err != nil {
				t.Fatal(err)
			}
			golanglden := map[string]string{}
			for _, f := range archive.Files {
				golanglden[f.Name] = strings.TrimSpace(string(f.Data))
			}

			// Collect the results of doc.Examples in a map keyed by example name.
			examples := map[string]*doc.Example{}
			for _, e := range doc.Examples(astFile) {
				examples[e.Name] = e
				// Treat missing sections in the golanglden as empty.
				for _, kind := range []string{"Play", "Output"} {
					key := e.Name + "." + kind
					if _, ok := golanglden[key]; !ok {
						golanglden[key] = ""
					}
				}
			}

			// Each section in the golanglden file corresponds to an example we expect
			// to see.
			for sectionName, want := range golanglden {
				name, kind, found := strings.Cut(sectionName, ".")
				if !found {
					t.Fatalf("bad section name %q, want EXAMPLE_NAME.KIND", sectionName)
				}
				ex := examples[name]
				if ex == nil {
					t.Fatalf("no example named %q", name)
				}

				var golangt string
				switch kind {
				case "Play":
					golangt = strings.TrimSpace(formatFile(t, fset, ex.Play))

				case "Output":
					golangt = strings.TrimSpace(ex.Output)
				default:
					t.Fatalf("bad section kind %q", kind)
				}

				if golangt != want {
					t.Errorf("%s mismatch:\n%s", sectionName,
						diff.Diff("want", []byte(want), "golangt", []byte(golangt)))
				}
			}
		})
	}
}

func formatFile(t *testing.T, fset *token.FileSet, n *ast.File) string {
	t.Helper()
	if n == nil {
		return "<nil>"
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, n); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// This example illustrates how to use NewFromFiles
// to compute package documentation with examples.
func ExampleNewFromFiles() {
	// src and test are two source files that make up
	// a package whose documentation will be computed.
	const src = `
// This is the package comment.
package p

import "fmt"

// This comment is associated with the Greet function.
func Greet(who string) {
	fmt.Printf("Hello, %s!\n", who)
}
`
	const test = `
package p_test

// This comment is associated with the ExampleGreet_world example.
func ExampleGreet_world() {
	Greet("world")
}
`

	// Create the AST by parsing src and test.
	fset := token.NewFileSet()
	files := []*ast.File{
		mustParse(fset, "src.golang", src),
		mustParse(fset, "src_test.golang", test),
	}

	// Compute package documentation with examples.
	p, err := doc.NewFromFiles(fset, files, "example.com/p")
	if err != nil {
		panic(err)
	}

	fmt.Printf("package %s - %s", p.Name, p.Doc)
	fmt.Printf("func %s - %s", p.Funcs[0].Name, p.Funcs[0].Doc)
	fmt.Printf(" ⤷ example with suffix %q - %s", p.Funcs[0].Examples[0].Suffix, p.Funcs[0].Examples[0].Doc)

	// Output:
	// package p - This is the package comment.
	// func Greet - This comment is associated with the Greet function.
	//  ⤷ example with suffix "world" - This comment is associated with the ExampleGreet_world example.
}

func TestClassifyExamples(t *testing.T) {
	const src = `
package p

const Const1 = 0
var   Var1   = 0

type (
	Type1     int
	Type1_Foo int
	Type1_foo int
	type2     int

	Embed struct { Type1 }
	Uembed struct { type2 }
)

func Func1()     {}
func Func1_Foo() {}
func Func1_foo() {}
func func2()     {}

func (Type1) Func1() {}
func (Type1) Func1_Foo() {}
func (Type1) Func1_foo() {}
func (Type1) func2() {}

func (type2) Func1() {}

type (
	Conflict          int
	Conflict_Conflict int
	Conflict_conflict int
)

func (Conflict) Conflict() {}

func GFunc[T any]() {}

type GType[T any] int

func (GType[T]) M() {}
`
	const test = `
package p_test

func ExampleConst1() {} // invalid - no support for consts and vars
func ExampleVar1()   {} // invalid - no support for consts and vars

func Example()               {}
func Example_()              {} // invalid - suffix must start with a lower-case letter
func Example_suffix()        {}
func Example_suffix_xX_X_x() {}
func Example_世界()           {} // invalid - suffix must start with a lower-case letter
func Example_123()           {} // invalid - suffix must start with a lower-case letter
func Example_BadSuffix()     {} // invalid - suffix must start with a lower-case letter

func ExampleType1()               {}
func ExampleType1_()              {} // invalid - suffix must start with a lower-case letter
func ExampleType1_suffix()        {}
func ExampleType1_BadSuffix()     {} // invalid - suffix must start with a lower-case letter
func ExampleType1_Foo()           {}
func ExampleType1_Foo_suffix()    {}
func ExampleType1_Foo_BadSuffix() {} // invalid - suffix must start with a lower-case letter
func ExampleType1_foo()           {}
func ExampleType1_foo_suffix()    {}
func ExampleType1_foo_Suffix()    {} // matches Type1, instead of Type1_foo
func Exampletype2()               {} // invalid - cannot match unexported

func ExampleFunc1()               {}
func ExampleFunc1_()              {} // invalid - suffix must start with a lower-case letter
func ExampleFunc1_suffix()        {}
func ExampleFunc1_BadSuffix()     {} // invalid - suffix must start with a lower-case letter
func ExampleFunc1_Foo()           {}
func ExampleFunc1_Foo_suffix()    {}
func ExampleFunc1_Foo_BadSuffix() {} // invalid - suffix must start with a lower-case letter
func ExampleFunc1_foo()           {}
func ExampleFunc1_foo_suffix()    {}
func ExampleFunc1_foo_Suffix()    {} // matches Func1, instead of Func1_foo
func Examplefunc1()               {} // invalid - cannot match unexported

func ExampleType1_Func1()               {}
func ExampleType1_Func1_()              {} // invalid - suffix must start with a lower-case letter
func ExampleType1_Func1_suffix()        {}
func ExampleType1_Func1_BadSuffix()     {} // invalid - suffix must start with a lower-case letter
func ExampleType1_Func1_Foo()           {}
func ExampleType1_Func1_Foo_suffix()    {}
func ExampleType1_Func1_Foo_BadSuffix() {} // invalid - suffix must start with a lower-case letter
func ExampleType1_Func1_foo()           {}
func ExampleType1_Func1_foo_suffix()    {}
func ExampleType1_Func1_foo_Suffix()    {} // matches Type1.Func1, instead of Type1.Func1_foo
func ExampleType1_func2()               {} // matches Type1, instead of Type1.func2

func ExampleEmbed_Func1()         {} // invalid - no support for forwarded methods from embedding exported type
func ExampleUembed_Func1()        {} // methods from embedding unexported types are OK
func ExampleUembed_Func1_suffix() {}

func ExampleConflict_Conflict()        {} // ambiguous with either Conflict or Conflict_Conflict type
func ExampleConflict_conflict()        {} // ambiguous with either Conflict or Conflict_conflict type
func ExampleConflict_Conflict_suffix() {} // ambiguous with either Conflict or Conflict_Conflict type
func ExampleConflict_conflict_suffix() {} // ambiguous with either Conflict or Conflict_conflict type

func ExampleGFunc() {}
func ExampleGFunc_suffix() {}

func ExampleGType_M() {}
func ExampleGType_M_suffix() {}
`

	// Parse literal source code as a *doc.Package.
	fset := token.NewFileSet()
	files := []*ast.File{
		mustParse(fset, "src.golang", src),
		mustParse(fset, "src_test.golang", test),
	}
	p, err := doc.NewFromFiles(fset, files, "example.com/p")
	if err != nil {
		t.Fatalf("doc.NewFromFiles: %v", err)
	}

	// Collect the association of examples to top-level identifiers.
	golangt := map[string][]string{}
	golangt[""] = exampleNames(p.Examples)
	for _, f := range p.Funcs {
		golangt[f.Name] = exampleNames(f.Examples)
	}
	for _, t := range p.Types {
		golangt[t.Name] = exampleNames(t.Examples)
		for _, f := range t.Funcs {
			golangt[f.Name] = exampleNames(f.Examples)
		}
		for _, m := range t.Methods {
			golangt[t.Name+"."+m.Name] = exampleNames(m.Examples)
		}
	}

	want := map[string][]string{
		"": {"", "suffix", "suffix_xX_X_x"}, // Package-level examples.

		"Type1":     {"", "foo_Suffix", "func2", "suffix"},
		"Type1_Foo": {"", "suffix"},
		"Type1_foo": {"", "suffix"},

		"Func1":     {"", "foo_Suffix", "suffix"},
		"Func1_Foo": {"", "suffix"},
		"Func1_foo": {"", "suffix"},

		"Type1.Func1":     {"", "foo_Suffix", "suffix"},
		"Type1.Func1_Foo": {"", "suffix"},
		"Type1.Func1_foo": {"", "suffix"},

		"Uembed.Func1": {"", "suffix"},

		// These are implementation dependent due to the ambiguous parsing.
		"Conflict_Conflict": {"", "suffix"},
		"Conflict_conflict": {"", "suffix"},

		"GFunc":   {"", "suffix"},
		"GType.M": {"", "suffix"},
	}

	for id := range golangt {
		if !reflect.DeepEqual(golangt[id], want[id]) {
			t.Errorf("classification mismatch for %q:\ngolangt  %q\nwant %q", id, golangt[id], want[id])
		}
		delete(want, id)
	}
	if len(want) > 0 {
		t.Errorf("did not find:\n%q", want)
	}
}

func exampleNames(exs []*doc.Example) (out []string) {
	for _, ex := range exs {
		out = append(out, ex.Suffix)
	}
	return out
}

func mustParse(fset *token.FileSet, filename, src string) *ast.File {
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}
