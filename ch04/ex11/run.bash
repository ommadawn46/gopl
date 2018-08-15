#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

export EDITOR="vim"
export GH_USER="GITHUB_USERNAME"
export GH_PASS="GITHUB_PASSWORD"

go run main.go -a read -o golang -r go
