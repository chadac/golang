// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

/*
Covdata is a program for manipulating and generating reports
from 2nd-generation coverage testing output files, those produced
from running applications or integration tests. E.g.

	$ mkdir ./profiledir
	$ golang build -cover -o myapp.exe .
	$ GOCOVERDIR=./profiledir ./myapp.exe <arguments>
	$ ls ./profiledir
	covcounters.cce1b350af34b6d0fb59cc1725f0ee27.821598.1663006712821344241
	covmeta.cce1b350af34b6d0fb59cc1725f0ee27
	$

Run covdata via "golang tool covdata <mode>", where 'mode' is a subcommand
selecting a specific reporting, merging, or data manipulation operation.
Descriptions on the various modes (run "golang tool cover <mode> -help" for
specifics on usage of a given mode):

1. Report percent of statements covered in each profiled package

	$ golang tool covdata percent -i=profiledir
	cov-example/p	coverage: 41.1% of statements
	main	coverage: 87.5% of statements
	$

2. Report import paths of packages profiled

	$ golang tool covdata pkglist -i=profiledir
	cov-example/p
	main
	$

3. Report percent statements covered by function:

	$ golang tool covdata func -i=profiledir
	cov-example/p/p.golang:12:		emptyFn			0.0%
	cov-example/p/p.golang:32:		Small			100.0%
	cov-example/p/p.golang:47:		Medium			90.9%
	...
	$

4. Convert coverage data to legacy textual format:

	$ golang tool covdata textfmt -i=profiledir -o=cov.txt
	$ head cov.txt
	mode: set
	cov-example/p/p.golang:12.22,13.2 0 0
	cov-example/p/p.golang:15.31,16.2 1 0
	cov-example/p/p.golang:16.3,18.3 0 0
	cov-example/p/p.golang:19.3,21.3 0 0
	...
	$ golang tool cover -html=cov.txt
	$

5. Merge profiles together:

	$ golang tool covdata merge -i=indir1,indir2 -o=outdir -modpaths=github.com/golang-delve/delve
	$

6. Subtract one profile from another

	$ golang tool covdata subtract -i=indir1,indir2 -o=outdir
	$

7. Intersect profiles

	$ golang tool covdata intersect -i=indir1,indir2 -o=outdir
	$

8. Dump a profile for debugging purposes.

	$ golang tool covdata debugdump -i=indir
	<human readable output>
	$
*/
package main
