// errorcheckdir

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// The gccgolang compiler had a bug: if one file in a package did a dot
// import, then an earlier file in the package would incorrectly
// resolve to the imported names rather than reporting undefined
// errors.

package ignored
