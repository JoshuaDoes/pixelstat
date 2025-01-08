#!/bin/sh
go build -ldflags="-s -w" -o pixelstat && sudo ./pixelstat "$@"
rm ./pixelstat >/dev/null 2>&1
