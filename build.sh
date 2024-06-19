#!/bin/bash

if [ $1 == "" ] || [ $1 == "MAC" ]; then
    gox -os "darwin" -arch "amd64" -cgo
    mv acmweb_darwin_amd64 ../../bin/acmsvr_darwin
fi
if [ $1 == "" ] || [ $1 == "LINUX" ]; then
    CC=x86_64-unknown-linux-gnu-gcc gox -os "linux" -arch "amd64" -cgo
    mv acmweb_linux_amd64 ../../bin/acmsvr
fi
if [ $1 == "" ] || [ $1 == "WIN64" ]; then
    CC=x86_64-w64-mingw32-gcc gox -os "windows" -arch "amd64" -cgo
    mv acmweb_windows_amd64.exe ../../bin/acmsvr.exe
fi
