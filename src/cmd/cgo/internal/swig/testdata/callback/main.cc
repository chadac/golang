// Copyright 2013 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This .cc file will be automatically compiled by the golang tool and
// included in the package.

#include <string>
#include "main.h"

std::string Caller::call() {
	if (callback_ != 0)
		return callback_->run();
	return "";
}
