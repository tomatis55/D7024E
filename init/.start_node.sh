#!/bin/bash

# Runs the code in main file
source ~/.profile
IP=$(hostname -I)
go run main.go ${IP}
echo "Exited node" 