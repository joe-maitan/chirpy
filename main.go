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
	
	/* A http.Server is a struct that describes a server configuration */
	server := http.Server{
		Handler: mux,
		Addr: 	 "127.0.0.1:" + port,

	} 

	mux.Handle("/", http.FileServer(http.Dir(".")))
	
	fmt.Printf("Server started on: %v...\n", port)
	
	/* ListenAndServe() blocks the main function until the server shuts down or an
	unexpected error crashes it. */
	err := server.ListenAndServe()
	log.Fatal(err)
} // End main() func