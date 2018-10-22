#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

## Install FTP client (MacOS High Sierra or higher)
# brew install inetutils

rm -rf ./demo/server ./demo/client ./demo/passwd

mkdir -p ./demo/server ./demo/client
touch ./demo/passwd

echo -e "SERVER FILE\nSERVER FILE\nSERVER FILE\nSERVER FILE\nSERVER FILE" > ./demo/server/serverfile.txt
echo -e "CLIENT FILE\nCLIENT FILE\nCLIENT FILE\nCLIENT FILE\nCLIENT FILE" > ./demo/client/clientfile.txt

cp ./demo/img/image1.png ./demo/server/serverfile.png
cp ./demo/img/image2.png ./demo/client/clientfile.png

## Create Users
go run main.go -passwd ./demo/passwd -adduser -user test1 -pass qwerty
go run main.go -passwd ./demo/passwd -adduser -user test2 -pass abcd1234
go run main.go -passwd ./demo/passwd -adduser -user test3 -pass p455w0rd
go run main.go -passwd ./demo/passwd -adduser -user test4 -pass letmein

## Start FTP Server
go run main.go -port 22221 -root ./demo/server -passwd ./demo/passwd &
sleep 3s

cd ./demo/client

## HELP MESSAGES
cat <<EOF | ftp -n &
open 127.0.0.1 22221
user test1 qwerty

rhelp
rhelp ABOR
rhelp CWD
rhelp DELE
rhelp EPRT
rhelp EPSV
rhelp HELP
rhelp LIST
rhelp MDTM
rhelp MKD
rhelp MODE
rhelp NLST
rhelp NOOP
rhelp PASS
rhelp PASV
rhelp PORT
rhelp PWD
rhelp QUIT
rhelp RETR
rhelp RMD
rhelp RNFR
rhelp RNTO
rhelp SIZE
rhelp STOR
rhelp STRU
rhelp TYPE
rhelp USER
bye
EOF

## PASSIVE MODE
cat <<EOF | ftp -n -p &
open 127.0.0.1 22221
user test2 abcd1234

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
user test3 p455w0rd

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
user test4 letmein

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
cd ../..

echo -e "\n$ ls -la ./demo/server"
ls -la ./demo/server

echo -e "\n$ ls -la ./demo/client"
ls -la ./demo/client
