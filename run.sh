#!/bin/bash

# set the environmental variables
set -a
source .env
set +a

# clean previously compiled files
go clean

# build a new go exec
go build .

# run the program
./chirpy
