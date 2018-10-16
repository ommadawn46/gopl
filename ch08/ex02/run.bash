#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

## Install FTP client (MacOS High Sierra or higher)
# brew install inetutils

rm -rf ./server ./client
mkdir -p ./server ./client

echo "SERVER FILE" > ./server/serverfile.txt
echo "CLIENT FILE" > ./client/clientfile.txt

cp ./img/image1.png ./server/serverfile.png
cp ./img/image2.png ./client/clientfile.png

go run main.go -port 22221 -root ./server &
sleep 3s

cd ./client

## PASSIVE MODE
cat <<EOF | ftp -n -p &
open 127.0.0.1 22221
user test1 qwerty

ascii
get serverfile.txt downloaded1.txt
put clientfile.txt uploaded1.txt

binary
get serverfile.png downloaded1.png
put clientfile.png uploaded1.png

size uploaded1.txt
size uploaded1.png

modtime uploaded1.txt
modtime uploaded1.png
bye
EOF

## ACTIVE MODE
cat <<EOF | ftp -n &
open 127.0.0.1 22221
user test2 abcd1234

ascii
get serverfile.txt downloaded2.txt
put clientfile.txt uploaded2.txt

binary
get serverfile.png downloaded2.png
put clientfile.png uploaded2.png

size uploaded2.txt
size uploaded2.png

modtime uploaded2.txt
modtime uploaded2.png
bye
EOF

## DIRECTORY OPERATIONS
cat <<EOF | ftp -n -p &
open 127.0.0.1 22221
user test3 p455w0rd

mkdir dir1
mkdir dir1/dir2
mkdir dir1/dir2/dir3
ls
cd dir1/dir2/dir3
pwd
put clientfile.png uploaded.png
nlist
rename uploaded.png renamed.png
nlist
delete renamed.png
nlist
cd /
rmdir dir1/dir2/dir3
rmdir dir1/dir2
rmdir dir1
ls
bye
EOF

sleep 2s
cd ..

echo -e "\n$ ls -la ./server"
ls -la ./server

echo -e "\n$ ls -la ./client"
ls -la ./client
