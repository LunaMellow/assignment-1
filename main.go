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
	router.HandleFunc("/", handler.Empty)
	router.HandleFunc(handler.DefaultPath, handler.Empty)
	router.HandleFunc(handler.InfoPath+"/{countryCode}", handler.Info)
	router.HandleFunc(handler.PopulationPath+"/{countryCode}", handler.Population)
	router.HandleFunc(handler.StatusPath, handler.Status)

	log.Println("Starting server on port " + port)
	serveErr := http.ListenAndServe(":"+port, router)
	if serveErr != nil {

		log.Fatal(serveErr.Error())
	}

}
