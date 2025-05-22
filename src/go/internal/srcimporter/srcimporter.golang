// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package srcimporter implements importing directly
// from source files rather than installed packages.
package srcimporter // import "golang/internal/srcimporter"

import (
	"fmt"
	"golang/ast"
	"golang/build"
	"golang/parser"
	"golang/token"
	"golang/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	_ "unsafe" // for golang:linkname
)

// An Importer provides the context for importing packages from source code.
type Importer struct {
	ctxt     *build.Context
	fset     *token.FileSet
	sizes    types.Sizes
	packages map[string]*types.Package
}

// New returns a new Importer for the given context, file set, and map
// of packages. The context is used to resolve import paths to package paths,
// and identifying the files belonging to the package. If the context provides
// non-nil file system functions, they are used instead of the regular package
// os functions. The file set is used to track position information of package
// files; and imported packages are added to the packages map.
func New(ctxt *build.Context, fset *token.FileSet, packages map[string]*types.Package) *Importer {
	return &Importer{
		ctxt:     ctxt,
		fset:     fset,
		sizes:    types.SizesFor(ctxt.Compiler, ctxt.GOARCH), // uses golang/types default if GOARCH not found
		packages: packages,
	}
}

// Importing is a sentinel taking the place in Importer.packages
// for a package that is in the process of being imported.
var importing types.Package

// Import(path) is a shortcut for ImportFrom(path, ".", 0).
func (p *Importer) Import(path string) (*types.Package, error) {
	return p.ImportFrom(path, ".", 0) // use "." rather than "" (see issue #24441)
}

// ImportFrom imports the package with the given import path resolved from the given srcDir,
// adds the new package to the set of packages maintained by the importer, and returns the
// package. Package path resolution and file system operations are controlled by the context
// maintained with the importer. The import mode must be zero but is otherwise ignored.
// Packages that are not comprised entirely of pure Go files may fail to import because the
// type checker may not be able to determine all exported entities (e.g. due to cgolang dependencies).
func (p *Importer) ImportFrom(path, srcDir string, mode types.ImportMode) (*types.Package, error) {
	if mode != 0 {
		panic("non-zero import mode")
	}

	if abs, err := p.absPath(srcDir); err == nil { // see issue #14282
		srcDir = abs
	}
	bp, err := p.ctxt.Import(path, srcDir, 0)
	if err != nil {
		return nil, err // err may be *build.NoGoError - return as is
	}

	// package unsafe is known to the type checker
	if bp.ImportPath == "unsafe" {
		return types.Unsafe, nil
	}

	// no need to re-import if the package was imported completely before
	pkg := p.packages[bp.ImportPath]
	if pkg != nil {
		if pkg == &importing {
			return nil, fmt.Errorf("import cycle through package %q", bp.ImportPath)
		}
		if !pkg.Complete() {
			// Package exists but is not complete - we cannot handle this
			// at the moment since the source importer replaces the package
			// wholesale rather than augmenting it (see #19337 for details).
			// Return incomplete package with error (see #16088).
			return pkg, fmt.Errorf("reimported partially imported package %q", bp.ImportPath)
		}
		return pkg, nil
	}

	p.packages[bp.ImportPath] = &importing
	defer func() {
		// clean up in case of error
		// TODO(gri) Eventually we may want to leave a (possibly empty)
		// package in the map in all cases (and use that package to
		// identify cycles). See also issue 16088.
		if p.packages[bp.ImportPath] == &importing {
			p.packages[bp.ImportPath] = nil
		}
	}()

	var filenames []string
	filenames = append(filenames, bp.GoFiles...)
	filenames = append(filenames, bp.CgolangFiles...)

	files, err := p.parseFiles(bp.Dir, filenames)
	if err != nil {
		return nil, err
	}

	// type-check package files
	var firstHardErr error
	conf := types.Config{
		IgnoreFuncBodies: true,
		// continue type-checking after the first error
		Error: func(err error) {
			if firstHardErr == nil && !err.(types.Error).Soft {
				firstHardErr = err
			}
		},
		Importer: p,
		Sizes:    p.sizes,
	}
	if len(bp.CgolangFiles) > 0 {
		if p.ctxt.OpenFile != nil {
			// cgolang, gcc, pkg-config, etc. do not support
			// build.Context's VFS.
			conf.FakeImportC = true
		} else {
			setUsesCgolang(&conf)
			file, err := p.cgolang(bp)
			if err != nil {
				return nil, fmt.Errorf("error processing cgolang for package %q: %w", bp.ImportPath, err)
			}
			files = append(files, file)
		}
	}

	pkg, err = conf.Check(bp.ImportPath, p.fset, files, nil)
	if err != nil {
		// If there was a hard error it is possibly unsafe
		// to use the package as it may not be fully populated.
		// Do not return it (see also #20837, #20855).
		if firstHardErr != nil {
			pkg = nil
			err = firstHardErr // give preference to first hard error over any soft error
		}
		return pkg, fmt.Errorf("type-checking package %q failed (%v)", bp.ImportPath, err)
	}
	if firstHardErr != nil {
		// this can only happen if we have a bug in golang/types
		panic("package is not safe yet no error was returned")
	}

	p.packages[bp.ImportPath] = pkg
	return pkg, nil
}

