#!/bin/bash

# Takes a single argument, the contents of the file you are uploading, and outputs the
# hash of the object, if it could be uploaded successfully.
function put(){
    echo "test1"

}

# Takes a hash as its only argument, and outputs the contents of the object and the
# node it was retrieved from, if it could be downloaded successfully.
function get(){
    echo "test2"
    echo $1
}

# Terminates the node
function exit(){
    echo "test3"
}

cd docker
echo "test"