#!/bin/bash

# build the go code
cd /home/emacar-8/gocode/src/init
go build -o main  
chmod +x main

# source /home/emacar-8/join.sh
echo "source /home/emacar-8/join.sh" >> ~/.bashrc

# to build or not to build?
cd
sudo docker build . -t kadlab

sudo docker rm $(sudo docker ps -q -a) --force

sudo docker network rm net
sudo docker network create net --subnet=172.20.0.0/16     # gateway ip = 172.20.0.1

# add super node                          # super node ip = 172.20.0.2
sudo docker container run -it --ip 172.20.0.2 --net net --name "c0" kadlab      # -d

# loop to start nodes 
# for i in {1..3}
# do
#     echo "Creating node c$i"
#     sudo docker container run -it -d --net net --name "c${i}" kadlab
# done


# Guide:
# Run ./start_network.sh in server to start up the network
# If you want to create multiple nodes automatically, uncomment the for loop and add a -d flag on row 21 between -it and --ip

# To see a list of all nodes and their status: 
# sudo docker ps -a

# To look into an already created node:
# sudo docker attach "name of node"     for example sudo docker attach c1

# To start a new node in a new terminal that connects to the standard super node:
# sudo docker container run -it --net net --name "c1" kadlab 
# or
# join "name of container"          for example join c7 or join c52

# To start a new node in a new terminal that connects to a node of your choice:
# sudo docker container run -it --env CONN_TO_IP=172.20.0.2 --env CONN_TO_ID=0000000000000000000000000000000000000001 --net net --name "c97" kadlab
# or
# join "name of container" "ip of node to connect to" "id of node to connect to"