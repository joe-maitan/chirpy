#!/bin/bash

# Set boot.dev url to be ferrari:8080/app or localhost:8080/app

set -a
source .env
set +a

go clean

go build .

./chirpy
