package search

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/orbitcacheservices/errors"
	"github.com/orbitcacheservices/models"
	"github.com/orbitcacheservices/utils"
)

var (
	SearchIndex bleve.Index
)

func getIndex(indexName string) bleve.Index {
	var index bleve.Index
	switch indexName {
	case "search":
		index = SearchIndex
	}
	return index
}

func CreateIndex() {
	var err error
	SearchIndex, err = bleve.Open("data/search.bleve")
	if err == bleve.ErrorIndexPathDoesNotExist {
		fmt.Println(fmt.Sprintf("Creating product new index %s ... ", "data/search.bleve"))
		// create a mapping
		indexMapping, err := buildEventIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		SearchIndex, err = bleve.New("data/search.bleve", indexMapping)
		if err != nil {
			fmt.Println("Failed product index field mapping", err)
		}
	}
}

func GetSearchResult(m models.SearchRequestModel) (*models.SearchResponseModel, errors.RestErrors) {
	q := bleve.NewBooleanQuery()

	for _, v := range m.Terms {
		mq := bleve.NewMatchQuery(v)
		q.AddMust(mq)
	}
	if len(m.PharseQueries) > 0 {
		for _, v := range m.PharseQueries {
			mpq := bleve.NewMatchPhraseQuery(v)
			q.AddMust(mpq)
		}
	}

	req := bleve.NewSearchRequest(q)

	if len(m.Fields) > 0 {
		req.Fields = m.Fields
	} else {
		req.Fields = []string{"*"}
	}
	if len(m.SortBy) > 0 {
		req.SortBy(m.SortBy)
	}
	for _, f := range m.Facets {
		req.AddFacet(f, bleve.NewFacetRequest(f, 10))
	}
	req.From = m.From
	req.Size = m.Size
	if m.Size == 0 {
		req.Size = 50
	}
	CreateIndex()
	index := getIndex(m.IndexName)
	fmt.Println("index name", m.IndexName, index.Name())
	res, err := index.Search(req)
	if err != nil {
		fmt.Println("Failed while execute the query", err)
		return nil, errors.NewInternalServerError("Failed while execute the query", err)
	}
	rm := models.SearchResponseModel{}
	resultRow := make([]map[string]interface{}, 0)

	rm.Fields = res.Request.Fields
	// resS, _ := json.Marshal(res)
	// fmt.Println("search res", string(resS))
	rm.Total = res.Total
	rm.Took = res.Took
	rm.Status = *res.Status
	for _, rv := range res.Hits {
		resultRow = append(resultRow, rv.Fields)
	}

	var searchResults []models.SearchResult
	for _, rr := range resultRow {
		tripCodes := rr["tripCode"].([]interface{})
		tripDate := rr["tripDate"].([]interface{})
		modifiedDate := rr["modifiedDate"].([]interface{})
		availableSeatCount := rr["availableSeatCount"].([]interface{})
		journeyMinutes := rr["journeyMinutes"].([]interface{})
		fares := getInterfaceArray(rr["fares"].([]interface{}))

		amenities := rr["amenities"].([]interface{})
		operator := rr["operator"].([]interface{})
		bus := rr["bus"].([]interface{})
		tax := rr["tax"].([]interface{})
		tripStatus := rr["tripStatus"].([]interface{})
		fromStation := rr["fromStation"].([]interface{})
		toStation := rr["toStation"].([]interface{})

		tripCount := len(tripCodes)
		for i := 0; i < tripCount; i++ {
			var searchResult models.SearchResult

			availableSeat, _ := strconv.ParseInt(fmt.Sprint(availableSeatCount[i]), 10, 64)
			searchResult.AvailableSeatCount = int(availableSeat)

			jm, _ := strconv.ParseInt(fmt.Sprint(journeyMinutes[i]), 10, 64)
			searchResult.JourneyMinutes = int(jm)

			amen, _ := utils.UnMarshalBinaryArrayBase([]byte(fmt.Sprint(amenities[i])))
			searchResult.Amenities = amen

			op, _ := utils.UnMarshalBinaryBase([]byte(fmt.Sprint(operator[i])))
			searchResult.Operator = op

			bs, _ := utils.UnMarshalBinaryBus([]byte(fmt.Sprint(bus[i])))
			searchResult.Bus = bs

			tx, _ := utils.UnMarshalBinaryTax([]byte(fmt.Sprint(tax[i])))
			searchResult.Tax = tx

			tps, _ := utils.UnMarshalBinaryBase([]byte(fmt.Sprint(tripStatus[i])))
			searchResult.TripStatus = tps

			fr, _ := utils.UnMarshalBinaryFares([]byte(fmt.Sprint(fares[i])))
			searchResult.Fares = fr

			searchResult.ModifiedDate = fmt.Sprint(modifiedDate[i])
			searchResult.TripCode = fmt.Sprint(tripCodes[i])
			searchResult.TripDate = fmt.Sprint(tripDate[i])

			fms, _ := utils.UnMarshalBinaryStation([]byte(fmt.Sprint(fromStation[i])))
			searchResult.FromStation = fms

			tms, _ := utils.UnMarshalBinaryStation([]byte(fmt.Sprint(toStation[i])))
			searchResult.ToStation = tms

			searchResults = append(searchResults, searchResult)
		}
	}
	rm.SearchResults = searchResults

	if len(res.Facets) > 0 {
		rm.Facets = make(map[string][]models.TermFacet)
	}
	for k, v := range res.Facets {
		terms := make([]models.TermFacet, 0)
		for _, fv := range v.Terms {
			terms = append(terms, models.TermFacet{Term: fv.Term, Count: fv.Count})
		}
		rm.Facets[k] = terms
	}
	//resbyte, err := json.Marshal(rm)
	if err != nil {
		fmt.Println("Failed while marshal final search result", err)
		return nil, errors.NewRestErrors("Failed while final search result", http.StatusInternalServerError, err.Error(), nil)
	}
	//resStr := string(resbyte)
	//fmt.Println("Final search", resStr)
	//header
	if len(resultRow) == 0 {
		return nil, nil
	}
	fieldPos := make(map[string]int)
	if rm.Fields[0] != "*" {
		fmt.Println("read from fields list")
		for idx, fn := range rm.Fields {
			fieldPos[fn] = idx
		}
	} else {
		idx := 0
		for f := range resultRow[0] {
			fieldPos[f] = idx
			idx = idx + 1
		}
	}

	// generateCSV(&m, fieldPos, resultRow)
	return &rm, nil
}

