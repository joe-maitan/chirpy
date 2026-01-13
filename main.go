package main

import (
	// "os"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := "8080" // os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	// mux.Handle("/assets", http.FileServer(http.Dir("./assets/logo.png")))
	
	/* A http.Server is a struct that describes a server configuration */
	server := http.Server{
		Handler: mux,
		Addr: 	 "127.0.0.1:" + port,

	} 

	fmt.Printf("Server started on: %v...\n", port)
	
	/* ListenAndServe() blocks the main function until the server shuts down or an
	unexpected error crashes it. */
	log.Fatal(server.ListenAndServe())
} // End main() func