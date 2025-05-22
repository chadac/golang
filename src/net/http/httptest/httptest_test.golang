// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package httptest

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	golangt := NewRequest("GET", "/", nil)
	want := &http.Request{
		Method:     "GET",
		Host:       "example.com",
		URL:        &url.URL{Path: "/"},
		Header:     http.Header{},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		RemoteAddr: "192.0.2.1:1234",
		RequestURI: "/",
	}
	golangt.Body = nil // before DeepEqual
	want = want.WithContext(context.Background())
	if !reflect.DeepEqual(golangt, want) {
		t.Errorf("Request mismatch:\n golangt: %#v\nwant: %#v", golangt, want)
	}
}

func TestNewRequestWithContext(t *testing.T) {
	for _, tt := range [...]struct {
		name string

		method, uri string
		body        io.Reader

		want     *http.Request
		wantBody string
	}{
		{
			name:   "Empty method means GET",
			method: "",
			uri:    "/",
			body:   nil,
			want: &http.Request{
				Method:     "GET",
				Host:       "example.com",
				URL:        &url.URL{Path: "/"},
				Header:     http.Header{},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				RemoteAddr: "192.0.2.1:1234",
				RequestURI: "/",
			},
			wantBody: "",
		},

		{
			name:   "GET with full URL",
			method: "GET",
			uri:    "http://foo.com/path/%2f/bar/",
			body:   nil,
			want: &http.Request{
				Method: "GET",
				Host:   "foo.com",
				URL: &url.URL{
					Scheme:  "http",
					Path:    "/path///bar/",
					RawPath: "/path/%2f/bar/",
					Host:    "foo.com",
				},
				Header:     http.Header{},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				RemoteAddr: "192.0.2.1:1234",
				RequestURI: "http://foo.com/path/%2f/bar/",
			},
			wantBody: "",
		},

		{
			name:   "GET with full https URL",
			method: "GET",
			uri:    "https://foo.com/path/",
			body:   nil,
			want: &http.Request{
				Method: "GET",
				Host:   "foo.com",
				URL: &url.URL{
					Scheme: "https",
					Path:   "/path/",
					Host:   "foo.com",
				},
				Header:     http.Header{},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				RemoteAddr: "192.0.2.1:1234",
				RequestURI: "https://foo.com/path/",
				TLS: &tls.ConnectionState{
					Version:           tls.VersionTLS12,
					HandshakeComplete: true,
					ServerName:        "foo.com",
				},
			},
			wantBody: "",
		},

		{
			name:   "Post with known length",
			method: "POST",
			uri:    "/",
			body:   strings.NewReader("foo"),
			want: &http.Request{
				Method:        "POST",
				Host:          "example.com",
				URL:           &url.URL{Path: "/"},
				Header:        http.Header{},
				Proto:         "HTTP/1.1",
				ContentLength: 3,
				ProtoMajor:    1,
				ProtoMinor:    1,
				RemoteAddr:    "192.0.2.1:1234",
				RequestURI:    "/",
			},
			wantBody: "foo",
		},

		{
			name:   "Post with unknown length",
			method: "POST",
			uri:    "/",
			body:   struct{ io.Reader }{strings.NewReader("foo")},
			want: &http.Request{
				Method:        "POST",
				Host:          "example.com",
				URL:           &url.URL{Path: "/"},
				Header:        http.Header{},
				Proto:         "HTTP/1.1",
				ContentLength: -1,
				ProtoMajor:    1,
				ProtoMinor:    1,
				RemoteAddr:    "192.0.2.1:1234",
				RequestURI:    "/",
			},
			wantBody: "foo",
		},

		{
			name:   "Post with NoBody",
			method: "POST",
			uri:    "/",
			body:   http.NoBody,
			want: &http.Request{
				Method:     "POST",
				Host:       "example.com",
				URL:        &url.URL{Path: "/"},
				Header:     http.Header{},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				RemoteAddr: "192.0.2.1:1234",
				RequestURI: "/",
			},
		},

		{
			name:   "OPTIONS *",
			method: "OPTIONS",
			uri:    "*",
			want: &http.Request{
				Method:     "OPTIONS",
				Host:       "example.com",
				URL:        &url.URL{Path: "*"},
				Header:     http.Header{},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				RemoteAddr: "192.0.2.1:1234",
				RequestURI: "*",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			golangt := NewRequestWithContext(context.Background(), tt.method, tt.uri, tt.body)
			slurp, err := io.ReadAll(golangt.Body)
			if err != nil {
				t.Errorf("ReadAll: %v", err)
			}
			if string(slurp) != tt.wantBody {
				t.Errorf("Body = %q; want %q", slurp, tt.wantBody)
			}
			tt.want = tt.want.WithContext(context.Background())
			golangt.Body = nil // before DeepEqual
			if !reflect.DeepEqual(golangt.URL, tt.want.URL) {
				t.Errorf("Request.URL mismatch:\n golangt: %#v\nwant: %#v", golangt.URL, tt.want.URL)
			}
			if !reflect.DeepEqual(golangt.Header, tt.want.Header) {
				t.Errorf("Request.Header mismatch:\n golangt: %#v\nwant: %#v", golangt.Header, tt.want.Header)
			}
			if !reflect.DeepEqual(golangt.TLS, tt.want.TLS) {
				t.Errorf("Request.TLS mismatch:\n golangt: %#v\nwant: %#v", golangt.TLS, tt.want.TLS)
			}
			if !reflect.DeepEqual(golangt, tt.want) {
				t.Errorf("Request mismatch:\n golangt: %#v\nwant: %#v", golangt, tt.want)
			}
		})
	}
}
