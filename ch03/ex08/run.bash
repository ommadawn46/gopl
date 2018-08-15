#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run mandelbrot.go > mandelbrot.png
echo 'output to mandelbrot.png'
