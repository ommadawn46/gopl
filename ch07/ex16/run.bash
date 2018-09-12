#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

go run evalServer.go &
sleep 3s

curl 'http://localhost:8000/?expr=(2+%2B+3)+*+5'; echo
curl 'http://localhost:8000/?expr=reduce(mul,+5,+3,+10,+20,+50)'; echo
curl 'http://localhost:8000/?expr=reduce(add,+5,+3,+17,+31)+*+reduce(add,+7,+11,+13,+27)'; echo
