#!/bin/sh
go build -ldflags="-s -w" -o pixelstat && sudo ./pixelstat "$@"
