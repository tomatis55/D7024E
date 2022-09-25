#!/bin/bash

source ~/.profile

# Runs the code in main file

IP=$(hostname -I)

if [ ! -z "$CONN_TO_IP" ]
then
    if [ ! -z "$CONN_TO_ID" ]
    then
    # echo "in double if"
    # echo ${CONN_TO_IP}
    # echo ${CONN_TO_ID}
    ./main ${IP} ${CONN_TO_ID} ${CONN_TO_IP}
    fi
else
    # echo "in else"
    ./main ${IP}
fi

echo "Exited node"