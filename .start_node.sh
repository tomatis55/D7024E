#!/bin/bash

# Runs the code in main file

go install
IP=$(hostname -I)
D7024E ${IP}

