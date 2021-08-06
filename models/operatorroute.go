package models

type OperatorRoute struct {
	Operator Base    `json:"operator"`
	Routes   []Route `json:"routes"`
}

type OperatorRouteBleve struct {
	OperatorCode string `json:"operatorCode"`
	ModifiedDate string `json:"modifiedDate"`
	Operator     string `json:"operator"`
	Routes       string `json:"routes"`
}
