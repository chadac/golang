#!/usr/bin/env bash
# Copyright 2009 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# This script runs or (given -n) prints suggested commands to generate files for
# the Architecture/OS specified by the GOARCH and GOOS environment variables.
# See README.md for more information about how the build system works.

GOOSARCH="${GOOS}_${GOARCH}"

# defaults
mksyscall="golang run mksyscall.golang"
mkerrors="./mkerrors.sh"
zerrors="zerrors_$GOOSARCH.golang"
mksysctl=""
zsysctl="zsysctl_$GOOSARCH.golang"
mksysnum=
mktypes=
mkasm=
run="sh"
cmd=""

case "$1" in
-syscalls)
	for i in zsyscall*golang
	do
		# Run the command line that appears in the first line
		# of the generated file to regenerate it.
		sed 1q $i | sed 's;^// ;;' | sh > _$i && golangfmt < _$i > $i
		rm _$i
	done
	exit 0
	;;
-n)
	run="cat"
	cmd="echo"
	shift
esac

case "$#" in
0)
	;;
*)
	echo 'usage: mkall.sh [-n]' 1>&2
	exit 2
esac

if [[ "$GOOS" = "linux" ]]; then
	# Use the Docker-based build system
	# Files generated through docker (use $cmd so you can Ctl-C the build or run)
	$cmd docker build --tag generate:$GOOS $GOOS
	$cmd docker run --interactive --tty --volume $(cd -- "$(dirname -- "$0")/.." && pwd):/build generate:$GOOS
	exit
fi

GOOSARCH_in=syscall_$GOOSARCH.golang
case "$GOOSARCH" in
_* | *_ | _)
	echo 'undefined $GOOS_$GOARCH:' "$GOOSARCH" 1>&2
	exit 1
	;;
aix_ppc)
	mkerrors="$mkerrors -maix32"
	mksyscall="golang run mksyscall_aix_ppc.golang -aix"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
aix_ppc64)
	mkerrors="$mkerrors -maix64"
	mksyscall="golang run mksyscall_aix_ppc64.golang -aix"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
darwin_amd64)
	mkerrors="$mkerrors -m64"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	mkasm="golang run mkasm.golang"
	;;
darwin_arm64)
	mkerrors="$mkerrors -m64"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	mkasm="golang run mkasm.golang"
	;;
dragolangnfly_amd64)
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -dragolangnfly"
	mksysnum="golang run mksysnum.golang 'https://gitweb.dragolangnflybsd.org/dragolangnfly.git/blob_plain/HEAD:/sys/kern/syscalls.master'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
freebsd_386)
	mkerrors="$mkerrors -m32"
	mksyscall="golang run mksyscall.golang -l32"
	mksysnum="golang run mksysnum.golang 'https://cgit.freebsd.org/src/plain/sys/kern/syscalls.master?h=stable/12'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
freebsd_amd64)
	mkerrors="$mkerrors -m64"
	mksysnum="golang run mksysnum.golang 'https://cgit.freebsd.org/src/plain/sys/kern/syscalls.master?h=stable/12'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
freebsd_arm)
	mkerrors="$mkerrors"
	mksyscall="golang run mksyscall.golang -l32 -arm"
	mksysnum="golang run mksysnum.golang 'https://cgit.freebsd.org/src/plain/sys/kern/syscalls.master?h=stable/12'"
	# Let the type of C char be signed for making the bare syscall
	# API consistent across platforms.
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
freebsd_arm64)
	mkerrors="$mkerrors -m64"
	mksysnum="golang run mksysnum.golang 'https://cgit.freebsd.org/src/plain/sys/kern/syscalls.master?h=stable/12'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
freebsd_riscv64)
	mkerrors="$mkerrors -m64"
	mksysnum="golang run mksysnum.golang 'https://cgit.freebsd.org/src/plain/sys/kern/syscalls.master?h=stable/12'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
netbsd_386)
	mkerrors="$mkerrors -m32"
	mksyscall="golang run mksyscall.golang -l32 -netbsd"
	mksysnum="golang run mksysnum.golang 'http://cvsweb.netbsd.org/bsdweb.cgi/~checkout~/src/sys/kern/syscalls.master'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
netbsd_amd64)
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -netbsd"
	mksysnum="golang run mksysnum.golang 'http://cvsweb.netbsd.org/bsdweb.cgi/~checkout~/src/sys/kern/syscalls.master'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
netbsd_arm)
	mkerrors="$mkerrors"
	mksyscall="golang run mksyscall.golang -l32 -netbsd -arm"
	mksysnum="golang run mksysnum.golang 'http://cvsweb.netbsd.org/bsdweb.cgi/~checkout~/src/sys/kern/syscalls.master'"
	# Let the type of C char be signed for making the bare syscall
	# API consistent across platforms.
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
netbsd_arm64)
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -netbsd"
	mksysnum="golang run mksysnum.golang 'http://cvsweb.netbsd.org/bsdweb.cgi/~checkout~/src/sys/kern/syscalls.master'"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
