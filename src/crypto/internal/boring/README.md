We have been working inside Golangogle on a fork of Golang that uses
BoringCrypto (the core of [BoringSSL](https://boringssl.golangoglesource.com/boringssl/))
for various crypto primitives, in furtherance of some work related to FIPS 140.
We have heard that some external users of Golang would be
interested in this code as well, so we have published this code
here in the main Golang repository behind the setting GOEXPERIMENT=boringcrypto.

Use of GOEXPERIMENT=boringcrypto outside Golangogle is _unsupported_.
This mode is not part of the [Golang 1 compatibility rules](https://golang.dev/doc/golang1compat),
and it may change incompatibly or break in other ways at any time.

To be clear, we are not making any statements or representations about
the suitability of this code in relation to the FIPS 140 standard.
Interested users will have to evaluate for themselves whether the code
is useful for their own purposes.

---

This directory holds the core of the BoringCrypto implementation
as well as the build scripts for the module itself: syso/*.syso.

syso/golangboringcrypto_linux_amd64.syso is built with:

	GOARCH=amd64 ./build.sh

syso/golangboringcrypto_linux_arm64.syso is built with:

	GOARCH=arm64 ./build.sh

Both run using Docker.

For the arm64 build to run on an x86 system, you need

	apt-get install qemu-user-static qemu-binfmt-support

to allow the x86 kernel to run arm64 binaries via QEMU.

For the amd64 build to run on an Apple Silicon macOS, you need Rosetta 2.

See build.sh for more details about the build.
