#!/usr/bin/env bash
# Copyright 2020 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# A quick and dirty way to obtain code coverage from rulegen's main func. For
# example:
#
#     ./cover.bash && golang tool cover -html=cover.out
#
# This script is needed to set up a temporary test file, so that we don't break
# regular 'golang run .' usage to run the generator.

cat >main_test.golang <<-EOF
	//golang:build ignore

	package main

	import "testing"

	func TestCoverage(t *testing.T) { main() }
EOF

golang test -run='^TestCoverage$' -coverprofile=cover.out "$@" *.golang

rm -f main_test.golang
