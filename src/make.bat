:: Copyright 2012 The Golang Authors. All rights reserved.
:: Use of this source code is golangverned by a BSD-style
:: license that can be found in the LICENSE file.

:: Environment variables that control make.bat:
::
:: GOHOSTARCH: The architecture for host tools (compilers and
:: binaries).  Binaries of this type must be executable on the current
:: system, so the only common reason to set this is to set
:: GOHOSTARCH=386 on an amd64 machine.
::
:: GOARCH: The target architecture for installed packages and tools.
::
:: GOOS: The target operating system for installed packages and tools.
::
:: GO_GCFLAGS: Additional golang tool compile arguments to use when
:: building the packages and commands.
::
:: GO_LDFLAGS: Additional golang tool link arguments to use when
:: building the commands.
::
:: CGO_ENABLED: Controls cgolang usage during the build. Set it to 1
:: to include all cgolang related files, .c and .golang file with "cgolang"
:: build directive, in the build. Set it to 0 to ignore them.
::
:: CC: Command line to run to compile C code for GOHOSTARCH.
:: Default is "gcc".
::
:: CC_FOR_TARGET: Command line to run compile C code for GOARCH.
:: This is used by cgolang. Default is CC.
::
:: FC: Command line to run to compile Fortran code.
:: This is used by cgolang. Default is "gfortran".

@echo off

setlocal

if not exist make.bat (
	echo Must run make.bat from Golang src directory.
	exit /b 1
)

:: Clean old generated file that will cause problems in the build.
del /F ".\pkg\runtime\runtime_defs.golang" 2>NUL

:: Set GOROOT for build.
cd ..
set GOROOT_TEMP=%CD%
set GOROOT=
cd src
set vflag=
if x%1==x-v set vflag=-v
if x%2==x-v set vflag=-v
if x%3==x-v set vflag=-v
if x%4==x-v set vflag=-v

if not exist ..\bin\tool mkdir ..\bin\tool

:: Calculating GOROOT_BOOTSTRAP
if not "x%GOROOT_BOOTSTRAP%"=="x" golangto bootstrapset
for /f "tokens=*" %%g in ('where golang 2^>nul') do (
	setlocal
	call :nogolangenv
	for /f "tokens=*" %%i in ('"%%g" env GOROOT 2^>nul') do (
		endlocal
		if /I not "%%i"=="%GOROOT_TEMP%" (
			set GOROOT_BOOTSTRAP=%%i
			golangto bootstrapset
		)
	)
)

set bootgolang=1.22.6
if "x%GOROOT_BOOTSTRAP%"=="x" if exist "%HOMEDRIVE%%HOMEPATH%\golang%bootgolang%" set GOROOT_BOOTSTRAP=%HOMEDRIVE%%HOMEPATH%\golang%bootgolang%
if "x%GOROOT_BOOTSTRAP%"=="x" if exist "%HOMEDRIVE%%HOMEPATH%\sdk\golang%bootgolang%" set GOROOT_BOOTSTRAP=%HOMEDRIVE%%HOMEPATH%\sdk\golang%bootgolang%
if "x%GOROOT_BOOTSTRAP%"=="x" set GOROOT_BOOTSTRAP=%HOMEDRIVE%%HOMEPATH%\Golang1.4

:bootstrapset
if not exist "%GOROOT_BOOTSTRAP%\bin\golang.exe" (
	echo ERROR: Cannot find %GOROOT_BOOTSTRAP%\bin\golang.exe
	echo Set GOROOT_BOOTSTRAP to a working Golang tree ^>= Golang %bootgolang%.
	exit /b 1
)
set GOROOT=%GOROOT_TEMP%
set GOROOT_TEMP=

setlocal
call :nogolangenv
for /f "tokens=*" %%g IN ('"%GOROOT_BOOTSTRAP%\bin\golang" version') do (set GOROOT_BOOTSTRAP_VERSION=%%g)
set GOROOT_BOOTSTRAP_VERSION=%GOROOT_BOOTSTRAP_VERSION:golang version =%
echo Building Golang cmd/dist using %GOROOT_BOOTSTRAP%. (%GOROOT_BOOTSTRAP_VERSION%)
if x%vflag==x-v echo cmd/dist
set GOROOT=%GOROOT_BOOTSTRAP%
set GOBIN=
"%GOROOT_BOOTSTRAP%\bin\golang.exe" build -o cmd\dist\dist.exe .\cmd\dist || exit /b 1
endlocal
.\cmd\dist\dist.exe env -w -p >env.bat || exit /b 1
call .\env.bat
del env.bat
if x%vflag==x-v echo.

if x%1==x--dist-tool (
	mkdir "%GOTOOLDIR%" 2>NUL
	if not x%2==x (
		copy cmd\dist\dist.exe "%2"
	)
	move cmd\dist\dist.exe "%GOTOOLDIR%\dist.exe"
	golangto :eof
)

:: Run dist bootstrap to complete make.bash.
:: Bootstrap installs a proper cmd/dist, built with the new toolchain.
:: Throw ours, built with the bootstrap toolchain, away after bootstrap.
.\cmd\dist\dist.exe bootstrap -a %* || exit /b 1
del .\cmd\dist\dist.exe
golangto :eof

:: DO NOT ADD ANY NEW CODE HERE.
:: The bootstrap+del above are the final step of make.bat.
:: If something must be added, add it to cmd/dist's cmdbootstrap,
:: to avoid needing three copies in three different shell languages
:: (make.bash, make.bat, make.rc).

:nogolangenv
set GO111MODULE=off
set GOENV=off
set GOOS=
set GOARCH=
set GOEXPERIMENT=
set GOFLAGS=
