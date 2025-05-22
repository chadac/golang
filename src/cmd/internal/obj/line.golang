// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package obj

import (
	"cmd/internal/golangobj"
	"cmd/internal/src"
)

// AddImport adds a package to the list of imported packages.
func (ctxt *Link) AddImport(pkg string, fingerprint golangobj.FingerprintType) {
	ctxt.Imports = append(ctxt.Imports, golangobj.ImportedPkg{Pkg: pkg, Fingerprint: fingerprint})
}

// getFileIndexAndLine returns the relative file index (local to the CU), and
// the relative line number for a position (i.e., as adjusted by a //line
// directive). This is the file/line visible in the final binary (pcfile, pcln,
// etc).
func (ctxt *Link) getFileIndexAndLine(xpos src.XPos) (int, int32) {
	pos := ctxt.InnermostPos(xpos)
	if !pos.IsKnown() {
		pos = src.Pos{}
	}
	return pos.FileIndex(), int32(pos.RelLine())
}
