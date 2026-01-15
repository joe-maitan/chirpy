package main

import (
	// "os"
	"fmt"
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	port := "8080" // os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	// mux.Handle("/assets", http.FileServer(http.Dir("./assets/logo.png")))
	
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handlerReadiness)

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