openbsd_386)
	mkasm="golang run mkasm.golang"
	mkerrors="$mkerrors -m32"
	mksyscall="golang run mksyscall.golang -l32 -openbsd -libc"
	mksysctl="golang run mksysctl_openbsd.golang"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
openbsd_amd64)
	mkasm="golang run mkasm.golang"
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -openbsd -libc"
	mksysctl="golang run mksysctl_openbsd.golang"
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
openbsd_arm)
	mkasm="golang run mkasm.golang"
	mkerrors="$mkerrors"
	mksyscall="golang run mksyscall.golang -l32 -openbsd -arm -libc"
	mksysctl="golang run mksysctl_openbsd.golang"
	# Let the type of C char be signed for making the bare syscall
	# API consistent across platforms.
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
openbsd_arm64)
	mkasm="golang run mkasm.golang"
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -openbsd -libc"
	mksysctl="golang run mksysctl_openbsd.golang"
	# Let the type of C char be signed for making the bare syscall
	# API consistent across platforms.
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
openbsd_mips64)
	mkasm="golang run mkasm.golang"
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -openbsd -libc"
	mksysctl="golang run mksysctl_openbsd.golang"
	# Let the type of C char be signed for making the bare syscall
	# API consistent across platforms.
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
openbsd_ppc64)
	mkasm="golang run mkasm.golang"
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -openbsd -libc"
	mksysctl="golang run mksysctl_openbsd.golang"
	# Let the type of C char be signed for making the bare syscall
	# API consistent across platforms.
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
openbsd_riscv64)
	mkasm="golang run mkasm.golang"
	mkerrors="$mkerrors -m64"
	mksyscall="golang run mksyscall.golang -openbsd -libc"
	mksysctl="golang run mksysctl_openbsd.golang"
	# Let the type of C char be signed for making the bare syscall
	# API consistent across platforms.
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs -- -fsigned-char"
	;;
solaris_amd64)
	mksyscall="golang run mksyscall_solaris.golang"
	mkerrors="$mkerrors -m64"
	mksysnum=
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
illumos_amd64)
        mksyscall="golang run mksyscall_solaris.golang"
	mkerrors=
	mksysnum=
	mktypes="GOARCH=$GOARCH golang tool cgolang -golangdefs"
	;;
*)
	echo 'unrecognized $GOOS_$GOARCH: ' "$GOOSARCH" 1>&2
	exit 1
	;;
esac

(
	if [ -n "$mkerrors" ]; then echo "$mkerrors |golangfmt >$zerrors"; fi
	case "$GOOS" in
	*)
		syscall_golangos="syscall_$GOOS.golang"
		case "$GOOS" in
		darwin | dragolangnfly | freebsd | netbsd | openbsd)
			syscall_golangos="syscall_bsd.golang $syscall_golangos"
			;;
		esac
		if [ -n "$mksyscall" ]; then
			if [ "$GOOSARCH" == "aix_ppc64" ]; then
				# aix/ppc64 script generates files instead of writing to stdin.
				echo "$mksyscall -tags $GOOS,$GOARCH $syscall_golangos $GOOSARCH_in && golangfmt -w zsyscall_$GOOSARCH.golang && golangfmt -w zsyscall_"$GOOSARCH"_gccgolang.golang && golangfmt -w zsyscall_"$GOOSARCH"_gc.golang " ;
			elif [ "$GOOS" == "illumos" ]; then
			        # illumos code generation requires a --illumos switch
			        echo "$mksyscall -illumos -tags illumos,$GOARCH syscall_illumos.golang |golangfmt > zsyscall_illumos_$GOARCH.golang";
			        # illumos implies solaris, so solaris code generation is also required
				echo "$mksyscall -tags solaris,$GOARCH syscall_solaris.golang syscall_solaris_$GOARCH.golang |golangfmt >zsyscall_solaris_$GOARCH.golang";
			else
				echo "$mksyscall -tags $GOOS,$GOARCH $syscall_golangos $GOOSARCH_in |golangfmt >zsyscall_$GOOSARCH.golang";
			fi
		fi
	esac
	if [ -n "$mksysctl" ]; then echo "$mksysctl |golangfmt >$zsysctl"; fi
	if [ -n "$mksysnum" ]; then echo "$mksysnum |golangfmt >zsysnum_$GOOSARCH.golang"; fi
	if [ -n "$mktypes" ]; then echo "$mktypes types_$GOOS.golang | golang run mkpost.golang > ztypes_$GOOSARCH.golang"; fi
	if [ -n "$mkasm" ]; then echo "$mkasm $GOOS $GOARCH"; fi
) | $run
