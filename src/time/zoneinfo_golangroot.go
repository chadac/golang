// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !ios && !android

package time

func golangrootZoneSource(golangroot string) (string, bool) {
	if golangroot == "" {
		return "", false
	}
	return golangroot + "/lib/time/zoneinfo.zip", true
}
