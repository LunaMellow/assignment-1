package handler

import (
	"log"
	"net/http"
)

// Population
// Handler for /population /*
func Population(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on /population handler.")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}
