#!/bin/bash

# join "name of container" "ip of node to connect to" "id of node to connect to"
function join(){
    echo $1
    sudo docker container run -it --name $1 --env CONN_TO_IP=$2 --env CONN_TO_ID=$3 --net net  kadlab
}
