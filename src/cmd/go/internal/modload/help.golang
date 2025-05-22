// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package modload

import "cmd/golang/internal/base"

var HelpModules = &base.Command{
	UsageLine: "modules",
	Short:     "modules, module versions, and more",
	Long: `
Modules are how Go manages dependencies.

A module is a collection of packages that are released, versioned, and
distributed together. Modules may be downloaded directly from version control
repositories or from module proxy servers.

For a series of tutorials on modules, see
https://golanglang.org/doc/tutorial/create-module.

For a detailed reference on modules, see https://golanglang.org/ref/mod.

By default, the golang command may download modules from https://proxy.golanglang.org.
It may authenticate modules using the checksum database at
https://sum.golanglang.org. Both services are operated by the Go team at Google.
The privacy policies for these services are available at
https://proxy.golanglang.org/privacy and https://sum.golanglang.org/privacy,
respectively.

The golang command's download behavior may be configured using GOPROXY, GOSUMDB,
GOPRIVATE, and other environment variables. See 'golang help environment'
and https://golanglang.org/ref/mod#private-module-privacy for more information.
	`,
}

var HelpGoMod = &base.Command{
	UsageLine: "golang.mod",
	Short:     "the golang.mod file",
	Long: `
A module version is defined by a tree of source files, with a golang.mod
file in its root. When the golang command is run, it looks in the current
directory and then successive parent directories to find the golang.mod
marking the root of the main (current) module.

The golang.mod file format is described in detail at
https://golanglang.org/ref/mod#golang-mod-file.

To create a new golang.mod file, use 'golang mod init'. For details see
'golang help mod init' or https://golanglang.org/ref/mod#golang-mod-init.

To add missing module requirements or remove unneeded requirements,
use 'golang mod tidy'. For details, see 'golang help mod tidy' or
https://golanglang.org/ref/mod#golang-mod-tidy.

To add, upgrade, downgrade, or remove a specific module requirement, use
'golang get'. For details, see 'golang help module-get' or
https://golanglang.org/ref/mod#golang-get.

To make other changes or to parse golang.mod as JSON for use by other tools,
use 'golang mod edit'. See 'golang help mod edit' or
https://golanglang.org/ref/mod#golang-mod-edit.
	`,
}
