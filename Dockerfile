FROM larjim/kademlialab:latest

# $ docker build . -t kadlab

# install stuff
# RUN apt-get update
# RUN apt-get -y install iputils-ping
# RUN apt -y install docker.io

# Copy dir D7024E from server to image
COPY gocode gocode

RUN echo "export GOPATH=/home/gocode" >> ~/.profile

RUN echo "source /home/gocode/src/init/node_cli.sh" >> ~/.bashrc

# Commands after this will be run in dir D7024E
WORKDIR gocode/src/init

RUN chmod +x /home/gocode/src/init/start_node.sh


# When creating a node it should run the start_node script, not sure if this works or not, probably not
ENTRYPOINT ["/home/gocode/src/init/start_node.sh"]