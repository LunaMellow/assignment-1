package handler

type CitiesResponse struct {
	Data []string `json:"data"`
}

// CountryInfo /info endpoint struct
type CountryInfo struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    []string          `json:"capital"`
	Cities     []string          `json:"cities"`
}

// CountryPopulation /population endpoint struct
type CountryPopulation struct {
	Mean   int `json:"mean"`
	Values []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	} `json:"values"`
}

// StatusResponse /status endpoint struct
type StatusResponse struct {
	CountriesNowAPI  int    `json:"countriesnowapi"`
	RestCountriesAPI int    `json:"restcountriesapi"`
	Version          string `json:"version"`
	Uptime           uint   `json:"uptime"`
}
