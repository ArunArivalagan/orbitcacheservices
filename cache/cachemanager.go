package cache

import (
	"fmt"
	"time"

	"github.com/orbitcacheservices/errors"
	"github.com/orbitcacheservices/logger"
	"github.com/orbitcacheservices/models"
	"github.com/orbitcacheservices/search"
	"github.com/orbitcacheservices/utils/date_utils"
)

func TripCreateOrUpdate(t models.TripBleve, modifidDate string, force bool) errors.RestErrors {
	//fmt.Println("productCreateOrUpdate:SKU", occM.Sku, occM.ID)
	key := fmt.Sprintf("%s_%s_%s", t.FromStationCode, t.ToStationCode, date_utils.GetFileDateLayoutFormat(t.TripDate))
	//fmt.Println("productCreateOrUpdate|key:", key)
	t.TripKey = key

	// isModified TODO
	// isModified := IsModified(t, key, modifidDate)
	isModified := 1

	//fmt.Println("CreateOrUpdateProductIntService status", key, isModified)
	if isModified == 1 || force {
		//fmt.Println("force to add bolt db", t.TripKey, string(occMbytes))
		//index only specific attributes for products
		t.ModifiedDate = date_utils.GeApiDBLayoutDateFormat(time.Now().UTC())

		search.SearchIndex.Index(string(key), t.SearchResults)
	}
	return nil
}

func OperatorRoutesCreateOrUpdate(t models.OperatorRouteBleve, modifidDate string, force bool) errors.RestErrors {
	key := fmt.Sprintf("%s_%s", t.OperatorCode, "Routes")

	// isModified TODO
	// isModified := IsModified(t, key, modifidDate)
	isModified := 1

	if isModified == 1 || force {
		t.ModifiedDate = date_utils.GeApiDBLayoutDateFormat(time.Now().UTC())

		search.OpRoutesIndex.Index(string(key), t)
	}
	return nil
}

func IsModified(t models.Trip, key string, modifidDate string) int {
	var err error
	isModified := 0
	srm := models.SearchRequestModel{PharseQueries: []string{fmt.Sprintf(`"%s"`, key)}, Fields: []string{"modifiedDate"}, IndexName: "TripSearchResult", From: 0, Size: 1}

	resM, err := search.GetSearchResult(srm)
	if err != nil {
		logger.ErrorLogger.Println("Failed while fetch indexed value|GetSearchResult ", err)
		return 0
	}
	if resM != nil {
		row := resM.SearchResults[0]
		lastmd, err := time.Parse(date_utils.ApiDbLayout, getInterfaceValToStr(row.ModifiedDate))
		if err != nil {
			logger.ErrorLogger.Println("Failed while parse time|Last ModifiedDate", err)
		}

		md, err := time.Parse(date_utils.ApiDbLayout, modifidDate)
		if err != nil {
			logger.ErrorLogger.Println("Failed while parse time|ModifiedDate", err)
		}

		if md.After(lastmd) {
			isModified = 1
		}
	} else {
		isModified = 1
	}
	return isModified
}

func getInterfaceValToStr(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
