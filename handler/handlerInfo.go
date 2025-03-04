package handler

import (
	"bytes"
	"encoding/json"
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

		countryCode := r.PathValue("countryCode")

		limit := r.URL.Query().Get("limit")

		limitInt, err := strconv.Atoi(limit)
		if err != nil && limit != "" {
			log.Println("Could not convert limit to int:", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		reqUrl := urlRestCountries + fmt.Sprintf("alpha/%s?fields=name,continents,population,languages,borders,flags,capital", countryCode)

		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("Error fetching country data:", err)
			http.Error(w, "Failed to fetch country data", http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("Error closing body:", err)
			}
		}(resp.Body)

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		var countryInfo CountryInfo
		if err := json.Unmarshal(respBody, &countryInfo); err != nil {
			log.Println("Error parsing JSON:", err)
			http.Error(w, "Failed to parse country data", http.StatusInternalServerError)
			return
		}

		cities, err := getCities(countryInfo.Name.Common, limitInt)
		if err != nil {
			log.Println("Error fetching cities:", err)
			http.Error(w, "Failed to fetch cities", http.StatusInternalServerError)
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

		if err := json.NewEncoder(w).Encode(countryInfoFormatted); err != nil {
			log.Println("Error encoding response:", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}

// getCities
// Get cities from second api /*
func getCities(country string, limit int) ([]string, error) {
	requestBody, _ := json.Marshal(map[string]string{
		"country": country,
	})

	req, err := http.NewRequest("POST", urlCountriesNow+"countries/cities", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body:", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
	}

	var citiesResp CountryInfo

	if err := json.Unmarshal(respBody, &citiesResp); err != nil {
		log.Println("Error parsing JSON:", err)
	}

	if limit != 0 && len(citiesResp.Cities) > limit {
		citiesResp.Cities = citiesResp.Cities[:limit]
	}

	return citiesResp.Cities, nil
}
