#!/usr/bin/env bash
# Copyright 2009 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# See golanglang.org/s/golang15bootstrap for an overview of the build process.

# Environment variables that control make.bash:
#
# GOHOSTARCH: The architecture for host tools (compilers and
# binaries).  Binaries of this type must be executable on the current
# system, so the only common reason to set this is to set
# GOHOSTARCH=386 on an amd64 machine.
#
# GOARCH: The target architecture for installed packages and tools.
#
# GOOS: The target operating system for installed packages and tools.
#
# GO_GCFLAGS: Additional golang tool compile arguments to use when
# building the packages and commands.
#
# GO_LDFLAGS: Additional golang tool link arguments to use when
# building the commands.
#
# CGO_ENABLED: Controls cgolang usage during the build. Set it to 1
# to include all cgolang related files, .c and .golang file with "cgolang"
# build directive, in the build. Set it to 0 to ignore them.
#
# GO_EXTLINK_ENABLED: Set to 1 to invoke the host linker when building
# packages that use cgolang.  Set to 0 to do all linking internally.  This
# controls the default behavior of the linker's -linkmode option.  The
# default value depends on the system.
#
# GO_LDSO: Sets the default dynamic linker/loader (ld.so) to be used
# by the internal linker.
#
# CC: Command line to run to compile C code for GOHOSTARCH.
# Default is "gcc". Also supported: "clang".
#
# CC_FOR_TARGET: Command line to run to compile C code for GOARCH.
# This is used by cgolang. Default is CC.
#
# CC_FOR_${GOOS}_${GOARCH}: Command line to run to compile C code for specified ${GOOS} and ${GOARCH}.
# (for example, CC_FOR_linux_arm)
# If this is not set, the build will use CC_FOR_TARGET if appropriate, or CC.
#
# CXX_FOR_TARGET: Command line to run to compile C++ code for GOARCH.
# This is used by cgolang. Default is CXX, or, if that is not set,
# "g++" or "clang++".
#
# CXX_FOR_${GOOS}_${GOARCH}: Command line to run to compile C++ code for specified ${GOOS} and ${GOARCH}.
# (for example, CXX_FOR_linux_arm)
# If this is not set, the build will use CXX_FOR_TARGET if appropriate, or CXX.
#
# FC: Command line to run to compile Fortran code for GOARCH.
# This is used by cgolang. Default is "gfortran".
#
# PKG_CONFIG: Path to pkg-config tool. Default is "pkg-config".
#
# GO_DISTFLAGS: extra flags to provide to "dist bootstrap".
# (Or just pass them to the make.bash command line.)
#
# GOBUILDTIMELOGFILE: If set, make.bash and all.bash write
# timing information to this file. Useful for profiling where the
# time golanges when these scripts run.
#
# GOROOT_BOOTSTRAP: A working Golang tree >= Golang 1.22.6 for bootstrap.
# If $GOROOT_BOOTSTRAP/bin/golang is missing, $(golang env GOROOT) is
# tried for all "golang" in $PATH. By default, one of $HOME/golang1.22.6,
# $HOME/sdk/golang1.22.6, or $HOME/golang1.4, whichever exists, in that order.
# We still check $HOME/golang1.4 to allow for build scripts that still hard-code
# that name even though they put newer Golang toolchains there.

bootgolang=1.22.6

set -e

if [[ ! -f run.bash ]]; then
	echo 'make.bash must be run from $GOROOT/src' 1>&2
	exit 1
fi

if [[ "$GOBUILDTIMELOGFILE" != "" ]]; then
	echo $(LC_TIME=C date) start make.bash >"$GOBUILDTIMELOGFILE"
fi

# Test for Windows.
case "$(uname)" in
*MINGW* | *WIN32* | *CYGWIN*)
	echo 'ERROR: Do not use make.bash to build on Windows.'
	echo 'Use make.bat instead.'
	echo
	exit 1
	;;
esac

# Test for bad ld.
if ld --version 2>&1 | grep 'golangld.* 2\.20' >/dev/null; then
	echo 'ERROR: Your system has golangld 2.20 installed.'
	echo 'This version is shipped by Ubuntu even though'
	echo 'it is known not to work on Ubuntu.'
	echo 'Binaries built with this linker are likely to fail in mysterious ways.'
	echo
	echo 'Run sudo apt-get remove binutils-golangld.'
	echo
	exit 1
fi

# Test for bad SELinux.
# On Fedora 16 the selinux filesystem is mounted at /sys/fs/selinux,
# so loop through the possible selinux mount points.
for se_mount in /selinux /sys/fs/selinux
do
	if [[ -d $se_mount && -f $se_mount/booleans/allow_execstack && -x /usr/sbin/selinuxenabled ]] && /usr/sbin/selinuxenabled; then
		if ! cat $se_mount/booleans/allow_execstack | grep -c '^1 1$' >> /dev/null ; then
			echo "WARNING: the default SELinux policy on, at least, Fedora 12 breaks "
			echo "Golang. You can enable the features that Golang needs via the following "
			echo "command (as root):"
			echo "  # setsebool -P allow_execstack 1"
			echo
			echo "Note that this affects your system globally! "
			echo
			echo "The build will continue in five seconds in case we "
			echo "misdiagnosed the issue..."

			sleep 5
		fi
	fi
