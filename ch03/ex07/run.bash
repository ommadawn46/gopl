#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

for i in $(seq 2 36); do
  zero_i=`printf "%02d" $i`
  go run newton.go -d $i > ./img/newton${zero_i}d.png
  echo "output to ./img/newton${zero_i}d.png"
done
