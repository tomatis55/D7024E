FROM larjim/kademlialab:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab

# RUN source ..node_cli.sh

RUN apt-get update
RUN apt-get -y install iputils-ping
RUN apt -y install docker.io
# add all code to node
# cd D7024E
# go install
# D7024E "id" "ip_address"

# go run main.go "id" "ip_address"  run start node file - skicka med id and ip address
# TODO use commit to create new image with everything already installed and files copied