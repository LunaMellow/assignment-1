package handler

import (
	"fmt"
	"log"
	"net/http"
)

// Empty
// Default handler /*
func Empty(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received %s request on default handler.", r.Method)

	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("Content-Type", "text/html")

	// Offer user information about paths and accepted parameters
	output := "This service does not provide any functionality on root path level. Please use paths:<br>" +
		"\n<a href=\"" + InfoPath + "\">" + InfoPath + "</a>" + "/{:two_letter_country_code}{?limit=10}<br>" +
		"\n<a href=\"" + PopulationPath + "\">" + PopulationPath + "</a>" + "/{:two_letter_country_code}{?limit={:startYear-endYear}}<br>" +
		"\n<a href=\"" + StatusPath + "\">" + StatusPath + "</a><br>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		log.Println("error in Empty handler: ", err)
		http.Error(w, "error when returning output", http.StatusInternalServerError)
		return
	}
}
