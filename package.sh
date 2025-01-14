#!/bin/sh

BIN=$(basename $PWD)
ZIP=$BIN-pkg.zip

go build -ldflags="-s -w" -o $BIN
rm $ZIP
zip $ZIP -x "*.git*" -r -v .
rm $BIN
