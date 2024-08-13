package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Place struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location GeoPoint `json:"location"`
}

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func main() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	indexName := "places"
	mapping := `{
	"mappings": {
			"properties": {
	    "name": {
	        "type":  "text"
	    },
	    "address": {
	        "type":  "text"
	    },
	    "phone": {
	        "type":  "text"
	    },
	    "location": {
	      "type": "geo_point"
	    }
	  }
	}
	}`

	_, err = client.Indices.Delete([]string{indexName})
	if err != nil {
		log.Printf("Error deleting index: %s", err)
	}
	esAnswer, err := client.Indices.Create(indexName,
		client.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))),
	)
	if err != nil {
		log.Fatalf("Error creating the index: %s", err)
	}
	defer esAnswer.Body.Close()

	fmt.Printf("Index %s created with mapping.\n", indexName)
	// time.Sleep(1 * time.Second)

	// exists, err := client.Indices.Exists([]string{indexName})
	// if exists.StatusCode == 200 {
	// 	fmt.Printf("Index %s exists.\n", indexName)
	// } else if exists.StatusCode == 404 {
	// 	fmt.Printf("Index %s does not exist.\n", indexName)
	// } else {
	// 	fmt.Printf("Received unexpected status code: %d\n", exists.StatusCode)
	// }

	const kolvoGorutine = 10
	var wg sync.WaitGroup
	placeChan := make(chan Place, 20)
	for i := 0; i < kolvoGorutine; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for place := range placeChan {
				data := esutil.NewJSONReader(place)
				res, err := client.Index(indexName, data, client.Index.WithDocumentID(strconv.Itoa((place.ID))))
				if err != nil {
					log.Printf("Error indexing document ID=%d: %s", place.ID, err)
					continue
				}
				defer res.Body.Close()
				if res.IsError() {
					log.Printf("Error response from Elasticsearch for document ID=%d: %s", place.ID, res.String())
				}
				res.Body.Close()
			}
		}()
	}
	csvFile_path := "../../materials/data.csv"
	csvFile, err := os.Open(csvFile_path)
	if err != nil {
		log.Fatalf("Error opening CSV file: %s", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = '\t'
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Error reading CSV header: %s", err)
	}

	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error reading CSV record: %s", err)
		}
		id, _ := strconv.Atoi(line[0])
		lon, _ := strconv.ParseFloat(line[4], 64)
		lat, _ := strconv.ParseFloat(line[5], 64)

		placeForSend := Place{
			ID:      id,
			Name:    line[1],
			Address: line[2],
			Phone:   line[3],
			Location: GeoPoint{
				Lat: lat,
				Lon: lon,
			},
		}
		placeChan <- placeForSend
	}
	close(placeChan)
	wg.Wait()
	fmt.Println("All records have been processed and indexed successfully.")
}
