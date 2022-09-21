#!/bin/bash

source ~/.profile

# Runs the code in main file
IP=$(hostname -I)
./main ${IP}
echo "Exited node"