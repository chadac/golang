#!/bin/bash
# Copyright 2022 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# This could be a golangod use for embed but golang/doc/comment
# is built into the bootstrap golang command, so it can't use embed.
# Also not using embed lets us emit a string array directly
# and avoid init-time work.

(
echo "// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by 'golang generate' DO NOT EDIT.
//golang:generate ./mkstd.sh

package comment

var stdPkgs = []string{"
golang list std | grep -v / | sort | sed 's/.*/"&",/'
echo "}"
) | golangfmt >std.golang.tmp && mv std.golang.tmp std.golang
