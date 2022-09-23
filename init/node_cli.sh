#!/bin/bash

# Takes a single argument, the contents of the file you are uploading, and outputs the
# hash of the object, if it could be uploaded successfully.
function put(){
    echo "in put"
    source ~/.profile
    cd /home/gocode/src/put
    go run put.go $1
    cd /home/gocode/src/init
}

# Takes a hash as its only argument, and outputs the contents of the object and the
# node it was retrieved from, if it could be downloaded successfully.
function get(){
    echo "in get"
    source ~/.profile
    cd /home/gocode/src/get
    go run get.go $1
    cd /home/gocode/src/init
}

# Terminates the node (kill it)
function exit2(){
    echo "in exit2"
    source ~/.profile
    cd /home/gocode/src/exit
    go run exit.go $1
    cd /home/gocode/src/init
    kill -s SIGTERM 1
}

function ping2(){
    source ~/.profile
    cd /home/gocode/src/ping
    go run ping.go $1
    cd /home/gocode/src/init
}


# Tomorrow: Command exec package with bufio to recieve input in go with bufio and then forward to the terminal with exec