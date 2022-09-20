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

# install stuff
# RUN apt-get update
# RUN apt-get -y install iputils-ping
# RUN apt -y install docker.io

# Copy dir gocode from server to image
COPY gocode/src/docker gocode/src/docker

RUN echo "export GOPATH=/home/gocode" >> ~/.profile

# Commands after this will be run in dir gocode/src/init
WORKDIR gocode/src/init

RUN chmod +x .start_node.sh

# When creating a node it should run the start_node script, not sure if this works or not, probably not
# Outcomment and run ./.start_node.sh manually to tets for sure
# CMD ["./.start_node.sh"]