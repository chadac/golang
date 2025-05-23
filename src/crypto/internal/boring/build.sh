#!/bin/bash
# Copyright 2022 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# This shell script uses Docker to run build-boring.sh and build-golangboring.sh,
# which build golangboringcrypto_linux_$GOARCH.syso according to the Security Policy.
# Currently, amd64 and arm64 are permitted.

set -e
set -o pipefail

GOARCH=${GOARCH:-$(golang env GOARCH)}
echo "# Building golangboringcrypto_linux_$GOARCH.syso. Set GOARCH to override." >&2

if ! which docker >/dev/null; then
	echo "# Docker not found. Inside Golangogle, see golang/installdocker." >&2
	exit 1
fi

platform=""
buildargs=""
case "$GOARCH" in
amd64)
	if ! docker run --rm -t amd64/ubuntu:focal uname -m >/dev/null 2>&1; then
		echo "# Docker cannot run amd64 binaries."
		exit 1
	fi
	platform="--platform linux/amd64"
	buildargs="--build-arg ubuntu=amd64/ubuntu"
	;;
arm64)
	if ! docker run --rm -t arm64v8/ubuntu:focal uname -m >/dev/null 2>&1; then
		echo "# Docker cannot run arm64 binaries. Try:"
		echo "	sudo apt-get install qemu binfmt-support qemu-user-static"
		echo "	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes"
		echo "	docker run --rm -t arm64v8/ubuntu:focal uname -m"
		exit 1
	fi
	platform="--platform linux/arm64/v8"
	buildargs="--build-arg ubuntu=arm64v8/ubuntu"
	;;
*)
	echo unknown GOARCH $GOARCH >&2
	exit 2
esac

docker build $platform $buildargs --build-arg GOARCH=$GOARCH -t golangboring:$GOARCH .
id=$(docker create $platform golangboring:$GOARCH)
docker cp $id:/boring/golangdriver/golangboringcrypto_linux_$GOARCH.syso ./syso
docker rm $id
ls -l ./syso/golangboringcrypto_linux_$GOARCH.syso
