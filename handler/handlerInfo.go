package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Info
// Handler for /info /*
func Info(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received %s request on /info handler.", r.Method)

	switch r.Method {
	case http.MethodGet:

		// Get parameters
		countryCode := r.PathValue("countryCode")
		limit := r.URL.Query().Get("limit")

		// Convert limit parameter to int
		limitInt, err := strconv.Atoi(limit)
		if limitInt <= 0 || err != nil && limit != "" {
			log.Println("could not convert limit to int:", err)
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		// Api request url
		reqUrl := urlRestCountries + fmt.Sprintf("alpha/%s?fields=name,continents,population,languages,borders,flags,capital", countryCode)

		// Fetch country data
		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("error fetching country data:", err)
			http.Error(w, "failed to fetch country data", http.StatusInternalServerError)
			return
		}

		// Defer close body and throw error if any
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("error closing body:", err)
				return
			}
		}(resp.Body)

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("error reading response body:", err)
			http.Error(w, "failed to read response", http.StatusInternalServerError)
			return
		}

		// Parse json country data
		var countryInfo CountryInfo
		if err := json.Unmarshal(respBody, &countryInfo); err != nil {
			log.Println("error parsing JSON:", err)
			http.Error(w, "failed to parse country data", http.StatusInternalServerError)
			return
		}

		// Get cities from second api
		cities, err := getCities(countryInfo.Name.Common, limitInt)
		if err != nil {
			log.Println("error fetching cities:", err)
			http.Error(w, "failed to fetch cities", http.StatusInternalServerError)
			return
		}

		countryInfoFormatted := CountryInfoFormatted{
			Name:       countryInfo.Name.Common,
			Continents: countryInfo.Continents,
			Population: countryInfo.Population,
			Languages:  countryInfo.Languages,
			Borders:    countryInfo.Borders,
			Flags:      countryInfo.Flags.Png,
			Capital:    countryInfo.Capital[0],
			Cities:     cities,
		}

		w.Header().Set("Content-Type", "application/json")

		// Encode json response
		if err := json.NewEncoder(w).Encode(countryInfoFormatted); err != nil {
			log.Println("error encoding response:", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

	default:
		log.Println("unsupported method received", r.Method)
		http.Error(w, "method not supported", http.StatusNotImplemented)
		return
	}
}

// getCities
// Get cities from second api /*
func getCities(country string, limit int) ([]string, error) {

	// Api post request body
	requestBody, _ := json.Marshal(map[string]string{
		"country": country,
	})

	// Send post request with provided body
	req, err := http.NewRequest("POST", urlCountriesNow+"countries/cities", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("error creating request:", err)
		return nil, errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	// Start client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error fetching countries data:", err)
		return nil, errors.New("failed to fetch countries data")
	}

	// Defer close body and throw error if any
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing body:", err)
		}
	}(resp.Body)

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
		return nil, errors.New("failed to read response body")
	}

	// Parse json countries data
	var citiesResp CountryInfo
	if err := json.Unmarshal(respBody, &citiesResp); err != nil {
		log.Println("error parsing JSON:", err)
		return nil, errors.New("failed to parse country data")
	}

	// Limit cities to provided user limit
	if limit != 0 && len(citiesResp.Cities) > limit {
		citiesResp.Cities = citiesResp.Cities[:limit]
	}

	return citiesResp.Cities, nil
}
