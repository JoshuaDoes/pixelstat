#!/bin/sh
go build -ldflags="-s -w" -o pixelstat && ./pixelstat "$@"
rm ./pixelstat >/dev/null 2>&1
