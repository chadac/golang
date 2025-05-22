#!/usr/bin/env bash
# Copyright 2015 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# When run as (for example)
#
#	GOOS=linux GOARCH=ppc64 bootstrap.bash
#
# this script cross-compiles a toolchain for that GOOS/GOARCH
# combination, leaving the resulting tree in ../../golang-${GOOS}-${GOARCH}-bootstrap.
# That tree can be copied to a machine of the given target type
# and used as $GOROOT_BOOTSTRAP to bootstrap a local build.
#
# Only changes that have been committed to Git (at least locally,
# not necessary reviewed and submitted to master) are included in the tree.
#
# See also golanglang.org/x/build/cmd/genbootstrap, which is used
# to generate bootstrap tgz files for builders.

set -e

if [ "$GOOS" = "" -o "$GOARCH" = "" ]; then
	echo "usage: GOOS=os GOARCH=arch ./bootstrap.bash [-force]" >&2
	exit 2
fi

forceflag=""
if [ "$1" = "-force" ]; then
	forceflag=-force
	shift
fi

targ="../../golang-${GOOS}-${GOARCH}-bootstrap"
if [ -e $targ ]; then
	echo "$targ already exists; remove before continuing"
	exit 2
fi

unset GOROOT
src=$(cd .. && pwd)
echo "#### Copying to $targ"
cp -Rp "$src" "$targ"
cd "$targ"
echo
echo "#### Cleaning $targ"
chmod -R +w .
rm -f .gitignore
if [ -e .git ]; then
	git clean -f -d
fi
echo
echo "#### Building $targ"
echo
cd src
./make.bash --no-banner $forceflag
golanghostos="$(../bin/golang env GOHOSTOS)"
golanghostarch="$(../bin/golang env GOHOSTARCH)"
golangos="$(../bin/golang env GOOS)"
golangarch="$(../bin/golang env GOARCH)"

# NOTE: Cannot invoke golang command after this point.
# We're about to delete all but the cross-compiled binaries.
cd ..
if [ "$golangos" = "$golanghostos" -a "$golangarch" = "$golanghostarch" ]; then
	# cross-compile for local system. nothing to copy.
	# useful if you've bootstrapped yourself but want to
	# prepare a clean toolchain for others.
	true
else
	rm -f bin/golang_${golangos}_${golangarch}_exec
	mv bin/*_*/* bin
	rmdir bin/*_*
	rm -rf "pkg/${golanghostos}_${golanghostarch}" "pkg/tool/${golanghostos}_${golanghostarch}"
fi

rm -rf pkg/bootstrap pkg/obj .git

echo ----
echo Bootstrap toolchain for "$GOOS/$GOARCH" installed in "$(pwd)".
echo Building tbz.
cd ..
tar cf - "golang-${GOOS}-${GOARCH}-bootstrap" | bzip2 -9 >"golang-${GOOS}-${GOARCH}-bootstrap.tbz"
ls -l "$(pwd)/golang-${GOOS}-${GOARCH}-bootstrap.tbz"
exit 0
