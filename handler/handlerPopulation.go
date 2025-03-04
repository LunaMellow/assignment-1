package handler

import (
	"log"
	"net/http"
)

// Population
// Handler for /population /*
func Population(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received %s request on /population handler.", r.Method)

	switch r.Method {
	case http.MethodGet:

		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}
