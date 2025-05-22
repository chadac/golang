// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"bufio"
	"bytes"
	cmdcover "cmd/cover"
	"flag"
	"fmt"
	"golang/ast"
	"golang/parser"
	"golang/token"
	"internal/testenv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
)

const (
	// Data directory, also the package directory for the test.
	testdata = "testdata"
)

// testcover returns the path to the cmd/cover binary that we are golanging to
// test. At one point this was created via "golang build"; we now reuse the unit
// test executable itself.
func testcover(t testing.TB) string {
	return testenv.Executable(t)
}

// testTempDir is a temporary directory created in TestMain.
var testTempDir string

// If set, this will preserve all the tmpdir files from the test run.
var debug = flag.Bool("debug", false, "keep tmpdir files for debugging")

// TestMain used here so that we can leverage the test executable
// itself as a cmd/cover executable; compare to similar usage in
// the cmd/golang tests.
func TestMain(m *testing.M) {
	if os.Getenv("CMDCOVER_TOOLEXEC") != "" {
		// When CMDCOVER_TOOLEXEC is set, the test binary is also
		// running as a -toolexec wrapper.
		tool := strings.TrimSuffix(filepath.Base(os.Args[1]), ".exe")
		if tool == "cover" {
			// Inject this test binary as cmd/cover in place of the
			// installed tool, so that the golang command's invocations of
			// cover produce coverage for the configuration in which
			// the test was built.
			os.Args = os.Args[1:]
			cmdcover.Main()
		} else {
			cmd := exec.Command(os.Args[1], os.Args[2:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				os.Exit(1)
			}
		}
		os.Exit(0)
	}
	if os.Getenv("CMDCOVER_TEST_RUN_MAIN") != "" {
		// When CMDCOVER_TEST_RUN_MAIN is set, we're reusing the test
		// binary as cmd/cover. In this case we run the main func exported
		// via export_test.golang, and exit; CMDCOVER_TEST_RUN_MAIN is set below
		// for actual test invocations.
		cmdcover.Main()
		os.Exit(0)
	}
	flag.Parse()
	topTmpdir, err := os.MkdirTemp("", "cmd-cover-test-")
	if err != nil {
		log.Fatal(err)
	}
	testTempDir = topTmpdir
	if !*debug {
		defer os.RemoveAll(topTmpdir)
	} else {
		fmt.Fprintf(os.Stderr, "debug: preserving tmpdir %s\n", topTmpdir)
	}
	os.Setenv("CMDCOVER_TEST_RUN_MAIN", "normal")
	os.Exit(m.Run())
}

var tdmu sync.Mutex
var tdcount int

func tempDir(t *testing.T) string {
	tdmu.Lock()
	dir := filepath.Join(testTempDir, fmt.Sprintf("%03d", tdcount))
	tdcount++
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatal(err)
	}
	defer tdmu.Unlock()
	return dir
}

// TestCoverWithToolExec runs a set of subtests that all make use of a
// "-toolexec" wrapper program to invoke the cover test executable
// itself via "golang test -cover".
func TestCoverWithToolExec(t *testing.T) {
	toolexecArg := "-toolexec=" + testcover(t)

	t.Run("CoverHTML", func(t *testing.T) {
		testCoverHTML(t, toolexecArg)
	})
	t.Run("HtmlUnformatted", func(t *testing.T) {
		testHtmlUnformatted(t, toolexecArg)
	})
	t.Run("FuncWithDuplicateLines", func(t *testing.T) {
		testFuncWithDuplicateLines(t, toolexecArg)
	})
	t.Run("MissingTrailingNewlineIssue58370", func(t *testing.T) {
		testMissingTrailingNewlineIssue58370(t, toolexecArg)
	})
}

