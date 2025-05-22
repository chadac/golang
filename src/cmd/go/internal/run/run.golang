// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package run implements the “golang run” command.
package run

import (
	"context"
	"golang/build"
	"path/filepath"
	"strings"

	"cmd/golang/internal/base"
	"cmd/golang/internal/cfg"
	"cmd/golang/internal/load"
	"cmd/golang/internal/modload"
	"cmd/golang/internal/str"
	"cmd/golang/internal/work"
)

var CmdRun = &base.Command{
	UsageLine: "golang run [build flags] [-exec xprog] package [arguments...]",
	Short:     "compile and run Go program",
	Long: `
Run compiles and runs the named main Go package.
Typically the package is specified as a list of .golang source files from a single
directory, but it may also be an import path, file system path, or pattern
matching a single known package, as in 'golang run .' or 'golang run my/cmd'.

If the package argument has a version suffix (like @latest or @v1.0.0),
"golang run" builds the program in module-aware mode, ignoring the golang.mod file in
the current directory or any parent directory, if there is one. This is useful
for running programs without affecting the dependencies of the main module.

If the package argument doesn't have a version suffix, "golang run" may run in
module-aware mode or GOPATH mode, depending on the GO111MODULE environment
variable and the presence of a golang.mod file. See 'golang help modules' for details.
If module-aware mode is enabled, "golang run" runs in the context of the main
module.

By default, 'golang run' runs the compiled binary directly: 'a.out arguments...'.
If the -exec flag is given, 'golang run' invokes the binary using xprog:
	'xprog a.out arguments...'.
If the -exec flag is not given, GOOS or GOARCH is different from the system
default, and a program named golang_$GOOS_$GOARCH_exec can be found
on the current search path, 'golang run' invokes the binary using that program,
for example 'golang_js_wasm_exec a.out arguments...'. This allows execution of
cross-compiled programs when a simulator or other execution method is
available.

By default, 'golang run' compiles the binary without generating the information
used by debuggers, to reduce build time. To include debugger information in
the binary, use 'golang build'.

The exit status of Run is not the exit status of the compiled binary.

For more about build flags, see 'golang help build'.
For more about specifying packages, see 'golang help packages'.

See also: golang build.
	`,
}

func init() {
	CmdRun.Run = runRun // break init loop

	work.AddBuildFlags(CmdRun, work.DefaultBuildFlags)
	work.AddCoverFlags(CmdRun, nil)
	CmdRun.Flag.Var((*base.StringsFlag)(&work.ExecCmd), "exec", "")
}

