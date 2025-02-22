package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Application version
var version = "v1"

// Application timer
var startTime = time.Now()

type StatusResponse struct {
	CountriesNowAPI  int     `json:"countriesnowapi"`
	RestCountriesAPI int     `json:"restcountriesapi"`
	Version          string  `json:"version"`
	Uptime           float64 `json:"uptime"`
}

func HandlerInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on /info handler.")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}

func HandlerPopulation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on /population handler.")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}

func HandlerStatus(w http.ResponseWriter, r *http.Request) {

	// Method GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported", http.StatusNotImplemented)
		return
	}

	urlCountriesNow := "http://129.241.150.113:3500/api/v0.1/countries"          // Default
	urlRestCountries := "http://129.241.150.113:8080/v3.1/all?fields=name,flags" // Filtered for faster response

	log.Println("Received " + r.Method + " request on /status handler.")

	// Create a client with a timeout
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()

	checkAPI := func(url string) (int, error) {

		// Get request
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error creating request: %v", err)
		}

		req.Header.Add("Content-Type", "application/json")

		// Send client request
		res, err := client.Do(req)
		if err != nil {
			return http.StatusServiceUnavailable, fmt.Errorf("error contacting API: %v", err)
		}
		// Close client at end of scope
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("error closing response body: %v", err)
			}
		}(res.Body)

		return res.StatusCode, nil
	}

	// Get API status codes
	statusCountriesNow, errCountriesNow := checkAPI(urlCountriesNow)
	if errCountriesNow != nil {
		log.Println("CountriesNow", errCountriesNow)
	}
	statusRestCountries, errRestCountries := checkAPI(urlRestCountries)
	if errRestCountries != nil {
		log.Println("RestCountries", errRestCountries)
	}

	uptime := time.Since(startTime).Round(time.Second)

	response := StatusResponse{
		CountriesNowAPI:  statusCountriesNow,
		RestCountriesAPI: statusRestCountries,
		Version:          version,
		Uptime:           uptime.Seconds(),
	}

	w.Header().Set("Content-Type", "application/json")

	// Send the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
