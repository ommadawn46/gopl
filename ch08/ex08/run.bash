#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

go run reverb.go &
sleep 1s

echo -e "Yahoo!" | nc localhost 8000