func getInterfaceArray(data []interface{}) []interface{} {
	if data == nil {
		var data1 []interface{}
		return data1
	}
	return data
}

func generateCSV(reqM *models.SearchRequestModel, fields map[string]int, resultRow []map[string]interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed while get the working dir", err)
		return
	}
	dt := time.Now().UTC()
	p := fmt.Sprintf("%s/csvs/%s/%d%d%d", wd, reqM.IndexName, dt.Year(), dt.Month(), dt.Day())

	if _, err := os.Stat(p); os.IsNotExist(err) {
		//fmt.Println("generateCSV", err)
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			fmt.Println("Failed while create file", err)
			return
		}
	}
	fn := fmt.Sprintf("search_%d_%d_%d%d%d%d.csv", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
	//fullFileName:=fmt.Sprintf("%s/%s", path, fn)
	file, err := os.Create(path.Join(p, fn))
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed creating file %s", path.Join(p, fn)), err)
	}
	defer file.Close()
	if err != nil {
		fmt.Println("Failed while create csv files", err)
	}
	// 2. Initialize the writer
	writer := csv.NewWriter(file)

	csvData := make([][]string, 0)

	l := len(fields)
	row := make([]string, l)
	for k, v := range fields {
		//	fmt.Println("field order", l, k, v)
		row[v] = k
	}

	//fmt.Println("print header", row)
	csvData = append(csvData, row)
	//product data
	for _, r := range resultRow {
		row = make([]string, len(fields))
		for f, val := range r {
			idx, isFound := fields[f]
			if isFound {
				//fmt.Println("col value", f, val, idx)
				row[idx] = fmt.Sprintf("%v", val)
			}
		}
		//fmt.Println("row len", len(row))
		csvData = append(csvData, row)
	}
	// 3. Write all the records
	err = writer.WriteAll(csvData) // returns error
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
}