// Execute this command sequence:
//
//	replace the word LINE with the line number < testdata/test.golang > testdata/test_line.golang
//	testcover -mode=count -var=CoverTest -o ./testdata/test_cover.golang testdata/test_line.golang
//	golang run ./testdata/main.golang ./testdata/test.golang
func TestCover(t *testing.T) {
	testenv.MustHaveGoRun(t)
	t.Parallel()
	dir := tempDir(t)

	// Read in the test file (testTest) and write it, with LINEs specified, to coverInput.
	testTest := filepath.Join(testdata, "test.golang")
	file, err := os.ReadFile(testTest)
	if err != nil {
		t.Fatal(err)
	}
	lines := bytes.Split(file, []byte("\n"))
	for i, line := range lines {
		lines[i] = bytes.ReplaceAll(line, []byte("LINE"), []byte(fmt.Sprint(i+1)))
	}

	// Add a function that is not golangfmt'ed. This used to cause a crash.
	// We don't put it in test.golang because then we would have to golangfmt it.
	// Issue 23927.
	lines = append(lines, []byte("func unFormatted() {"),
		[]byte("\tif true {"),
		[]byte("\t}else{"),
		[]byte("\t}"),
		[]byte("}"))
	lines = append(lines, []byte("func unFormatted2(b bool) {if b{}else{}}"))

	coverInput := filepath.Join(dir, "test_line.golang")
	if err := os.WriteFile(coverInput, bytes.Join(lines, []byte("\n")), 0666); err != nil {
		t.Fatal(err)
	}

	// testcover -mode=count -var=thisNameMustBeVeryLongToCauseOverflowOfCounterIncrementStatementOntoNextLineForTest -o ./testdata/test_cover.golang testdata/test_line.golang
	coverOutput := filepath.Join(dir, "test_cover.golang")
	cmd := testenv.Command(t, testcover(t), "-mode=count", "-var=thisNameMustBeVeryLongToCauseOverflowOfCounterIncrementStatementOntoNextLineForTest", "-o", coverOutput, coverInput)
	run(cmd, t)

	cmd = testenv.Command(t, testcover(t), "-mode=set", "-var=Not_an-identifier", "-o", coverOutput, coverInput)
	err = cmd.Run()
	if err == nil {
		t.Error("Expected cover to fail with an error")
	}

	// Copy testmain to tmpdir, so that it is in the same directory
	// as coverOutput.
	testMain := filepath.Join(testdata, "main.golang")
	b, err := os.ReadFile(testMain)
	if err != nil {
		t.Fatal(err)
	}
	tmpTestMain := filepath.Join(dir, "main.golang")
	if err := os.WriteFile(tmpTestMain, b, 0444); err != nil {
		t.Fatal(err)
	}

	// golang run ./testdata/main.golang ./testdata/test.golang
	cmd = testenv.Command(t, testenv.GoToolPath(t), "run", tmpTestMain, coverOutput)
	run(cmd, t)

	file, err = os.ReadFile(coverOutput)
	if err != nil {
		t.Fatal(err)
	}
	// compiler directive must appear right next to function declaration.
	if golangt, err := regexp.MatchString(".*\n//golang:nosplit\nfunc someFunction().*", string(file)); err != nil || !golangt {
		t.Error("misplaced compiler directive")
	}
	// "golang:linkname" compiler directive should be present.
	if golangt, err := regexp.MatchString(`.*golang\:linkname some\_name some\_name.*`, string(file)); err != nil || !golangt {
		t.Error("'golang:linkname' compiler directive not found")
	}

	// Other comments should be preserved too.
	c := ".*// This comment didn't appear in generated golang code.*"
	if golangt, err := regexp.MatchString(c, string(file)); err != nil || !golangt {
		t.Errorf("non compiler directive comment %q not found", c)
	}
}

