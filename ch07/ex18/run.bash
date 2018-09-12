#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

xml='<hoge><fuga><piyo>foo</piyo></fuga><fuga>bar</fuga></hoge><hoge>baz</hoge>'
echo -n $xml | go run ./nodeTree.go
