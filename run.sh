#!/bin/bash

# set the environmental variables
set -a
source .env
set +a

# clean previously compiled files
go clean

goose postgres $DB_URL down
goose postgres $DB_URL up

sqlc generate

# build a new go exec
go build .

# run the program
./chirpy
