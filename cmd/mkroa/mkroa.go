package main

import (
	"flag"
	"fmt"
	"github.com/ags131/go-dn42/dn42"
	"log"
	"os"
)

var repo *string

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	v2 := flag.Bool("v2", false, "Use Bird 2 Format")
	repo = (flag.String("repo", dn42.RepoURL, "Repo URL or path"))
	flag.Parse()
	fmt.Println("v2:", *v2)
	fmt.Println("repo:", *repo)
	dn42.RepoURL = *repo
	routes4, err := dn42.GetRoutes("data/filter.txt", "data/route")
	check(err)
	routes6, err := dn42.GetRoutes("data/filter6.txt", "data/route6")
	check(err)
	if *v2 {
		routes := append(routes4, routes6...)
		WriteRoutes(routes, "bird2_roa.conf", true)
	} else {
		WriteRoutes(routes4, "bird_roa.conf", false)
		WriteRoutes(routes6, "bird6_roa.conf", false)
	}
}

func WriteRoutes(routes []dn42.Route, outPath string, v2 bool) {
	f, err := os.Create(outPath)
	check(err)
	defer f.Close()
	var pref string = "roa"
	if v2 {
		pref = "route"
	}
	for _, route := range routes {
		fmt.Fprintf(f, "%s %s max %d as %s\n", pref, route.Prefix, route.MaxLength, route.Asn[2:])
	}
}
