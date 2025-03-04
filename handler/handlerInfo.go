package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Info
// Handler for /info /*
func Info(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		log.Printf("Received %s request on /info handler.", r.Method)

		countryCode := r.PathValue("countryCode")

		reqUrl := urlRestCountries + fmt.Sprintf("alpha/%s?fields=name,continents,population,languages,borders,flag,capital", countryCode)

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

		cities, err := getCities(countryInfo.Name.Common)
		if err != nil {
			log.Println("Error fetching cities:", err)
			http.Error(w, "Failed to fetch cities", http.StatusInternalServerError)
			return
		}

		if len(cities) > 10 {
			cities = cities[:10]
		}
		countryInfo.Cities = cities

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(countryInfo); err != nil {
			log.Println("Error encoding response:", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}

func getCities(country string) ([]string, error) {
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
			log.Println(err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var citiesResp CitiesResponse
	if err := json.Unmarshal(respBody, &citiesResp); err != nil {
		return nil, err
	}

	if len(citiesResp.Data) > 10 {
		citiesResp.Data = citiesResp.Data[:10]
	}

	return citiesResp.Data, nil
}
