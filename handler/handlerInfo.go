package handler

import (
	"log"
	"net/http"
)

// Info
// Handler for /info /*
func Info(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on /info handler.")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}
