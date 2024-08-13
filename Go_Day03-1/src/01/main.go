package main

import (
	"01/db"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func clientExist(client *elasticsearch.Client, indexName string) (bool, int, error) {
	flag := false
	exists, err := client.Indices.Exists([]string{indexName})
	if exists.StatusCode == 200 {
		fmt.Printf("%s exists", indexName)
		flag = true
	} else if exists.StatusCode == 404 {
		fmt.Printf("%s doesn't exists", indexName)
		flag = false
	} else {
		fmt.Printf("Received unexpected status code: %d\n", exists.StatusCode)
		flag = false
	}
	return flag, exists.StatusCode, err
}

func main() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}
	indexName := "places"

	// clientExist(client, indexName)

	store := db.NewElasticsearchStore(client, indexName)
	var places []db.Place
	var total int
	places, total, err = store.GetPlaces(1000, 0)
	if err != nil {
		log.Fatalf("Error GetPlaces: %v", err)
	}
	log.Printf("Total places: %d\n", total)
	for _, place := range places {
		log.Printf("Place: %+v\n", place)
	}

	// store.GetPlaces(2, 0)
}
