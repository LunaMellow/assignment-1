package main

import (
	"assignment-1/handlers"
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

	// /{:two_letter_country_code}{?limit=10}
	router.HandleFunc("/info", handlers.HandlerInfo)
	// /{:two_letter_country_code}{?limit={:startYear-endYear}}
	router.HandleFunc("/population", handlers.HandlerPopulation)
	router.HandleFunc("/status", handlers.HandlerStatus)

	// Start HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, router))

}
