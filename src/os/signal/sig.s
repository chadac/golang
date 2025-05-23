// Copyright 2012 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// The runtime package uses //golang:linkname to push a few functions into this
// package but we still need a .s file so the Golang tool does not pass -complete
// to the golang tool compile so the latter does not complain about Golang functions
// with no bodies.
