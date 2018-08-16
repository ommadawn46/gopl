#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run xkcd.go programming language go
go run xkcd.go 2012/8/10
