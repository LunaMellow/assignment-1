package handler

import (
	"fmt"
	"log"
	"net/http"
)

// Empty
// Default handler /*
func Empty(w http.ResponseWriter, r *http.Request) {

	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	log.Println("Received " + r.Method + " request on default handler.")

	// Offer information for redirection to paths
	output := "This service does not provide any functionality on root path level. Please use paths <a href=\"" +
		INFO_PATH + "\">" + INFO_PATH + "</a> or <a href=\"" + POPULATION_PATH + "\">" + POPULATION_PATH +
		"</a> or <a href=\"" + STATUS_PATH + "\">" + STATUS_PATH + "</a>."

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
