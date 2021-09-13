package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orbitcacheservices/logger"
	"github.com/orbitcacheservices/models"
	"github.com/orbitcacheservices/search"
	"github.com/orbitcacheservices/service"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	// router.Use(verifyAccessToken)
	router.HandleFunc("/orbitcacheservices", welcomePage)
	router.HandleFunc("/orbitcacheservices/operator/routes", getOrbitOperatorRoutes).Methods("GET")
	router.HandleFunc("/orbitcacheservices/search/{fromCode}/{toCode}/{tripDate}", getBitsSearchResult).Methods("GET")
	router.HandleFunc("/orbitcacheservices/{operatorCode}/{username}/{apiToken}/busmap/{tripCode}/{fromCode}/{toCode}/{travelDate}", getBitsBusMap).Methods("GET")

	http.ListenAndServe(":8080", router)
}

func welcomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func getOrbitOperatorRoutes(w http.ResponseWriter, r *http.Request) {
	var searchModel models.SearchRequestModel
	json.NewDecoder(r.Body).Decode(&searchModel)

	resp, error := search.GetOrbitOperatorRoutes(searchModel)
	if error != nil {
		logger.ErrorLogger.Println(error.Error())
		json.NewEncoder(w).Encode(models.Failure(error.Status(), error.Message()))
	}
	json.NewEncoder(w).Encode(models.Success(resp))
}

func getBitsSearchResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fromStationCode := vars["fromCode"]
	toStationCode := vars["toCode"]
	tripDate := vars["tripDate"]
	searchResults := service.GetBitsSearchResult(fromStationCode, toStationCode, tripDate)
	json.NewEncoder(w).Encode(models.Success(searchResults))
}

func getBitsBusMap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripCode := vars["tripCode"]
	fromStationCode := vars["fromCode"]
	toStationCode := vars["toCode"]
	travelDate := vars["travelDate"]

	var operator models.Operator
	operator.Code = vars["operatorCode"]
	operator.Username = vars["username"]
	operator.ApiToken = vars["apiToken"]

	resp := service.GetBusmap(tripCode, fromStationCode, toStationCode, travelDate, operator)

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
