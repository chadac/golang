// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// The runtime package contains an uninitialized definition
// for runtime·iscgolang. Override it to tell the runtime we're here.
// There are various function pointers that should be set too,
// but those depend on dynamic linker magic to get initialized
// correctly, and sometimes they break. This variable is a
// backup: it depends only on old C style static linking rules.

package cgolang

import _ "unsafe" // for golang:linkname

//golang:linkname _iscgolang runtime.iscgolang
var _iscgolang bool = true
