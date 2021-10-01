package config

import (
	"log"

	"github.com/JavakarBits/orbitcacheservices/search"
)

func StartApplication() {
	log.Println("Application Starting...")
	search.CreateIndex()
}
