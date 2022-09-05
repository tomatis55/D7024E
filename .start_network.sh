#!/bin/bash

# do not build

cd docker
sudo docker swarm init
sudo docker stack deploy 50 --compose-file "docker-compose.yml" 
sudo docker ps -a