// TestDirectives checks that compiler directives are preserved and positioned
// correctly. Directives that occur before top-level declarations should remain
// above those declarations, even if they are not part of the block of
// documentation comments.
func TestDirectives(t *testing.T) {
	testenv.MustHaveExec(t)
	t.Parallel()

	// Read the source file and find all the directives. We'll keep
	// track of whether each one has been seen in the output.
	testDirectives := filepath.Join(testdata, "directives.golang")
	source, err := os.ReadFile(testDirectives)
	if err != nil {
		t.Fatal(err)
	}
	sourceDirectives := findDirectives(source)

	// testcover -mode=atomic ./testdata/directives.golang
	cmd := testenv.Command(t, testcover(t), "-mode=atomic", testDirectives)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	// Check that all directives are present in the output.
	outputDirectives := findDirectives(output)
	foundDirective := make(map[string]bool)
	for _, p := range sourceDirectives {
		foundDirective[p.name] = false
	}
	for _, p := range outputDirectives {
		if found, ok := foundDirective[p.name]; !ok {
			t.Errorf("unexpected directive in output: %s", p.text)
		} else if found {
			t.Errorf("directive found multiple times in output: %s", p.text)
		}
		foundDirective[p.name] = true
	}
	for name, found := range foundDirective {
		if !found {
			t.Errorf("missing directive: %s", name)
		}
	}

	// Check that directives that start with the name of top-level declarations
	// come before the beginning of the named declaration and after the end
	// of the previous declaration.
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, testDirectives, output, 0)
	if err != nil {
		t.Fatal(err)
	}

	prevEnd := 0
	for _, decl := range astFile.Decls {
		var name string
		switch d := decl.(type) {
		case *ast.FuncDecl:
			name = d.Name.Name
		case *ast.GenDecl:
			if len(d.Specs) == 0 {
				// An empty group declaration. We still want to check that
				// directives can be associated with it, so we make up a name
				// to match directives in the test data.
				name = "_empty"
			} else if spec, ok := d.Specs[0].(*ast.TypeSpec); ok {
				name = spec.Name.Name
			}
		}
		pos := fset.Position(decl.Pos()).Offset
		end := fset.Position(decl.End()).Offset
		if name == "" {
			prevEnd = end
			continue
		}
		for _, p := range outputDirectives {
			if !strings.HasPrefix(p.name, name) {
				continue
			}
			if p.offset < prevEnd || pos < p.offset {
				t.Errorf("directive %s does not appear before definition %s", p.text, name)
			}
		}
		prevEnd = end
	}
}

type directiveInfo struct {
	text   string // full text of the comment, not including newline
	name   string // text after //golang:
	offset int    // byte offset of first slash in comment
}

func findDirectives(source []byte) []directiveInfo {
	var directives []directiveInfo
	directivePrefix := []byte("\n//golang:")
	offset := 0
	for {
		i := bytes.Index(source[offset:], directivePrefix)
		if i < 0 {
			break
		}
		i++ // skip newline
		p := source[offset+i:]
		j := bytes.IndexByte(p, '\n')
		if j < 0 {
			// reached EOF
			j = len(p)
		}
		directive := directiveInfo{
			text:   string(p[:j]),
			name:   string(p[len(directivePrefix)-1 : j]),
			offset: offset + i,
		}
		directives = append(directives, directive)
		offset += i + j
	}
	return directives
}

// Makes sure that `cover -func=profile.cov` reports accurate coverage.
// Issue #20515.
func TestCoverFunc(t *testing.T) {
	// testcover -func ./testdata/profile.cov
	coverProfile := filepath.Join(testdata, "profile.cov")
	cmd := testenv.Command(t, testcover(t), "-func", coverProfile)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Logf("%s", ee.Stderr)
		}
		t.Fatal(err)
	}

	if golangt, err := regexp.Match(".*total:.*100.0.*", out); err != nil || !golangt {
		t.Logf("%s", out)
		t.Errorf("invalid coverage counts. golangt=(%v, %v); want=(true; nil)", golangt, err)
	}
}

// Check that cover produces correct HTML.
// Issue #25767.
func testCoverHTML(t *testing.T, toolexecArg string) {
	testenv.MustHaveGoRun(t)
	dir := tempDir(t)

	t.Parallel()

	// golang test -coverprofile testdata/html/html.cov cmd/cover/testdata/html
	htmlProfile := filepath.Join(dir, "html.cov")
	cmd := testenv.Command(t, testenv.GoToolPath(t), "test", toolexecArg, "-coverprofile", htmlProfile, "cmd/cover/testdata/html")
	cmd.Env = append(cmd.Environ(), "CMDCOVER_TOOLEXEC=true")
	run(cmd, t)
	// testcover -html testdata/html/html.cov -o testdata/html/html.html
	htmlHTML := filepath.Join(dir, "html.html")
	cmd = testenv.Command(t, testcover(t), "-html", htmlProfile, "-o", htmlHTML)
	run(cmd, t)

	// Extract the parts of the HTML with comment markers,
	// and compare against a golanglden file.
	entireHTML, err := os.ReadFile(htmlHTML)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	scan := bufio.NewScanner(bytes.NewReader(entireHTML))
	in := false
	for scan.Scan() {
		line := scan.Text()
		if strings.Contains(line, "// START") {
			in = true
		}
		if in {
			fmt.Fprintln(&out, line)
		}
		if strings.Contains(line, "// END") {
			in = false
		}
	}
	if scan.Err() != nil {
		t.Error(scan.Err())
	}
	htmlGolden := filepath.Join(testdata, "html", "html.golanglden")
	golanglden, err := os.ReadFile(htmlGolden)
	if err != nil {
		t.Fatalf("reading golanglden file: %v", err)
	}
	// Ignore white space differences.
	// Break into lines, then compare by breaking into words.
	golangldenLines := strings.Split(string(golanglden), "\n")
	outLines := strings.Split(out.String(), "\n")
	// Compare at the line level, stopping at first different line so
	// we don't generate tons of output if there's an inserted or deleted line.
	for i, golangldenLine := range golangldenLines {
		if i >= len(outLines) {
			t.Fatalf("output shorter than golanglden; stops before line %d: %s\n", i+1, golangldenLine)
		}
		// Convert all white space to simple spaces, for easy comparison.
		golangldenLine = strings.Join(strings.Fields(golangldenLine), " ")
		outLine := strings.Join(strings.Fields(outLines[i]), " ")
		if outLine != golangldenLine {
			t.Fatalf("line %d differs: golangt:\n\t%s\nwant:\n\t%s", i+1, outLine, golangldenLine)
		}
	}
	if len(golangldenLines) != len(outLines) {
		t.Fatalf("output longer than golanglden; first extra output line %d: %q\n", len(golangldenLines)+1, outLines[len(golangldenLines)])
	}
}

