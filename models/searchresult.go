package models

type SearchResult struct {
	TripCode           string  `json:"tripCode"`
	TripDate           string  `json:"tripDate"`
	ModifiedDate       string  `json:"modifiedDate"`
	AvailableSeatCount int     `json:"availableSeatCount"`
	JourneyMinutes     int     `json:"journeyMinutes"`
	Fares              []int   `json:"fares,omitempty"`
	Bus                Bus     `json:"bus"`
	FromStation        Station `json:"fromStation"`
	ToStation          Station `json:"toStation"`
	Operator           Base    `json:"operator"`
	TripStatus         Base    `json:"tripStatus"`
	Amenities          []Base  `json:"amenities"`
	Tax                Tax     `json:"tax"`
}

type SearchResultBleve struct {
	TripCode           string `json:"tripCode"`
	TripDate           string `json:"tripDate"`
	ModifiedDate       string `json:"modifiedDate"`
	AvailableSeatCount int    `json:"availableSeatCount"`
	JourneyMinutes     int    `json:"journeyMinutes"`
	Fares              string `json:"fares,omitempty"`
	Bus                string `json:"bus"`
	FromStation        string `json:"fromStation"`
	ToStation          string `json:"toStation"`
	Operator           string `json:"operator"`
	TripStatus         string `json:"tripStatus"`
	Amenities          string `json:"amenities"`
	Tax                string `json:"tax"`
}
