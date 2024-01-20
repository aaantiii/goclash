package clash

type Location struct {
	LocalizedName string `json:"localizedName,omitempty"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	IsCountry     bool   `json:"isCountry"`
	CountryCode   string `json:"countryCode"`
}

type Language struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	LanguageCode string `json:"languageCode"`
}
