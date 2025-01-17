#!/bin/sh

BIN=$(basename $PWD)-ncurses
ZIP=$BIN-pkg.zip

go build -ldflags="-s -w" -o $BIN -tags='ncurses'
rm $ZIP
zip $ZIP -x "*.git*" -r -v .
rm $BIN
