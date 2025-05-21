// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package race_test

import (
	"internal/testenv"
	"os"
	"os/exec"
	"testing"
)

func TestNoRaceCgolangSync(t *testing.T) {
	cmd := exec.Command(testenv.GoToolPath(t), "run", "-race", "cgolang_test_main.golang")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("program exited with error: %v\n", err)
	}
}
