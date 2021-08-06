package models

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

type SearchResponseModel struct {
	Fields        []string               `json:"fields"`
	Facets        map[string][]TermFacet `json:"facetResult"`
	Total         uint64                 `json:"total"`
	Took          time.Duration          `json:"took"`
	Status        bleve.SearchStatus     `json:"status"`
	SearchResults []SearchResult         `json:"searchResults"`
}

type SearchRouteResponseModel struct {
	Fields         []string               `json:"fields"`
	Facets         map[string][]TermFacet `json:"facetResult"`
	Total          uint64                 `json:"total"`
	Took           time.Duration          `json:"took"`
	Status         bleve.SearchStatus     `json:"status"`
	OperatorRoutes []OperatorRoute        `json:"searchResults"`
}

type TermFacet struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}
