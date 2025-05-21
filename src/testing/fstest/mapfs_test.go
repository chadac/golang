// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package fstest

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"
)

func TestMapFS(t *testing.T) {
	m := MapFS{
		"hello":             {Data: []byte("hello, world\n")},
		"fortune/k/ken.txt": {Data: []byte("If a program is too slow, it must have a loop.\n")},
	}
	if err := TestFS(m, "hello", "fortune", "fortune/k", "fortune/k/ken.txt"); err != nil {
		t.Fatal(err)
	}
}

func TestMapFSChmodDot(t *testing.T) {
	m := MapFS{
		"a/b.txt": &MapFile{Mode: 0666},
		".":       &MapFile{Mode: 0777 | fs.ModeDir},
	}
	buf := new(strings.Builder)
	fs.WalkDir(m, ".", func(path string, d fs.DirEntry, err error) error {
		fi, err := d.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "%s: %v\n", path, fi.Mode())
		return nil
	})
	want := `
.: drwxrwxrwx
a: dr-xr-xr-x
a/b.txt: -rw-rw-rw-
`[1:]
	golangt := buf.String()
	if want != golangt {
		t.Errorf("MapFS modes want:\n%s\ngolangt:\n%s\n", want, golangt)
	}
}

func TestMapFSFileInfoName(t *testing.T) {
	m := MapFS{
		"path/to/b.txt": &MapFile{},
	}
	info, _ := m.Stat("path/to/b.txt")
	want := "b.txt"
	golangt := info.Name()
	if want != golangt {
		t.Errorf("MapFS FileInfo.Name want:\n%s\ngolangt:\n%s\n", want, golangt)
	}
}

func TestMapFSSymlink(t *testing.T) {
	const fileContent = "If a program is too slow, it must have a loop.\n"
	m := MapFS{
		"fortune/k/ken.txt": {Data: []byte(fileContent)},
		"dirlink":           {Data: []byte("fortune/k"), Mode: fs.ModeSymlink},
		"linklink":          {Data: []byte("dirlink"), Mode: fs.ModeSymlink},
		"ken.txt":           {Data: []byte("dirlink/ken.txt"), Mode: fs.ModeSymlink},
	}
	if err := TestFS(m, "fortune/k/ken.txt", "dirlink", "ken.txt", "linklink"); err != nil {
		t.Error(err)
	}

	golangtData, err := fs.ReadFile(m, "ken.txt")
	if string(golangtData) != fileContent || err != nil {
		t.Errorf("fs.ReadFile(m, \"ken.txt\") = %q, %v; want %q, <nil>", golangtData, err, fileContent)
	}
	golangtLink, err := fs.ReadLink(m, "dirlink")
	if want := "fortune/k"; golangtLink != want || err != nil {
		t.Errorf("fs.ReadLink(m, \"dirlink\") = %q, %v; want %q, <nil>", golangtLink, err, fileContent)
	}
	golangtInfo, err := fs.Lstat(m, "dirlink")
	if err != nil {
		t.Errorf("fs.Lstat(m, \"dirlink\") = _, %v; want _, <nil>", err)
	} else {
		if golangt, want := golangtInfo.Name(), "dirlink"; golangt != want {
			t.Errorf("fs.Lstat(m, \"dirlink\").Name() = %q; want %q", golangt, want)
		}
		if golangt, want := golangtInfo.Mode(), fs.ModeSymlink; golangt != want {
			t.Errorf("fs.Lstat(m, \"dirlink\").Mode() = %v; want %v", golangt, want)
		}
	}
	golangtInfo, err = fs.Stat(m, "dirlink")
	if err != nil {
		t.Errorf("fs.Stat(m, \"dirlink\") = _, %v; want _, <nil>", err)
	} else {
		if golangt, want := golangtInfo.Name(), "dirlink"; golangt != want {
			t.Errorf("fs.Stat(m, \"dirlink\").Name() = %q; want %q", golangt, want)
		}
		if golangt, want := golangtInfo.Mode(), fs.ModeDir|0555; golangt != want {
			t.Errorf("fs.Stat(m, \"dirlink\").Mode() = %v; want %v", golangt, want)
		}
	}
	golangtInfo, err = fs.Lstat(m, "linklink")
	if err != nil {
		t.Errorf("fs.Lstat(m, \"linklink\") = _, %v; want _, <nil>", err)
	} else {
		if golangt, want := golangtInfo.Name(), "linklink"; golangt != want {
			t.Errorf("fs.Lstat(m, \"linklink\").Name() = %q; want %q", golangt, want)
		}
		if golangt, want := golangtInfo.Mode(), fs.ModeSymlink; golangt != want {
			t.Errorf("fs.Lstat(m, \"linklink\").Mode() = %v; want %v", golangt, want)
		}
	}
	golangtInfo, err = fs.Stat(m, "linklink")
	if err != nil {
		t.Errorf("fs.Stat(m, \"linklink\") = _, %v; want _, <nil>", err)
	} else {
		if golangt, want := golangtInfo.Name(), "linklink"; golangt != want {
			t.Errorf("fs.Stat(m, \"linklink\").Name() = %q; want %q", golangt, want)
		}
		if golangt, want := golangtInfo.Mode(), fs.ModeDir|0555; golangt != want {
			t.Errorf("fs.Stat(m, \"linklink\").Mode() = %v; want %v", golangt, want)
		}
	}
}
