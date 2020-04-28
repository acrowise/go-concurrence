package main

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)

	// [200 OK] {
	//   "name" : "node-1",
	//   "cluster_name" : "go-elasticsearch"
	// ...
}
