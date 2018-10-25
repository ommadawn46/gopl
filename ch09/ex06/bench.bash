#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

echo -e "GOMAXPROCS 1"
GOMAXPROCS=1 go test -benchmem -bench .
echo -e "\nGOMAXPROCS 2"
GOMAXPROCS=2 go test -benchmem -bench .
echo -e "\nGOMAXPROCS 4"
GOMAXPROCS=4 go test -benchmem -bench .
echo -e "\nGOMAXPROCS 8"
GOMAXPROCS=8 go test -benchmem -bench .
echo -e "\nGOMAXPROCS 16"
GOMAXPROCS=16 go test -benchmem -bench .
echo -e "\nGOMAXPROCS 32"
GOMAXPROCS=32 go test -benchmem -bench .
