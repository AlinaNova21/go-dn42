package main

import (
	"encoding/json"
	"github.com/ags131/go-dn42/dn42"
	"io/ioutil"
	"log"
)

func main() {
	filters4, err := dn42.ParseFilter("data/filter.txt")
	if err != nil {
		log.Fatal(err)
	}
	filters6, err := dn42.ParseFilter("data/filter6.txt")
	if err != nil {
		log.Fatal(err)
	}

	route4, err := dn42.ParseRoutes("data/route", filters4)
	if err != nil {
		log.Fatal(err)
	}
	route6, err := dn42.ParseRoutes("data/route6", filters6)
	if err != nil {
		log.Fatal(err)
	}

	var rpki struct {
		Roas []dn42.Route `json:"roas"`
	}

	rpki.Roas = append(route4, route6...)

	b, _ := json.Marshal(rpki)
	err = ioutil.WriteFile("export.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
