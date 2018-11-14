#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run packagedeps.go archive/zip image/png
go run packagedeps.go crypto/sha...
