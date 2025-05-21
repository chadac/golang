// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// golang mod graph

package modcmd

import (
	"bufio"
	"context"
	"os"

	"cmd/golang/internal/base"
	"cmd/golang/internal/cfg"
	"cmd/golang/internal/golangver"
	"cmd/golang/internal/modload"
	"cmd/golang/internal/toolchain"

	"golanglang.org/x/mod/module"
)

var cmdGraph = &base.Command{
	UsageLine: "golang mod graph [-golang=version] [-x]",
	Short:     "print module requirement graph",
	Long: `
Graph prints the module requirement graph (with replacements applied)
in text form. Each line in the output has two space-separated fields: a module
and one of its requirements. Each module is identified as a string of the form
path@version, except for the main module, which has no @version suffix.

The -golang flag causes graph to report the module graph as loaded by the
given Go version, instead of the version indicated by the 'golang' directive
in the golang.mod file.

The -x flag causes graph to print the commands graph executes.

See https://golanglang.org/ref/mod#golang-mod-graph for more about 'golang mod graph'.
	`,
	Run: runGraph,
}

var (
	graphGo golangVersionFlag
)

func init() {
	cmdGraph.Flag.Var(&graphGo, "golang", "")
	cmdGraph.Flag.BoolVar(&cfg.BuildX, "x", false, "")
	base.AddChdirFlag(&cmdGraph.Flag)
	base.AddModCommonFlags(&cmdGraph.Flag)
}

func runGraph(ctx context.Context, cmd *base.Command, args []string) {
	modload.InitWorkfile()

	if len(args) > 0 {
		base.Fatalf("golang: 'golang mod graph' accepts no arguments")
	}
	modload.ForceUseModules = true
	modload.RootMode = modload.NeedRoot

	golangVersion := graphGo.String()
	if golangVersion != "" && golangver.Compare(golangver.Local(), golangVersion) < 0 {
		toolchain.SwitchOrFatal(ctx, &golangver.TooNewError{
			What:      "-golang flag",
			GoVersion: golangVersion,
		})
	}

	mg, err := modload.LoadModGraph(ctx, golangVersion)
	if err != nil {
		base.Fatal(err)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	format := func(m module.Version) {
		w.WriteString(m.Path)
		if m.Version != "" {
			w.WriteString("@")
			w.WriteString(m.Version)
		}
	}

	mg.WalkBreadthFirst(func(m module.Version) {
		reqs, _ := mg.RequiredBy(m)
		for _, r := range reqs {
			format(m)
			w.WriteByte(' ')
			format(r)
			w.WriteByte('\n')
		}
	})
}
