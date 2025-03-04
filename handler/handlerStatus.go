package handler

import (
	"assignment-1/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Status
// Handler for /status /*
func Status(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received %s request on /status handler.", r.Method)

	switch r.Method {
	case http.MethodGet:

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

		// CountriesNow
		statusCountriesNow, errCountriesNow := checkAPI(urlCountriesNow + "countries")
		if errCountriesNow != nil {
			log.Println("CountriesNow", errCountriesNow)
		}

		// REST Countries
		statusRestCountries, errRestCountries := checkAPI(urlRestCountries + "all?fields=name,flags")
		if errRestCountries != nil {
			log.Println("RestCountries", errRestCountries)
		}

		response := StatusResponse{
			CountriesNowAPI:  statusCountriesNow,
			RestCountriesAPI: statusRestCountries,
			Version:          util.Version,
			Uptime:           util.Uptime(),
		}

		w.Header().Set("Content-Type", "application/json")

		// Send the response
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error encoding JSON response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
		log.Println("Unsupported method received", r.Method)
		return
	}
}
