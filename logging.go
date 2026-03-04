package main

import (
	"log"
)

func logStatement(fileName, functionName, errorMessage string, errorObj error) {
	log.Printf("%v - %v - %v: %v", fileName, functionName, errorMessage, errorObj)
} // End logStatement() func