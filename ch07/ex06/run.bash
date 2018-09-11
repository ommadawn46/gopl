#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run tempflag.go -temp 0K
go run tempflag.go -temp 0F
go run tempflag.go -temp 100K
go run tempflag.go -temp 100F