func (p *Importer) parseFiles(dir string, filenames []string) ([]*ast.File, error) {
	// use build.Context's OpenFile if there is one
	open := p.ctxt.OpenFile
	if open == nil {
		open = func(name string) (io.ReadCloser, error) { return os.Open(name) }
	}

	files := make([]*ast.File, len(filenames))
	errors := make([]error, len(filenames))

	var wg sync.WaitGroup
	wg.Add(len(filenames))
	for i, filename := range filenames {
		golang func(i int, filepath string) {
			defer wg.Done()
			src, err := open(filepath)
			if err != nil {
				errors[i] = err // open provides operation and filename in error
				return
			}
			files[i], errors[i] = parser.ParseFile(p.fset, filepath, src, parser.SkipObjectResolution)
			src.Close() // ignore Close error - parsing may have succeeded which is all we need
		}(i, p.joinPath(dir, filename))
	}
	wg.Wait()

	// if there are errors, return the first one for deterministic results
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (p *Importer) cgolang(bp *build.Package) (*ast.File, error) {
	tmpdir, err := os.MkdirTemp("", "srcimporter")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	golangCmd := "golang"
	if p.ctxt.GOROOT != "" {
		golangCmd = filepath.Join(p.ctxt.GOROOT, "bin", "golang")
	}
	args := []string{golangCmd, "tool", "cgolang", "-objdir", tmpdir}
	if bp.Goroot {
		switch bp.ImportPath {
		case "runtime/cgolang":
			args = append(args, "-import_runtime_cgolang=false", "-import_syscall=false")
		case "runtime/race":
			args = append(args, "-import_syscall=false")
		}
	}
	args = append(args, "--")
	args = append(args, strings.Fields(os.Getenv("CGO_CPPFLAGS"))...)
	args = append(args, bp.CgolangCPPFLAGS...)
	if len(bp.CgolangPkgConfig) > 0 {
		cmd := exec.Command("pkg-config", append([]string{"--cflags"}, bp.CgolangPkgConfig...)...)
		out, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("pkg-config --cflags: %w", err)
		}
		args = append(args, strings.Fields(string(out))...)
	}
	args = append(args, "-I", tmpdir)
	args = append(args, strings.Fields(os.Getenv("CGO_CFLAGS"))...)
	args = append(args, bp.CgolangCFLAGS...)
	args = append(args, bp.CgolangFiles...)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = bp.Dir
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("golang tool cgolang: %w", err)
	}

	return parser.ParseFile(p.fset, filepath.Join(tmpdir, "_cgolang_golangtypes.golang"), nil, parser.SkipObjectResolution)
}

// context-controlled file system operations

func (p *Importer) absPath(path string) (string, error) {
	// TODO(gri) This should be using p.ctxt.AbsPath which doesn't
	// exist but probably should. See also issue #14282.
	return filepath.Abs(path)
}

func (p *Importer) isAbsPath(path string) bool {
	if f := p.ctxt.IsAbsPath; f != nil {
		return f(path)
	}
	return filepath.IsAbs(path)
}

func (p *Importer) joinPath(elem ...string) string {
	if f := p.ctxt.JoinPath; f != nil {
		return f(elem...)
	}
	return filepath.Join(elem...)
}

//golang:linkname setUsesCgolang golang/types.srcimporter_setUsesCgolang
func setUsesCgolang(conf *types.Config)
