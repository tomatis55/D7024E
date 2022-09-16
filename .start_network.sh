#!/bin/bash

# to build or not to build?
sudo docker build . -t kadlab

sudo docker rm c0 c1 --force

sudo docker network rm net
sudo docker network create net      # gateway ip = 172.20.0.1

# add super node                          # super node ip = 172.20.0.2
sudo docker container run -it -d --net net --name "c0" kadlab

# loop to start nodes - not completely done yet
# for i in {1..2}
# do
#     echo "Welcome $i times"
#     sudo docker container run -it -d --net net --name "c${i}" kadlab
# done



# Guide:
# Run ./.start_network.sh in server
# In the container terminal, run ./.start_node.sh

# Change the content in main.go
# In a new server terminal, run the following:
# sudo docker container run -it -d --net net --name "c1" kadlab
# sudo docker exec -it c1 /bin/bash
# In the node terminal, run ./.start_node.sh