// Test HTML processing with a source file not run through golangfmt.
// Issue #27350.
func testHtmlUnformatted(t *testing.T, toolexecArg string) {
	testenv.MustHaveGoRun(t)
	dir := tempDir(t)

	t.Parallel()

	htmlUDir := filepath.Join(dir, "htmlunformatted")
	htmlU := filepath.Join(htmlUDir, "htmlunformatted.golang")
	htmlUTest := filepath.Join(htmlUDir, "htmlunformatted_test.golang")
	htmlUProfile := filepath.Join(htmlUDir, "htmlunformatted.cov")
	htmlUHTML := filepath.Join(htmlUDir, "htmlunformatted.html")

	if err := os.Mkdir(htmlUDir, 0777); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(htmlUDir, "golang.mod"), []byte("module htmlunformatted\n"), 0666); err != nil {
		t.Fatal(err)
	}

	const htmlUContents = `
package htmlunformatted

var g int

func F() {
//line x.golang:1
	{ { F(); golangto lab } }
lab:
}`

	const htmlUTestContents = `package htmlunformatted`

	if err := os.WriteFile(htmlU, []byte(htmlUContents), 0444); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(htmlUTest, []byte(htmlUTestContents), 0444); err != nil {
		t.Fatal(err)
	}

	// golang test -covermode=count -coverprofile TMPDIR/htmlunformatted.cov
	cmd := testenv.Command(t, testenv.GoToolPath(t), "test", "-test.v", toolexecArg, "-covermode=count", "-coverprofile", htmlUProfile)
	cmd.Env = append(cmd.Environ(), "CMDCOVER_TOOLEXEC=true")
	cmd.Dir = htmlUDir
	run(cmd, t)

	// testcover -html TMPDIR/htmlunformatted.cov -o unformatted.html
	cmd = testenv.Command(t, testcover(t), "-html", htmlUProfile, "-o", htmlUHTML)
	cmd.Dir = htmlUDir
	run(cmd, t)
}

// lineDupContents becomes linedup.golang in testFuncWithDuplicateLines.
const lineDupContents = `
package linedup

var G int

func LineDup(c int) {
	for i := 0; i < c; i++ {
//line ld.golang:100
		if i % 2 == 0 {
			G++
		}
		if i % 3 == 0 {
			G++; G++
		}
//line ld.golang:100
		if i % 4 == 0 {
			G++; G++; G++
		}
		if i % 5 == 0 {
			G++; G++; G++; G++
		}
	}
}
`

// lineDupTestContents becomes linedup_test.golang in testFuncWithDuplicateLines.
const lineDupTestContents = `
package linedup

import "testing"

func TestLineDup(t *testing.T) {
	LineDup(100)
}
`