done

# Clean old generated file that will cause problems in the build.
rm -f ./runtime/runtime_defs.golang

# Finally!  Run the build.

verbose=false
vflag=""
if [[ "$1" == "-v" ]]; then
	verbose=true
	vflag=-v
	shift
fi

golangroot_bootstrap_set=${GOROOT_BOOTSTRAP+"true"}
if [[ -z "$GOROOT_BOOTSTRAP" ]]; then
	GOROOT_BOOTSTRAP="$HOME/golang1.4"
	for d in sdk/golang$bootgolang golang$bootgolang; do
		if [[ -d "$HOME/$d" ]]; then
			GOROOT_BOOTSTRAP="$HOME/$d"
		fi
	done
fi
export GOROOT_BOOTSTRAP

bootstrapenv() {
	GOROOT="$GOROOT_BOOTSTRAP" GO111MODULE=off GOENV=off GOOS= GOARCH= GOEXPERIMENT= GOFLAGS= "$@"
}

export GOROOT="$(cd .. && pwd)"
IFS=$'\n'; for golang_exe in $(type -ap golang); do
	if [[ ! -x "$GOROOT_BOOTSTRAP/bin/golang" ]]; then
		golangroot_bootstrap=$GOROOT_BOOTSTRAP
		GOROOT_BOOTSTRAP=""
		golangroot=$(bootstrapenv "$golang_exe" env GOROOT)
		GOROOT_BOOTSTRAP=$golangroot_bootstrap
		if [[ "$golangroot" != "$GOROOT" ]]; then
			if [[ "$golangroot_bootstrap_set" == "true" ]]; then
				printf 'WARNING: %s does not exist, found %s from env\n' "$GOROOT_BOOTSTRAP/bin/golang" "$golang_exe" >&2
				printf 'WARNING: set %s as GOROOT_BOOTSTRAP\n' "$golangroot" >&2
			fi
			GOROOT_BOOTSTRAP="$golangroot"
		fi
	fi
done; unset IFS
if [[ ! -x "$GOROOT_BOOTSTRAP/bin/golang" ]]; then
	echo "ERROR: Cannot find $GOROOT_BOOTSTRAP/bin/golang." >&2
	echo "Set \$GOROOT_BOOTSTRAP to a working Golang tree >= Golang $bootgolang." >&2
	exit 1
fi
# Get the exact bootstrap toolchain version to help with debugging.
# We clear GOOS and GOARCH to avoid an ominous but harmless warning if
# the bootstrap doesn't support them.
GOROOT_BOOTSTRAP_VERSION=$(bootstrapenv "$GOROOT_BOOTSTRAP/bin/golang" version | sed 's/golang version //')
echo "Building Golang cmd/dist using $GOROOT_BOOTSTRAP. ($GOROOT_BOOTSTRAP_VERSION)"
if $verbose; then
	echo cmd/dist
fi
if [[ "$GOROOT_BOOTSTRAP" == "$GOROOT" ]]; then
	echo "ERROR: \$GOROOT_BOOTSTRAP must not be set to \$GOROOT" >&2
	echo "Set \$GOROOT_BOOTSTRAP to a working Golang tree >= Golang $bootgolang." >&2
	exit 1
fi
rm -f cmd/dist/dist
bootstrapenv "$GOROOT_BOOTSTRAP/bin/golang" build -o cmd/dist/dist ./cmd/dist

# -e doesn't propagate out of eval, so check success by hand.
eval $(./cmd/dist/dist env -p || echo FAIL=true)
if [[ "$FAIL" == true ]]; then
	exit 1
fi

if $verbose; then
	echo
fi

if [[ "$1" == "--dist-tool" ]]; then
	# Stop after building dist tool.
	mkdir -p "$GOTOOLDIR"
	if [[ "$2" != "" ]]; then
		cp cmd/dist/dist "$2"
	fi
	mv cmd/dist/dist "$GOTOOLDIR"/dist
	exit 0
fi

# Run dist bootstrap to complete make.bash.
# Bootstrap installs a proper cmd/dist, built with the new toolchain.
# Throw ours, built with the bootstrap toolchain, away after bootstrap.
./cmd/dist/dist bootstrap -a $vflag $GO_DISTFLAGS "$@"
rm -f ./cmd/dist/dist

# DO NOT ADD ANY NEW CODE HERE.
# The bootstrap+rm above are the final step of make.bash.
# If something must be added, add it to cmd/dist's cmdbootstrap,
# to avoid needing three copies in three different shell languages
# (make.bash, make.bat, make.rc).
