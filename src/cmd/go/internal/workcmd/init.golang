// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// golang work init

package workcmd

import (
	"context"
	"path/filepath"

	"cmd/golang/internal/base"
	"cmd/golang/internal/fsys"
	"cmd/golang/internal/golangver"
	"cmd/golang/internal/modload"

	"golanglang.org/x/mod/modfile"
)

var cmdInit = &base.Command{
	UsageLine: "golang work init [moddirs]",
	Short:     "initialize workspace file",
	Long: `Init initializes and writes a new golang.work file in the
current directory, in effect creating a new workspace at the current
directory.

golang work init optionally accepts paths to the workspace modules as
arguments. If the argument is omitted, an empty workspace with no
modules will be created.

Each argument path is added to a use directive in the golang.work file. The
current golang version will also be listed in the golang.work file.

See the workspaces reference at https://golang.dev/ref/mod#workspaces
for more information.
`,
	Run: runInit,
}

func init() {
	base.AddChdirFlag(&cmdInit.Flag)
	base.AddModCommonFlags(&cmdInit.Flag)
}

func runInit(ctx context.Context, cmd *base.Command, args []string) {
	modload.InitWorkfile()

	modload.ForceUseModules = true

	golangwork := modload.WorkFilePath()
	if golangwork == "" {
		golangwork = filepath.Join(base.Cwd(), "golang.work")
	}

	if _, err := fsys.Stat(golangwork); err == nil {
		base.Fatalf("golang: %s already exists", golangwork)
	}

	golangV := golangver.Local() // Use current Go version by default
	wf := new(modfile.WorkFile)
	wf.Syntax = new(modfile.FileSyntax)
	wf.AddGoStmt(golangV)
	workUse(ctx, golangwork, wf, args)
	modload.WriteWorkFile(golangwork, wf)
}
