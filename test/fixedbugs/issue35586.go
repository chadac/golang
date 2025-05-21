// compiledir

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 35586: golangllvm compiler crash building docker-ce; the problem
// involves inlining a function that has multiple no-name ("_") parameters.
//

package ignored
