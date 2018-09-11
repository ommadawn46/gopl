#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run ../../ch01/fetch.go http://gopl.io | go run main.go
