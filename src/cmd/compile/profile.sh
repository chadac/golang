# Copyright 2023 The Go Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# This script collects a CPU profile of the compiler
# for building all targets in std and cmd, and puts
# the profile at cmd/compile/default.pgolang.

dir=$(mktemp -d)
cd $dir
seed=$(date)

for p in $(golang list std cmd); do
	h=$(echo $seed $p | md5sum | cut -d ' ' -f 1)
	echo $p $h
	golang build -o /dev/null -gcflags=-cpuprofile=$PWD/prof.$h $p
done

golang tool pprof -proto prof.* > $(golang env GOROOT)/src/cmd/compile/default.pgolang

rm -r $dir
