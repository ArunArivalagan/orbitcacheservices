package models

type Trip struct {
	TripKey         string         `json:"tripkey"`
	ModifiedDate    string         `json:"modifiedDate"`
	TripDate        string         `json:"tripDate"`
	FromStationCode string         `json:"fromStationCode"`
	ToStationCode   string         `json:"toStationCode"`
	SearchResults   []SearchResult `json:"data"`
}

type TripBleve struct {
	TripKey         string              `json:"tripkey"`
	ModifiedDate    string              `json:"modifiedDate"`
	TripDate        string              `json:"tripDate"`
	FromStationCode string              `json:"fromStationCode"`
	ToStationCode   string              `json:"toStationCode"`
	SearchResults   []SearchResultBleve `json:"data"`
}
