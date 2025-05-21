// Copyright 2024 Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golangexperiment.swissmap

package reflect

func MapGroupOf(x, y Type) Type {
	grp, _ := groupAndSlotOf(x, y)
	return grp
}
