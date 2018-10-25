#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

echo "C" | go run crawl.go "https://tour.golang.org" "http://www.gopl.io/" "https://golang.org/" "https://github.com/golang"
