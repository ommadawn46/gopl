#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run ../../ch01/fetch.go http://gopl.io/ch1/helloworld?go-get=1

# <meta name="go-import" content="gopl.io git https://github.com/adonovan/gopl.io">
