// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build ignore

#ifdef WIN32
// A Windows DLL is unable to call an arbitrary function in
// the main executable. Work around that by making the main
// executable pass the callback function pointer to us.
void (*golangCallback)(void);
__declspec(dllexport) void setCallback(void *f)
{
	golangCallback = (void (*)())f;
}
__declspec(dllexport) void sofunc(void);
#elif defined(_AIX)
// AIX doesn't allow the creation of a shared object with an
// undefined symbol. It's possible to bypass this problem by
// using -Wl,-G and -Wl,-brtl option which allows run-time linking.
// However, that's not how most of AIX shared object works.
// Therefore, it's better to consider golangCallback as a pointer and
// to set up during an init function.
void (*golangCallback)(void);
void setCallback(void *f) { golangCallback = f; }
#else
extern void golangCallback(void);
void setCallback(void *f) { (void)f; }
#endif

// OpenBSD and older Darwin lack TLS support
#if !defined(__OpenBSD__) && !defined(__APPLE__)
__thread int tlsvar = 12345;
#endif

void sofunc(void)
{
	golangCallback();
}
