#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run main.go demo/data/archive.zip demo/zip
echo
go run main.go demo/data/archive.tar demo/tar
