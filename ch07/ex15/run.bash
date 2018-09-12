#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

echo '(A + 3) * B'
echo 'A=2 B=5'
cat << EOS | go run evalCli.go
(A + 3) * B
2
5
EOS

echo -e "\nreduce(max, A, 10, B, 5, C)"
echo 'A=1 B=100 C=20'
cat << EOS | go run evalCli.go
reduce(max, A, 10, B, 5, C)
1
100
20
EOS

echo -e "\nreduce(mul, A, 3, B) + reduce(mul, 7, C, 13)"
echo 'A=2 B=5 C=11'
cat << EOS | go run evalCli.go
reduce(mul, A, 3, B) + reduce(mul, 7, C, 13)
2
5
11
EOS
