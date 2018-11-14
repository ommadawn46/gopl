#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

cat demo/gopher_org.png | go run imgconv.go -f jpg > demo/gopher.jpg
cat demo/gopher.jpg | go run imgconv.go -f gif > demo/gopher.gif
cat demo/gopher.gif | go run imgconv.go -f png > demo/gopher.png
