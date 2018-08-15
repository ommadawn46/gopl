#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

echo -n 'sha256: ' && go run shadigest.go -m ABCDEFGH
echo -n 'sha384: ' && go run shadigest.go -a sha384 -m ABCDEFGH
echo -n 'sha512: ' && go run shadigest.go -a sha512 -m ABCDEFGH
