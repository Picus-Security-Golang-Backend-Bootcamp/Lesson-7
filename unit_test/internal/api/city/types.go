package city

type CreateCityRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"countryCode"`
}

type CityResponse struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"countryCode"`
}
