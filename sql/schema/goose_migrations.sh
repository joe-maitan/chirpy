#!/bin/bash

# Make sure to export $GOOSE_DB_URL before running this script

goose postgres $GOOSE_DB_URL down
goose postgres $GOOSE_DB_URL up

