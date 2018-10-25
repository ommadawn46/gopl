#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run chanpipe.go -n 8192
echo
go run chanpipe.go -n 16384
echo
go run chanpipe.go -n 32768
echo
go run chanpipe.go -n 65536
echo
go run chanpipe.go -n 131072
echo
go run chanpipe.go -n 262144
echo
go run chanpipe.go -n 524288
echo
go run chanpipe.go -n 1048576
echo
go run chanpipe.go -n 2097152
echo
go run chanpipe.go -n 4194304
echo
go run chanpipe.go -n 8388608
