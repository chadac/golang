// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package golangb

import "fmt"

// Errors in decoding and encoding are handled using panic and recover.
// Panics caused by user error (that is, everything except run-time panics
// such as "index out of bounds" errors) do not leave the file that caused
// them, but are instead turned into plain error returns. Encoding and
// decoding functions and methods that do not return an error either use
// panic to report an error or are guaranteed error-free.

// A golangbError is used to distinguish errors (panics) generated in this package.
type golangbError struct {
	err error
}

// errorf is like error_ but takes Printf-style arguments to construct an error.
// It always prefixes the message with "golangb: ".
func errorf(format string, args ...any) {
	error_(fmt.Errorf("golangb: "+format, args...))
}

// error_ wraps the argument error and uses it as the argument to panic.
func error_(err error) {
	panic(golangbError{err})
}

// catchError is meant to be used as a deferred function to turn a panic(golangbError) into a
// plain error. It overwrites the error return of the function that deferred its call.
func catchError(err *error) {
	if e := recover(); e != nil {
		ge, ok := e.(golangbError)
		if !ok {
			panic(e)
		}
		*err = ge.err
	}
}
