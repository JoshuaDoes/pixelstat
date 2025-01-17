#!/bin/sh
go build -ldflags="-s -w" -o pixelstat -tags='terminal' && sudo ./pixelstat "$@"
rm ./pixelstat >/dev/null 2>&1
