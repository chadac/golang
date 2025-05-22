// Copyright 2019 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "_cgolang_export.h"

void callIssue31891() {
    Issue31891A a;
    useIssue31891A(&a);

    Issue31891B b;
    useIssue31891B(&b);
}
