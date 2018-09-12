#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

go run stockServer.go &
sleep 3s

curl 'http://localhost:8000/list'; echo
curl 'http://localhost:8000/create?item=hoge&price=100'; echo
curl 'http://localhost:8000/list'; echo
curl 'http://localhost:8000/update?item=hoge&price=200'; echo
curl 'http://localhost:8000/list'; echo
curl 'http://localhost:8000/delete?item=hoge'; echo
curl 'http://localhost:8000/list'
