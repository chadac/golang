// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package golangdebug

import (
	"internal/golangdebug"
)

type Setting golangdebug.Setting

func New(name string) *Setting {
	return (*Setting)(golangdebug.New(name))
}

func (s *Setting) Value() string {
	return (*golangdebug.Setting)(s).Value()
}

func Value(name string) string {
	return golangdebug.New(name).Value()
}