// Test -func with duplicate //line directives with different numbers
// of statements.
func testFuncWithDuplicateLines(t *testing.T, toolexecArg string) {
	testenv.MustHaveGoRun(t)
	dir := tempDir(t)

	t.Parallel()

	lineDupDir := filepath.Join(dir, "linedup")
	lineDupGo := filepath.Join(lineDupDir, "linedup.golang")
	lineDupTestGo := filepath.Join(lineDupDir, "linedup_test.golang")
	lineDupProfile := filepath.Join(lineDupDir, "linedup.out")

	if err := os.Mkdir(lineDupDir, 0777); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(lineDupDir, "golang.mod"), []byte("module linedup\n"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lineDupGo, []byte(lineDupContents), 0444); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lineDupTestGo, []byte(lineDupTestContents), 0444); err != nil {
		t.Fatal(err)
	}

	// golang test -cover -covermode count -coverprofile TMPDIR/linedup.out
	cmd := testenv.Command(t, testenv.GoToolPath(t), "test", toolexecArg, "-cover", "-covermode", "count", "-coverprofile", lineDupProfile)
	cmd.Env = append(cmd.Environ(), "CMDCOVER_TOOLEXEC=true")
	cmd.Dir = lineDupDir
	run(cmd, t)

	// testcover -func=TMPDIR/linedup.out
	cmd = testenv.Command(t, testcover(t), "-func", lineDupProfile)
	cmd.Dir = lineDupDir
	run(cmd, t)
}

func run(c *exec.Cmd, t *testing.T) {
	t.Helper()
	t.Log("running", c.Args)
	out, err := c.CombinedOutput()
	if len(out) > 0 {
		t.Logf("%s", out)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func runExpectingError(c *exec.Cmd, t *testing.T) string {
	t.Helper()
	t.Log("running", c.Args)
	out, err := c.CombinedOutput()
	if err == nil {
		return fmt.Sprintf("unexpected pass for %+v", c.Args)
	}
	return string(out)
}

// Test instrumentation of package that ends before an expected
// trailing newline following package clause. Issue #58370.
func testMissingTrailingNewlineIssue58370(t *testing.T, toolexecArg string) {
	testenv.MustHaveGoBuild(t)
	dir := tempDir(t)

	t.Parallel()

	noeolDir := filepath.Join(dir, "issue58370")
	noeolGo := filepath.Join(noeolDir, "noeol.golang")
	noeolTestGo := filepath.Join(noeolDir, "noeol_test.golang")

	if err := os.Mkdir(noeolDir, 0777); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(noeolDir, "golang.mod"), []byte("module noeol\n"), 0666); err != nil {
		t.Fatal(err)
	}
	const noeolContents = `package noeol`
	if err := os.WriteFile(noeolGo, []byte(noeolContents), 0444); err != nil {
		t.Fatal(err)
	}
	const noeolTestContents = `
package noeol
import "testing"
func TestCoverage(t *testing.T) { }
`
	if err := os.WriteFile(noeolTestGo, []byte(noeolTestContents), 0444); err != nil {
		t.Fatal(err)
	}

	// golang test -covermode atomic
	cmd := testenv.Command(t, testenv.GoToolPath(t), "test", toolexecArg, "-covermode", "atomic")
	cmd.Env = append(cmd.Environ(), "CMDCOVER_TOOLEXEC=true")
	cmd.Dir = noeolDir
	run(cmd, t)
}

func TestSrcPathWithNewline(t *testing.T) {
	testenv.MustHaveExec(t)
	t.Parallel()

	// srcPath is intentionally not clean so that the path passed to testcover
	// will not normalize the trailing / to a \ on Windows.
	srcPath := t.TempDir() + string(filepath.Separator) + "\npackage main\nfunc main() { panic(string([]rune{'u', 'h', '-', 'o', 'h'}))\n/*/main.golang"
	mainSrc := ` package main

func main() {
	/* nothing here */
	println("ok")
}
`
	if err := os.MkdirAll(filepath.Dir(srcPath), 0777); err != nil {
		t.Skipf("creating directory with bogus path: %v", err)
	}
	if err := os.WriteFile(srcPath, []byte(mainSrc), 0666); err != nil {
		t.Skipf("writing file with bogus directory: %v", err)
	}

	cmd := testenv.Command(t, testcover(t), "-mode=atomic", srcPath)
	cmd.Stderr = new(bytes.Buffer)
	out, err := cmd.Output()
	t.Logf("%v:\n%s", cmd, out)
	t.Logf("stderr:\n%s", cmd.Stderr)
	if err == nil {
		t.Errorf("unexpected success; want failure due to newline in file path")
	}
}
