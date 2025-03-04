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
	w.Header().Set("Content-Type", "text/html")

	log.Printf("Received %s request on default handler.", r.Method)

	// Offer information for redirection to paths
	output := "This service does not provide any functionality on root path level. Please use paths <a href=\"" +
		InfoPath + "\">" + InfoPath + "</a> or <a href=\"" + PopulationPath + "\">" + PopulationPath +
		"</a> or <a href=\"" + StatusPath + "\">" + StatusPath + "</a>."

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
