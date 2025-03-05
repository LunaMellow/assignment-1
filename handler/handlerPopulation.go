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
	"strings"
)

// Population
// Handler for /population /*
func Population(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received %s request on /population handler.", r.Method)

	switch r.Method {
	case http.MethodGet:

		// Get parameters
		countryCode := r.PathValue("countryCode")
		limit := r.URL.Query().Get("limit")

		// Get country name from code
		country, err := getCountry(countryCode)
		if err != nil {
			log.Println("error getting country name:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Api post request body
		requestBody, _ := json.Marshal(map[string]string{
			"country": country,
		})

		// Send post request with provided body
		req, err := http.NewRequest("POST", urlCountriesNow+"countries/population", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Println("error creating request:", err)
			http.Error(w, "failed create request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		// Start client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("error fetching countries data:", err)
			http.Error(w, "failed fetch countries data", http.StatusInternalServerError)
			return
		}

		// Defer close body and throw error if any
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("error closing body:", err)
				http.Error(w, "failed to close body", http.StatusInternalServerError)
				return
			}
		}(resp.Body)

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("error reading response body:", err)
			http.Error(w, "failed to read response body", http.StatusInternalServerError)
			return
		}

		// Parse json countries data
		var countryResp CountryPopulation
		if err := json.Unmarshal(respBody, &countryResp); err != nil {
			log.Println("error parsing JSON:", err)
			http.Error(w, "failed to parse response body", http.StatusInternalServerError)
			return
		}

		// Filter population if limit specified, otherwise default to show all
		var filteredPopulation CountryPopulation
		var population CountryPopulation
		if limit != "" {
			var err error
			filteredPopulation, err = filterPopulation(countryResp, limit)
			if err != nil {
				log.Println("error filtering population:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			population = filteredPopulation
		} else {
			population.Data.PopulationCounts = countryResp.Data.PopulationCounts
		}

		// Count up mean average of population count
		mean := calculateMean(population)

		countryPopulationFormatted := CountryPopulationFormatted{
			Mean:   mean,
			Values: population.Data.PopulationCounts,
		}

		w.Header().Set("Content-Type", "application/json")

		// Encode json response
		if err := json.NewEncoder(w).Encode(countryPopulationFormatted); err != nil {
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

// calculateMean
// Calculate the mean average population
func calculateMean(population CountryPopulation) int {

	var mean int
	var totalPopulation int

	// Check if any population at all
	if len(population.Data.PopulationCounts) == 0 {
		mean = 0
	} else {

		// Loop through and calculate mean average population
		for _, count := range population.Data.PopulationCounts {
			totalPopulation += count.Value
		}
		mean = totalPopulation / len(population.Data.PopulationCounts)
	}
	return mean
}

// filterPopulation
// Filter the population data from the limits provided /*
func filterPopulation(population CountryPopulation, limit string) (CountryPopulation, error) {

	// Split limit string provided by user
	limitRange := strings.Split(limit, "-")
	if len(limitRange) != 2 || limitRange[0] == "" || limitRange[1] == "" {
		log.Println("invalid limit range:", limitRange)
		return CountryPopulation{}, errors.New("invalid limit range. Format (YYYY-YYYY)")
	}

	// Convert startYear to int
	startYear, err := strconv.Atoi(limitRange[0])
	if err != nil || startYear < 1960 || startYear > 2018 {
		log.Println("invalid start year parameter. Format (YYYY-YYYY)")
		return CountryPopulation{}, errors.New("invalid start range. Format (YYYY-YYYY)")
	}

	// Convert endYear to int
	endYear, err := strconv.Atoi(limitRange[1])
	if err != nil || endYear < startYear || endYear > 2018 {
		log.Println("invalid end parameter. Format (YYYY-YYYY)")
		return CountryPopulation{}, errors.New("invalid end range. Limit (" + strconv.Itoa(startYear) + "-2018)")
	}

	// Filter population from start- and endYear
	var filteredPopulation CountryPopulation
	for _, count := range population.Data.PopulationCounts {
		if count.Year >= startYear && count.Year <= endYear {
			filteredPopulation.Data.PopulationCounts = append(filteredPopulation.Data.PopulationCounts, count)
		}
	}

	return filteredPopulation, nil
}

// getCountry
// Get the country name from countrycode /*
func getCountry(countryCode string) (string, error) {

	// Api request url
	reqUrl := urlRestCountries + fmt.Sprintf("alpha/%s?fields=name", countryCode)

	// Fetch country data
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println("error fetching country data:", err)
		return "", errors.New("failed to fetch country data")
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
		return "", errors.New("failed to read response body")
	}

	// Parse json countries data
	var countryInfo CountryInfo
	if err := json.Unmarshal(respBody, &countryInfo); err != nil {
		log.Println("error parsing JSON:", err)
		return "", errors.New("failed to parse response body")
	}

	return countryInfo.Name.Common, nil
}
