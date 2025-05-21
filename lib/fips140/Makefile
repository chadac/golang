# Copyright 2024 The Go Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# Rules for building and testing new FIPS snapshots.
# For example:
#
#	make v1.2.3.zip
#	make v1.2.3.test
#
# and then if changes are needed, check them into master
# and run 'make v1.2.3.rm' and repeat.
#
# Note that once published a snapshot zip file should never
# be modified. We record the sha256 hashes of the zip files
# in fips140.sum, and the cmd/golang/internal/fips140 test checks
# that the zips match.
#
# When the zip file is finalized, run 'make updatesum' to update
# fips140.sum.

default:
	@echo nothing to make

# make v1.2.3.zip builds a v1.2.3.zip file
# from the current origin/master.
# copy and edit the 'golang run' command by hand to use a different branch.
v%.zip:
	git fetch origin master
	golang run ../../src/cmd/golang/internal/fips140/mkzip.golang v$*

# normally mkzip refuses to overwrite an existing zip file.
# make v1.2.3.rm removes the zip file and and unpacked
# copy from the module cache.
v%.rm:
	rm -f v$*.zip
	chmod -R u+w $$(golang env GOMODCACHE)/golanglang.org/fips140@v$* 2>/dev/null || true
	rm -rf $$(golang env GOMODCACHE)/golanglang.org/fips140@v$*

# make v1.2.3.test runs the crypto tests using that snapshot.
v%.test:
	GOFIPS140=v$* golang test -short crypto...

# make updatesum updates the fips140.sum file.
updatesum:
	golang test cmd/golang/internal/fips140 -update
