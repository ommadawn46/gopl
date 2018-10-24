#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

TZ=US/Eastern go run clock.go -port=8010 &
TZ=Asia/Tokyo go run clock.go -port=8020 &
TZ=Europe/London go run clock.go -port=8030 &
sleep 1s

go run clockwall.go NewYork=localhost:8010 Tokyo=localhost:8020 London=localhost:8030
