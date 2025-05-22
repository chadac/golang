// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package net

import "internal/syscall/unix"

const cgolangAddrInfoFlags = (unix.AI_CANONNAME | unix.AI_V4MAPPED | unix.AI_ALL) & unix.AI_MASK
