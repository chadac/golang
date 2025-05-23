// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build boringcrypto

package main_test

import "testing"

func TestBoringInternalLink(t *testing.T) {
	tg := testgolang(t)
	defer tg.cleanup()
	tg.parallel()
	tg.tempFile("main.golang", `package main
		import "crypto/sha1"
		func main() {
			sha1.New()
		}`)
	tg.run("build", "-ldflags=-w -extld=false", tg.path("main.golang"))
	tg.run("build", "-ldflags=-extld=false", tg.path("main.golang"))
}
