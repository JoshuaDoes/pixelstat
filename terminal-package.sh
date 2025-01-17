#!/bin/sh

BIN=$(basename $PWD)-terminal
ZIP=$BIN-pkg.zip

go build -ldflags="-s -w" -o $BIN -tags='terminal'
rm *.zip
zip $ZIP -x "*.git*" -r -v .
rm $BIN
