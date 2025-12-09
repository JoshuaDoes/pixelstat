#!/bin/sh
go build -ldflags="-s -w" -o pixelstat -tags='ncurses' && ./pixelstat "$@"
rm ./pixelstat >/dev/null 2>&1
