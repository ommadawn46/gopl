#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

 go run crawl.go -depth 3 "http://www.gopl.io/"
