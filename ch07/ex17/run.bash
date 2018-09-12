#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

url='http://www.w3.org/TR/2006/REC-xml11-20060816'
go run ../../ch01/fetch.go $url | go run ./xmlselect.go back div1 h2
