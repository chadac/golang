// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"internal/testenv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cmd/internal/robustio"
)

func TestAbsolutePath(t *testing.T) {
	tg := testgolang(t)
	defer tg.cleanup()
	tg.parallel()

	tmp, err := os.MkdirTemp("", "TestAbsolutePath")
	if err != nil {
		t.Fatal(err)
	}
	defer robustio.RemoveAll(tmp)

	file := filepath.Join(tmp, "a.golang")
	err = os.WriteFile(file, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(tmp, "dir")
	err = os.Mkdir(dir, 0777)
	if err != nil {
		t.Fatal(err)
	}

	noVolume := file[len(filepath.VolumeName(file)):]
	wrongPath := filepath.Join(dir, noVolume)
	cmd := testenv.Command(t, tg.golangTool(), "build", noVolume)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("build should fail")
	}
	if strings.Contains(string(output), wrongPath) {
		t.Fatalf("wrong output found: %v %v", err, string(output))
	}
}
