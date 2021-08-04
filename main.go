package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orbitcacheservices/cache"
	"github.com/orbitcacheservices/logger"
	"github.com/orbitcacheservices/models"
	"github.com/orbitcacheservices/search"
	"github.com/orbitcacheservices/utils/date_utils"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	// router.Use(verifyAccessToken)
	router.HandleFunc("/orbitcacheservices", welcomePage)

	router.HandleFunc("/orbitcacheservices/push/search", pushSearchResult).Methods("POST")
	router.HandleFunc("/orbitcacheservices/search", getSearchResult).Methods("POST")

	http.ListenAndServe(":8080", router)
}

func welcomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func pushSearchResult(w http.ResponseWriter, r *http.Request) {
	var trip models.Trip
	json.NewDecoder(r.Body).Decode(&trip)

	search.CreateIndex()

	var searchResults []models.SearchResultBleve
	for _, sr := range trip.SearchResults {
		var searchResult models.SearchResultBleve

		a, _ := json.Marshal(sr.Amenities)
		searchResult.Amenities = string(a)

		searchResult.AvailableSeatCount = sr.AvailableSeatCount
		searchResult.JourneyMinutes = sr.JourneyMinutes

		op, _ := json.Marshal(sr.Operator)
		searchResult.Operator = string(op)

		tx, _ := json.Marshal(sr.Tax)
		searchResult.Tax = string(tx)

		tps, _ := json.Marshal(sr.TripStatus)
		searchResult.TripStatus = string(tps)

		bs, _ := json.Marshal(sr.Bus)
		searchResult.Bus = string(bs)

		fr, _ := json.Marshal(sr.Fares)
		searchResult.Fares = string(fr)

		searchResult.ModifiedDate = sr.ModifiedDate
		searchResult.TripCode = sr.TripCode
		searchResult.TripDate = sr.TripDate

		fs, _ := json.Marshal(sr.FromStation)
		searchResult.FromStation = string(fs)

		ts, _ := json.Marshal(sr.ToStation)
		searchResult.ToStation = string(ts)

		searchResults = append(searchResults, searchResult)
	}
	var tripBleve models.TripBleve
	tripBleve.FromStationCode = trip.FromStationCode
	tripBleve.ModifiedDate = trip.ModifiedDate
	tripBleve.SearchResults = searchResults
	tripBleve.ToStationCode = trip.ToStationCode
	tripBleve.TripDate = trip.TripDate
	tripBleve.TripKey = trip.TripKey

	error := cache.TripCreateOrUpdate(tripBleve, date_utils.GetNowDBFormat(), true)
	if error != nil {
		logger.ErrorLogger.Println(error.Error())
		json.NewEncoder(w).Encode(models.Failure(error.Status(), error.Message()))
	}

	json.NewEncoder(w).Encode(models.Success(nil))
}

func getSearchResult(w http.ResponseWriter, r *http.Request) {
	var searchModel models.SearchRequestModel
	json.NewDecoder(r.Body).Decode(&searchModel)

	resp, error := search.GetSearchResult(searchModel)
	if error != nil {
		logger.ErrorLogger.Println(error.Error())
		json.NewEncoder(w).Encode(models.Failure(error.Status(), error.Message()))
	}
	json.NewEncoder(w).Encode(models.Success(resp))
}

// func verifyAccessToken(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var header = r.Header.Get("x-access-token")
// 		header = strings.TrimSpace(header)

// 		if header != "hxxjfehp79q69nzp" {
// 			w.WriteHeader(http.StatusForbidden)
// 			json.NewEncoder(w).Encode(entity.Failure(int(consts.Unauthorized), consts.Unauthorized.Error()))
// 			return
// 		}
// 		next.SeindexNameseHTTP(w, r)
// 	})
// }
