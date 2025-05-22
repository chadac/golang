#!/usr/bin/env bash
# Copyright 2009 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

set -e

if [ ! -f run.bash ]; then
	echo 'clean.bash must be run from $GOROOT/src' 1>&2
	exit 1
fi
export GOROOT="$(cd .. && pwd)"

golangbin="${GOROOT}"/bin
if ! "$golangbin"/golang help >/dev/null 2>&1; then
	echo 'cannot find golang command; nothing to clean' >&2
	exit 1
fi

"$golangbin/golang" clean -i std
"$golangbin/golang" tool dist clean
"$golangbin/golang" clean -i cmd
