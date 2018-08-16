#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

export EDITOR="vim"
export GH_USER="GITHUB_USERNAME"
export GH_PASS="GITHUB_PASSWORD"

go run main.go -a read -o golang -r go
#go run main.go -a create -o $GH_USER -r test
#go run main.go -a edit -o $GH_USER -r test -n 1
#go run main.go -a close -o $GH_USER -r test -n 1
