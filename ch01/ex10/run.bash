#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go run fetchall.go 'https://cve.mitre.org/data/downloads/allitems.html'
mv allitems.html allitems-1.html

go run fetchall.go 'https://cve.mitre.org/data/downloads/allitems.html'
mv allitems.html allitems-2.html

diff allitems-1.html allitems-2.html
