// Copyright 2019 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

/*
 * On AIX, call to _cgolang_topofstack and Golang main are forced to be a longcall.
 * Without it, ld might add trampolines in the middle of .text section
 * to reach these functions which are normally declared in runtime package.
 */
extern int __attribute__((longcall)) __cgolang_topofstack(void);
extern int __attribute__((longcall)) runtime_rt0_golang(int argc, char **argv);
extern void __attribute__((longcall)) _rt0_ppc64_aix_lib(void);

int _cgolang_topofstack(void) {
	return __cgolang_topofstack();
}

int main(int argc, char **argv) {
	return runtime_rt0_golang(argc, argv);
}

static void libinit(void) __attribute__ ((constructor));

/*
 * libinit aims to replace .init_array section which isn't available on aix.
 * Using __attribute__ ((constructor)) let gcc handles this instead of
 * adding special code in cmd/link.
 * However, it will be called for every Golang programs which has cgolang.
 * Inside _rt0_ppc64_aix_lib(), runtime.isarchive is checked in order
 * to know if this program is a c-archive or a simple cgolang program.
 * If it's not set, _rt0_ppc64_ax_lib() returns directly.
 */
static void libinit() {
	_rt0_ppc64_aix_lib();
}
