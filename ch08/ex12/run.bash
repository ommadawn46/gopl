#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

go run chat.go &
sleep 1s

echo "Hoge 1" | nc localhost 8000 &
sleep 0.2s
echo "Hoge 2" | nc localhost 8000 &
sleep 0.2s
echo "Hoge 3" | nc localhost 8000 &
sleep 0.2s
echo "Hoge 4" | nc localhost 8000 &
sleep 0.2s
echo "Hoge 5" | nc localhost 8000 &
sleep 0.2s
