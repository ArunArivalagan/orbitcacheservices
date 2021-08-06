package utils

import (
	"encoding/json"

	"github.com/orbitcacheservices/models"
)

func UnMarshalBinaryArraySearchRequest(data []byte) (models.SearchRequestModel, error) {
	var req models.SearchRequestModel
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryArraySearchResult(data []byte) (models.Trip, error) {
	var req models.Trip
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryArrayBase(data []byte) ([]models.Base, error) {
	var req []models.Base
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryBase(data []byte) (models.Base, error) {
	var req models.Base
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryFares(data []byte) ([]int, error) {
	var req []int
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryStationPointList(data []byte) ([]models.StationPoint, error) {
	var req []models.StationPoint
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryTax(data []byte) (models.Tax, error) {
	var req models.Tax
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryBus(data []byte) (models.Bus, error) {
	var req models.Bus
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryStation(data []byte) (models.Station, error) {
	var req models.Station
	err := json.Unmarshal(data, &req)
	return req, err
}

func UnMarshalBinaryOperatorRouteList(data []byte) ([]models.Route, error) {
	var req []models.Route
	err := json.Unmarshal(data, &req)
	return req, err
}
