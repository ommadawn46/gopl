#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run fetch.go "https://tour.golang.org" "http://www.gopl.io/" "https://golang.org/" "https://github.com/golang"
