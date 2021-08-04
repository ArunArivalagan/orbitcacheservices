package models

type Bus struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	BusType      string `json:"busType"`
	CategoryCode string `json:"categoryCode"`
}
