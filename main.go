package main

import (
	"log"
	"net/http"
)

func main() {
	// Define the port we want our server to listen on
	const port = "8080"

	// Create a new ServeMux (router) to handle different HTTP paths
	mux := http.NewServeMux()

	// Initialize the server struct with the address and our router
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Print a message to the console so we know the server is starting
	log.Printf("Serving on port: %s\n", port)

	// Start the server and listen for incoming requests.
	// If the server fails to start, log.Fatal will exit the program.
	log.Fatal(srv.ListenAndServe())
}