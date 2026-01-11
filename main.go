package main

import (
	// "os"
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	port := "8080" // os.Getenv("PORT")
	
	server := http.Server{
		Handler: mux,
		Addr: 	 ":" + port,

	} 
	
	fmt.Printf("Server started on: %v...\n", port)
	err := server.ListenAndServe()
	log.Fatal(err)
} // End main() func