#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run countWordsAndImages.go 'https://ja.wikipedia.org/wiki/Go_(プログラミング言語)'
