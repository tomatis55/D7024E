#!/bin/bash

# to build or not to build?
cd
sudo docker build /home/romerm-8/gocode/ -t kadlab

sudo docker rm $(sudo docker ps -q -a) --force

sudo docker network rm net
sudo docker network create net --subnet=172.20.0.0/16     # gateway ip = 172.20.0.1

# add super node                          # super node ip = 172.20.0.2
sudo docker container run -it --ip 172.20.0.2 --net net --name "c0" kadlab

# loop to start nodes - not completely done yet
# for i in {1..2}
# do
#     echo "Welcome $i times"
#     sudo docker container run -it -d --net net --name "c${i}" kadlab
# done



# Guide:
# Run /home/romerm-8/gocode/src/init/.start_network.sh in server
# In the container terminal, run ./.start_node.sh

# In a new server terminal, run the following:
# sudo docker container run -it -d --net net --name "c1" kadlab
# sudo docker exec -it c1 /bin/bash
# In the node terminal, run ./.start_node.sh