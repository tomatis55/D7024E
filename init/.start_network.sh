#!/bin/bash

# build the go code
source /home/romerm-8/.profile

cd /home/romerm-8/gocode/src/init
go build -o main.go  
chmod +x main.go

source join.sh

# to build or not to build?
cd
sudo docker build /home/romerm-8/gocode/ -t kadlab

sudo docker rm $(sudo docker ps -q -a) --force

sudo docker network rm net
sudo docker network create net --subnet=172.20.0.0/16     # gateway ip = 172.20.0.1

# add super node                          # super node ip = 172.20.0.2
sudo docker container run -it --ip 172.20.0.2 --net net --name "c0" kadlab

# loop to start nodes 
# for i in {1..3}
# do
#     echo "Creating node c$i"
#     sudo docker container run -it -d --net net --name "c${i}" kadlab
# done

# sudo docker container run -it --net net --name "c1" kadlab

# sudo docker container run -it --env CONN_TO_IP=172.20.0.3 --env CONN_TO_ID=0000000000000000000000000000000000000001 --net net --name "c97" kadlab

# Guide:
# Run /home/romerm-8/gocode/src/init/.start_network.sh in server
# In the container terminal, run ./.start_node.sh

# In a new server terminal, run the following:
# sudo docker container run -it -d --net net --name "c1" kadlab
# sudo docker exec -it c1 /bin/bash
# In the node terminal, run ./.start_node.sh