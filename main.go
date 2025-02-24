package main

import (
	"assignment-1/handler"
	"log"
	"net/http"
	"os"
)

func main() {

	// Extract PORT variable from the environment variables
	port := os.Getenv("PORT")

	// Override port with default port if not provided (e.g. local deployment)
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Instantiate the router
	router := http.NewServeMux()

	// Set up handler endpoints
	router.HandleFunc(handler.DEFAULT_PATH, handler.Empty)
	router.HandleFunc(handler.INFO_PATH, handler.Info)
	router.HandleFunc(handler.POPULATION_PATH, handler.Population)
	router.HandleFunc(handler.STATUS_PATH, handler.Status)

	// Start HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, router))

}
