package main

import (
	"encoding/json"
	"github.com/ags131/go-dn42/dn42"
	"io/ioutil"
	"log"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	routes4, err := dn42.GetRoutes("data/filter.txt", "data/route")
	check(err)
	routes6, err := dn42.GetRoutes("data/filter6.txt", "data/route6")
	check(err)

	var rpki struct {
		Roas []dn42.Route `json:"roas"`
	}

	rpki.Roas = append(routes4, routes6...)

	b, _ := json.Marshal(rpki)
	err = ioutil.WriteFile("export.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
