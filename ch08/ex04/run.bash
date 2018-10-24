#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

go run reverb.go &
sleep 1s

echo -e "Hoge1\nHoge2\nHoge3\nHoge4\nHoge5\nHoge6\nHoge7\nHoge8\nHoge9\nHoge10" | go run ../ex03/netcat.go
