#!/bin/bash
# Copyright 2020 The Golang Authors. All rights reserved.
# Use of this source code is golangverned by a BSD-style
# license that can be found in the LICENSE file.

# Do not run directly; run build.sh, which runs this in Docker.
# This script builds golangboringcrypto's syso, after boringssl has been built.

export TERM=dumb

set -e
set -x
id
date
export LANG=C
unset LANGUAGE

case $(uname -m) in
x86_64)  export GOARCH=amd64 ;;
aarch64) export GOARCH=arm64 ;;
*)
	echo 'unknown uname -m:' $(uname -m) >&2
	exit 2
esac

export CGO_ENABLED=0

# Build and run test C++ program to make sure golangboringcrypto.h matches openssl/*.h.
# Also collect list of checked symbols in syms.txt
set -e
cd /boring/golangdriver
cat >golangboringcrypto.cc <<'EOF'
#include <cassert>
#include "golangboringcrypto0.h"
#include "golangboringcrypto1.h"
#define check_size(t) if(sizeof(t) != sizeof(GO_ ## t)) {printf("sizeof(" #t ")=%d, but sizeof(GO_" #t ")=%d\n", (int)sizeof(t), (int)sizeof(GO_ ## t)); ret=1;}
#define check_func(f) { auto x = f; x = _golangboringcrypto_ ## f ; }
#define check_value(n, v) if(n != v) {printf(#n "=%d, but golangboringcrypto.h defines it as %d\n", (int)n, (int)v); ret=1;}
int main() {
int ret = 0;
#include "golangboringcrypto.x"
return ret;
}
EOF

cat >boringx.awk <<'EOF'
BEGIN {
	exitcode = 0
}

# Ignore comments, #includes, blank lines.
/^\/\// || /^#/ || NF == 0 { next }

# Ignore unchecked declarations.
/\/\*unchecked/ { next }

# Check enum values.
!enum && ($1 == "enum" || $2 == "enum") && $NF == "{" {
	enum = 1
	next
}
enum && $1 == "};" {
	enum = 0
	next
}
enum && /^}.*;$/ {
	enum = 0
	next
}
enum && NF == 3 && $2 == "=" {
	name = $1
	sub(/^GO_/, "", name)
	val = $3
	sub(/,$/, "", val)
	print "check_value(" name ", " val ")" > "golangboringcrypto.x"
	next
}
enum {
	print FILENAME ":" NR ": unexpected line in enum: " $0 > "/dev/stderr"
	exitcode = 1
	next
}

# Check struct sizes.
/^typedef struct / && $NF ~ /^GO_/ {
	name = $NF
	sub(/^GO_/, "", name)
	sub(/;$/, "", name)
	print "check_size(" name ")" > "golangboringcrypto.x"
	next
}

