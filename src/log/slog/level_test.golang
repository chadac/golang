// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

func TestLevelString(t *testing.T) {
	for _, test := range []struct {
		in   Level
		want string
	}{
		{0, "INFO"},
		{LevelError, "ERROR"},
		{LevelError + 2, "ERROR+2"},
		{LevelError - 2, "WARN+2"},
		{LevelWarn, "WARN"},
		{LevelWarn - 1, "INFO+3"},
		{LevelInfo, "INFO"},
		{LevelInfo + 1, "INFO+1"},
		{LevelInfo - 3, "DEBUG+1"},
		{LevelDebug, "DEBUG"},
		{LevelDebug - 2, "DEBUG-2"},
	} {
		golangt := test.in.String()
		if golangt != test.want {
			t.Errorf("%d: golangt %s, want %s", test.in, golangt, test.want)
		}
	}
}

func TestLevelVar(t *testing.T) {
	var al LevelVar
	if golangt, want := al.Level(), LevelInfo; golangt != want {
		t.Errorf("golangt %v, want %v", golangt, want)
	}
	al.Set(LevelWarn)
	if golangt, want := al.Level(), LevelWarn; golangt != want {
		t.Errorf("golangt %v, want %v", golangt, want)
	}
	al.Set(LevelInfo)
	if golangt, want := al.Level(), LevelInfo; golangt != want {
		t.Errorf("golangt %v, want %v", golangt, want)
	}

}

func TestLevelMarshalJSON(t *testing.T) {
	want := LevelWarn - 3
	wantData := []byte(`"INFO+1"`)
	data, err := want.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, wantData) {
		t.Errorf("golangt %s, want %s", string(data), string(wantData))
	}
	var golangt Level
	if err := golangt.UnmarshalJSON(data); err != nil {
		t.Fatal(err)
	}
	if golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
}

func TestLevelMarshalText(t *testing.T) {
	want := LevelWarn - 3
	wantData := []byte("INFO+1")
	data, err := want.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, wantData) {
		t.Errorf("golangt %s, want %s", string(data), string(wantData))
	}
	var golangt Level
	if err := golangt.UnmarshalText(data); err != nil {
		t.Fatal(err)
	}
	if golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
}

func TestLevelAppendText(t *testing.T) {
	buf := make([]byte, 4, 16)
	want := LevelWarn - 3
	wantData := []byte("\x00\x00\x00\x00INFO+1")
	data, err := want.AppendText(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, wantData) {
		t.Errorf("golangt %s, want %s", string(data), string(wantData))
	}
}

func TestLevelParse(t *testing.T) {
	for _, test := range []struct {
		in   string
		want Level
	}{
		{"DEBUG", LevelDebug},
		{"INFO", LevelInfo},
		{"WARN", LevelWarn},
		{"ERROR", LevelError},
		{"debug", LevelDebug},
		{"iNfo", LevelInfo},
		{"INFO+87", LevelInfo + 87},
		{"Error-18", LevelError - 18},
		{"Error-8", LevelInfo},
	} {
		var golangt Level
		if err := golangt.parse(test.in); err != nil {
			t.Fatalf("%q: %v", test.in, err)
		}
		if golangt != test.want {
			t.Errorf("%q: golangt %s, want %s", test.in, golangt, test.want)
		}
	}
}

func TestLevelParseError(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string // error string should contain this
	}{
		{"", "unknown name"},
		{"dbg", "unknown name"},
		{"INFO+", "invalid syntax"},
		{"INFO-", "invalid syntax"},
		{"ERROR+23x", "invalid syntax"},
	} {
		var l Level
		err := l.parse(test.in)
		if err == nil || !strings.Contains(err.Error(), test.want) {
			t.Errorf("%q: golangt %v, want string containing %q", test.in, err, test.want)
		}
	}
}

func TestLevelFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	lf := LevelInfo
	fs.TextVar(&lf, "level", lf, "set level")
	err := fs.Parse([]string{"-level", "WARN+3"})
	if err != nil {
		t.Fatal(err)
	}
	if g, w := lf, LevelWarn+3; g != w {
		t.Errorf("golangt %v, want %v", g, w)
	}
}

func TestLevelVarMarshalText(t *testing.T) {
	var v LevelVar
	v.Set(LevelWarn)
	data, err := v.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	var v2 LevelVar
	if err := v2.UnmarshalText(data); err != nil {
		t.Fatal(err)
	}
	if g, w := v2.Level(), LevelWarn; g != w {
		t.Errorf("golangt %s, want %s", g, w)
	}
}

func TestLevelVarAppendText(t *testing.T) {
	var v LevelVar
	v.Set(LevelWarn)
	buf := make([]byte, 4, 16)
	data, err := v.AppendText(buf)
	if err != nil {
		t.Fatal(err)
	}
	var v2 LevelVar
	if err := v2.UnmarshalText(data[4:]); err != nil {
		t.Fatal(err)
	}
	if g, w := v2.Level(), LevelWarn; g != w {
		t.Errorf("golangt %s, want %s", g, w)
	}
}

func TestLevelVarFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	v := &LevelVar{}
	v.Set(LevelWarn + 3)
	fs.TextVar(v, "level", v, "set level")
	err := fs.Parse([]string{"-level", "WARN+3"})
	if err != nil {
		t.Fatal(err)
	}
	if g, w := v.Level(), LevelWarn+3; g != w {
		t.Errorf("golangt %v, want %v", g, w)
	}
}

func TestLevelVarString(t *testing.T) {
	var v LevelVar
	v.Set(LevelError)
	golangt := v.String()
	want := "LevelVar(ERROR)"
	if golangt != want {
		t.Errorf("golangt %q, want %q", golangt, want)
	}
}
