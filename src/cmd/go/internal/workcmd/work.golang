// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package workcmd implements the “golang work” command.
package workcmd

import (
	"cmd/golang/internal/base"
)

var CmdWork = &base.Command{
	UsageLine: "golang work",
	Short:     "workspace maintenance",
	Long: `Work provides access to operations on workspaces.

Note that support for workspaces is built into many other commands, not
just 'golang work'.

See 'golang help modules' for information about Go's module system of which
workspaces are a part.

See https://golang.dev/ref/mod#workspaces for an in-depth reference on
workspaces.

See https://golang.dev/doc/tutorial/workspaces for an introductory
tutorial on workspaces.

A workspace is specified by a golang.work file that specifies a set of
module directories with the "use" directive. These modules are used as
root modules by the golang command for builds and related operations.  A
workspace that does not specify modules to be used cannot be used to do
builds from local modules.

golang.work files are line-oriented. Each line holds a single directive,
made up of a keyword followed by arguments. For example:

	golang 1.18

	use ../foo/bar
	use ./baz

	replace example.com/foo v1.2.3 => example.com/bar v1.4.5

The leading keyword can be factored out of adjacent lines to create a block,
like in Go imports.

	use (
	  ../foo/bar
	  ./baz
	)

The use directive specifies a module to be included in the workspace's
set of main modules. The argument to the use directive is the directory
containing the module's golang.mod file.

The golang directive specifies the version of Go the file was written at. It
is possible there may be future changes in the semantics of workspaces
that could be controlled by this version, but for now the version
specified has no effect.

The replace directive has the same syntax as the replace directive in a
golang.mod file and takes precedence over replaces in golang.mod files.  It is
primarily intended to override conflicting replaces in different workspace
modules.

To determine whether the golang command is operating in workspace mode, use
the "golang env GOWORK" command. This will specify the workspace file being
used.
`,

	Commands: []*base.Command{
		cmdEdit,
		cmdInit,
		cmdSync,
		cmdUse,
		cmdVendor,
	},
}
