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
	output := "This service does not provide any functionality on root path level. Please use paths:" +
		"\n<a href=\"" + InfoPath + "\">" + InfoPath + "</a>" + "/{:two_letter_country_code}{?limit=10}<br>" +
		"\n<a href=\"" + PopulationPath + "\">" + PopulationPath + "</a>" + "/{:two_letter_country_code}{?limit={:startYear-endYear}}<br>" +
		"\n<a href=\"" + StatusPath + "\">" + StatusPath + "</a><br>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
