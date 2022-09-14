#!/bin/bash

# do not build

# cd docker
# sudo docker swarm init
# sudo docker stack deploy 50 --compose-file "docker-compose.yml" 

docker network rm net
sudo docker network create net      # gateway ip = 172.20.0.1

# add node
sudo docker container run -t --net net --name "c1" kadlab   # ip = 172.20.0.2

# loop to start nodes??
