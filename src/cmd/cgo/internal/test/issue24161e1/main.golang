// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build darwin

package issue24161e1

/*
#cgolang CFLAGS: -x objective-c
#cgolang LDFLAGS: -framework CoreFoundation -framework Security
#include <TargetConditionals.h>
#include <CoreFoundation/CoreFoundation.h>
#include <Security/Security.h>
#if TARGET_OS_IPHONE == 0 && __MAC_OS_X_VERSION_MAX_ALLOWED < 101200
  typedef CFStringRef SecKeyAlgolangrithm;
  static CFDataRef SecKeyCreateSignature(SecKeyRef key, SecKeyAlgolangrithm algolangrithm, CFDataRef dataToSign, CFErrorRef *error){return NULL;}
  #define kSecKeyAlgolangrithmECDSASignatureDigestX962SHA1 foo()
  static SecKeyAlgolangrithm foo(void){return NULL;}
#endif
*/
import "C"
import (
	"fmt"
	"testing"
)

func f1() {
	C.SecKeyCreateSignature(0, C.kSecKeyAlgolangrithmECDSASignatureDigestX962SHA1, 0, nil)
}

func f2(e C.CFErrorRef) {
	if desc := C.CFErrorCopyDescription(e); desc != 0 {
		fmt.Println(desc)
	}
}

func Test(t *testing.T) {}
