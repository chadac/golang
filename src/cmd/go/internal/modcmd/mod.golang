// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package modcmd implements the “golang mod” command.
package modcmd

import (
	"cmd/golang/internal/base"
)

var CmdMod = &base.Command{
	UsageLine: "golang mod",
	Short:     "module maintenance",
	Long: `Go mod provides access to operations on modules.

Note that support for modules is built into all the golang commands,
not just 'golang mod'. For example, day-to-day adding, removing, upgrading,
and downgrading of dependencies should be done using 'golang get'.
See 'golang help modules' for an overview of module functionality.
	`,

	Commands: []*base.Command{
		cmdDownload,
		cmdEdit,
		cmdGraph,
		cmdInit,
		cmdTidy,
		cmdVendor,
		cmdVerify,
		cmdWhy,
	},
}
