#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

go run chat.go &
sleep 1s

echo "user1" | nc localhost 8000 &
sleep 0.2s
echo "user2" | nc localhost 8000 &
sleep 0.2s
echo "user3" | nc localhost 8000 &
sleep 0.2s
echo "user4" | nc localhost 8000 &
sleep 0.2s
echo "user5" | nc localhost 8000 &
sleep 0.2s