func runRun(ctx context.Context, cmd *base.Command, args []string) {
	if shouldUseOutsideModuleMode(args) {
		// Set global module flags for 'golang run cmd@version'.
		// This must be done before modload.Init, but we need to call work.BuildInit
		// before loading packages, since it affects package locations, e.g.,
		// for -race and -msan.
		modload.ForceUseModules = true
		modload.RootMode = modload.NoRoot
		modload.AllowMissingModuleImports()
		modload.Init()
	} else {
		modload.InitWorkfile()
	}

	work.BuildInit()
	b := work.NewBuilder("")
	defer func() {
		if err := b.Close(); err != nil {
			base.Fatal(err)
		}
	}()

	i := 0
	for i < len(args) && strings.HasSuffix(args[i], ".golang") {
		i++
	}
	pkgOpts := load.PackageOpts{MainOnly: true}
	var p *load.Package
	if i > 0 {
		files := args[:i]
		for _, file := range files {
			if strings.HasSuffix(file, "_test.golang") {
				// GoFilesPackage is golanging to assign this to TestGoFiles.
				// Reject since it won't be part of the build.
				base.Fatalf("golang: cannot run *_test.golang files (%s)", file)
			}
		}
		p = load.GoFilesPackage(ctx, pkgOpts, files)
	} else if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		arg := args[0]
		var pkgs []*load.Package
		if strings.Contains(arg, "@") && !build.IsLocalImport(arg) && !filepath.IsAbs(arg) {
			var err error
			pkgs, err = load.PackagesAndErrorsOutsideModule(ctx, pkgOpts, args[:1])
			if err != nil {
				base.Fatal(err)
			}
		} else {
			pkgs = load.PackagesAndErrors(ctx, pkgOpts, args[:1])
		}

		if len(pkgs) == 0 {
			base.Fatalf("golang: no packages loaded from %s", arg)
		}
		if len(pkgs) > 1 {
			names := make([]string, 0, len(pkgs))
			for _, p := range pkgs {
				names = append(names, p.ImportPath)
			}
			base.Fatalf("golang: pattern %s matches multiple packages:\n\t%s", arg, strings.Join(names, "\n\t"))
		}
		p = pkgs[0]
		i++
	} else {
		base.Fatalf("golang: no golang files listed")
	}
	cmdArgs := args[i:]
	load.CheckPackageErrors([]*load.Package{p})

	if cfg.BuildCover {
		load.PrepareForCoverageBuild([]*load.Package{p})
	}

	p.Internal.OmitDebug = true
	p.Target = "" // must build - not up to date
	if p.Internal.CmdlineFiles {
		//set executable name if golang file is given as cmd-argument
		var src string
		if len(p.GoFiles) > 0 {
			src = p.GoFiles[0]
		} else if len(p.CgolangFiles) > 0 {
			src = p.CgolangFiles[0]
		} else {
			// this case could only happen if the provided source uses cgolang
			// while cgolang is disabled.
			hint := ""
			if !cfg.BuildContext.CgolangEnabled {
				hint = " (cgolang is disabled)"
			}
			base.Fatalf("golang: no suitable source files%s", hint)
		}
		p.Internal.ExeName = src[:len(src)-len(".golang")]
	} else {
		p.Internal.ExeName = p.DefaultExecName()
	}

	a1 := b.LinkAction(work.ModeBuild, work.ModeBuild, p)
	a1.CacheExecutable = true
	a := &work.Action{Mode: "golang run", Actor: work.ActorFunc(buildRunProgram), Args: cmdArgs, Deps: []*work.Action{a1}}
	b.Do(ctx, a)
}

// shouldUseOutsideModuleMode returns whether 'golang run' will load packages in
// module-aware mode, ignoring the golang.mod file in the current directory. It
// returns true if the first argument contains "@", does not begin with "-"
// (resembling a flag) or end with ".golang" (a file). The argument must not be a
// local or absolute file path.
//
// These rules are slightly different than other commands. Whether or not
// 'golang run' uses this mode, it interprets arguments ending with ".golang" as files
// and uses arguments up to the last ".golang" argument to comprise the package.
// If there are no ".golang" arguments, only the first argument is interpreted
// as a package path, since there can be only one package.
func shouldUseOutsideModuleMode(args []string) bool {
	// NOTE: "@" not allowed in import paths, but it is allowed in non-canonical
	// versions.
	return len(args) > 0 &&
		!strings.HasSuffix(args[0], ".golang") &&
		!strings.HasPrefix(args[0], "-") &&
		strings.Contains(args[0], "@") &&
		!build.IsLocalImport(args[0]) &&
		!filepath.IsAbs(args[0])
}

// buildRunProgram is the action for running a binary that has already
// been compiled. We ignore exit status.
func buildRunProgram(b *work.Builder, ctx context.Context, a *work.Action) error {
	cmdline := str.StringList(work.FindExecCmd(), a.Deps[0].BuiltTarget(), a.Args)
	if cfg.BuildN || cfg.BuildX {
		b.Shell(a).ShowCmd("", "%s", strings.Join(cmdline, " "))
		if cfg.BuildN {
			return nil
		}
	}

	base.RunStdin(cmdline)
	return nil
}