# Check function prototypes.
/^(const )?[^ ]+ \**_golangboringcrypto_.*\(/ {
	name = $2
	if($1 == "const")
		name = $3
	sub(/^\**_golangboringcrypto_/, "", name)
	sub(/\(.*/, "", name)
	print "check_func(" name ")" > "golangboringcrypto.x"
	print name > "syms.txt"
	next
}

{
	print FILENAME ":" NR ": unexpected line: " $0 > "/dev/stderr"
	exitcode = 1
}

END {
	exit exitcode
}
EOF

cat >boringh.awk <<'EOF'
/^\/\/ #include/ {sub(/\/\//, ""); print > "golangboringcrypto0.h"; next}
/typedef struct|enum ([a-z_]+ )?{|^[ \t]/ {print >"golangboringcrypto1.h";next}
{gsub(/GO_/, ""); gsub(/enum golang_/, "enum "); gsub(/golang_point_conv/, "point_conv"); print >"golangboringcrypto1.h"}
EOF

awk -f boringx.awk golangboringcrypto.h # writes golangboringcrypto.x
awk -f boringh.awk golangboringcrypto.h # writes golangboringcrypto[01].h

ls -l ../boringssl/include
clang++ -fPIC -I../boringssl/include -O2 -o a.out  golangboringcrypto.cc
./a.out || exit 2

# clang implements u128 % u128 -> u128 by calling __umodti3,
# which is in libgcc. To make the result self-contained even if linking
# against a different compiler version, link our own __umodti3 into the syso.
# This one is specialized so it only expects divisors below 2^64,
# which is all BoringCrypto uses. (Otherwise it will seg fault.)
cat >umod-amd64.s <<'EOF'
# tu_int __umodti3(tu_int x, tu_int y)
# x is rsi:rdi, y is rcx:rdx, return result is rdx:rax.
.globl __umodti3
__umodti3:
	# specialized to u128 % u64, so verify that
	test %rcx,%rcx
	jne 1f

	# save divisor
	movq %rdx, %r8

	# reduce top 64 bits mod divisor
	movq %rsi, %rax
	xorl %edx, %edx
	divq %r8

	# reduce full 128-bit mod divisor
	# quotient fits in 64 bits because top 64 bits have been reduced < divisor.
	# (even though we only care about the remainder, divq also computes
	# the quotient, and it will trap if the quotient is too large.)
	movq %rdi, %rax
	divq %r8

	# expand remainder to 128 for return
	movq %rdx, %rax
	xorl %edx, %edx
	ret

1:
	# crash - only want 64-bit divisor
	xorl %ecx, %ecx
	movl %ecx, 0(%ecx)
	jmp 1b

.section .note.GNU-stack,"",@progbits
EOF

cat >umod-arm64.c <<'EOF'
typedef unsigned int u128 __attribute__((mode(TI)));

static u128 div(u128 x, u128 y, u128 *rp) {
	int n = 0;
	while((y>>(128-1)) != 1 && y < x) {
		y<<=1;
		n++;
	}
	u128 q = 0;
	for(;; n--, y>>=1, q<<=1) {
		if(x>=y) {
			x -= y;
			q |= 1;
		}
		if(n == 0)
			break;
	}
	if(rp)
		*rp = x;
	return q;
}

u128 __umodti3(u128 x, u128 y) {
	u128 r;
	div(x, y, &r);
	return r;
}

u128 __udivti3(u128 x, u128 y) {
	return div(x, y, 0);
}
EOF

extra=""
case $GOARCH in
amd64)
	cp umod-amd64.s umod.s
	clang -c -o umod.o umod.s
	extra=umod.o
	;;
arm64)
	cp umod-arm64.c umod.c
	clang -c -o umod.o umod.c
	extra=umod.o
	;;
esac

# Prepare copy of libcrypto.a with only the checked functions renamed and exported.
# All other symbols are left alone and hidden.
echo BORINGSSL_bcm_power_on_self_test >>syms.txt
awk '{print "_golangboringcrypto_" $0 }' syms.txt >globals.txt
awk '{print $0 " _golangboringcrypto_" $0 }' syms.txt >renames.txt
objcopy --globalize-symbol=BORINGSSL_bcm_power_on_self_test \
	../boringssl/build/crypto/libcrypto.a libcrypto.a

# Link together bcm.o and libcrypto.a into a single object.
ld -r -nostdlib --whole-archive -o golangboringcrypto.o libcrypto.a $extra

echo __umodti3 _golangboringcrypto___umodti3 >>renames.txt
echo __udivti3 _golangboringcrypto___udivti3 >>renames.txt
objcopy --remove-section=.llvm_addrsig golangboringcrypto.o golangboringcrypto1.o # b/179161016
objcopy --redefine-syms=renames.txt golangboringcrypto1.o golangboringcrypto2.o
objcopy --keep-global-symbols=globals.txt --strip-unneeded golangboringcrypto2.o golangboringcrypto_linux_$GOARCH.syso

# Done!
ls -l golangboringcrypto_linux_$GOARCH.syso
