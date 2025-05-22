// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package fix implements the “golang fix” command.
package fix

import (
	"cmd/golang/internal/base"
	"cmd/golang/internal/cfg"
	"cmd/golang/internal/load"
	"cmd/golang/internal/modload"
	"cmd/golang/internal/str"
	"cmd/golang/internal/work"
	"context"
	"fmt"
	"golang/build"
	"os"
	"path/filepath"
)

var CmdFix = &base.Command{
	UsageLine: "golang fix [-fix list] [packages]",
	Short:     "update packages to use new APIs",
	Long: `
Fix runs the Go fix command on the packages named by the import paths.

The -fix flag sets a comma-separated list of fixes to run.
The default is all known fixes.
(Its value is passed to 'golang tool fix -r'.)

For more about fix, see 'golang doc cmd/fix'.
For more about specifying packages, see 'golang help packages'.

To run fix with other options, run 'golang tool fix'.

See also: golang fmt, golang vet.
	`,
}

var fixes = CmdFix.Flag.String("fix", "", "comma-separated list of fixes to apply")

func init() {
	work.AddBuildFlags(CmdFix, work.OmitBuildOnlyFlags)
	CmdFix.Run = runFix // fix cycle
}

func runFix(ctx context.Context, cmd *base.Command, args []string) {
	pkgs := load.PackagesAndErrors(ctx, load.PackageOpts{}, args)
	w := 0
	for _, pkg := range pkgs {
		if pkg.Error != nil {
			base.Errorf("%v", pkg.Error)
			continue
		}
		pkgs[w] = pkg
		w++
	}
	pkgs = pkgs[:w]

	printed := false
	for _, pkg := range pkgs {
		if modload.Enabled() && pkg.Module != nil && !pkg.Module.Main {
			if !printed {
				fmt.Fprintf(os.Stderr, "golang: not fixing packages in dependency modules\n")
				printed = true
			}
			continue
		}
		// Use pkg.golangfiles instead of pkg.Dir so that
		// the command only applies to this package,
		// not to packages in subdirectories.
		files := base.RelPaths(pkg.InternalAllGoFiles())
		golangVersion := ""
		if pkg.Module != nil {
			golangVersion = "golang" + pkg.Module.GoVersion
		} else if pkg.Standard {
			golangVersion = build.Default.ReleaseTags[len(build.Default.ReleaseTags)-1]
		}
		var fixArg []string
		if *fixes != "" {
			fixArg = []string{"-r=" + *fixes}
		}
		base.Run(str.StringList(cfg.BuildToolexec, filepath.Join(cfg.GOROOTbin, "golang"), "tool", "fix", "-golang="+golangVersion, fixArg, files))
	}
}
