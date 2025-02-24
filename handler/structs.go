package handler

type StatusResponse struct {
	CountriesNowAPI  int     `json:"countriesnowapi"`
	RestCountriesAPI int     `json:"restcountriesapi"`
	Version          string  `json:"version"`
	Uptime           float64 `json:"uptime"`
}
