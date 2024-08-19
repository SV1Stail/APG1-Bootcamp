package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
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
type Store interface {
	GetPlaces(limit int, offset int) ([]Place, int, error)
}

type EsStore struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchStore(client *elasticsearch.Client, index string) *EsStore {
	return &EsStore{
		client: client,
		index:  index,
	}
}

func (s *EsStore) GetNearesRestaurant(lat, lon float64, limit int) ([]Place, error) {
	ctx := context.Background()

	query := map[string]interface{}{
		"size": limit,
		"sort": map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": lat,
					"lon": lon,
				},
				"order":           "asc",
				"unit":            "km",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("ошибка поиска: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	kolvoGorurines := 2
	hitChan := make(chan interface{}, kolvoGorurines+5)
	placeChan := make(chan Place, kolvoGorurines+5)
	var wg sync.WaitGroup
	var places []Place
	go func() {
		for _, h := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			hitChan <- h
		}
		close(hitChan)
	}()

	for i := 0; i < kolvoGorurines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for hit := range hitChan {
				source := hit.(map[string]interface{})["_source"]
				var place Place
				buf, err := json.Marshal(source)
				if err != nil {
					fmt.Printf("ERROR! problem marshaling source: %v", source)
					continue
				}
				if err := json.Unmarshal(buf, &place); err != nil {
					fmt.Printf("ERROR! problem Unmarshalling source: %v", source)
					continue
				}
				placeChan <- place
			}

		}()
	}
	go func() {
		wg.Wait()
		close(placeChan)
	}()
	for place := range placeChan {
		places = append(places, place)
	}
	return places, err
}

func (s *EsStore) GetPlaces(limit int, offset int) ([]Place, int, error) {
	ctx := context.Background()
	query := map[string]interface{}{
		"from": offset, // Смещение
		"size": limit,  // Лимит
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}
	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("ошибка поиска: %s", res.String())
	}

	// Обработка ответа
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	// fmt.Printf("%v\n", r)
	var wg sync.WaitGroup
	const kolvoGorutine = 10
	chanPlaces := make(chan Place, 15)
	chanHit := make(chan interface{}, 15)
	var places []Place

	go func() {
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			chanHit <- hit
		}
		close(chanHit)
	}()

	for i := 0; i < kolvoGorutine; i++ {
		wg.Add(1)
		// fmt.Println(3)
		go MarshalUnmarshal(&wg, chanHit, chanPlaces)
	}

	go func() {
		wg.Wait()
		close(chanPlaces)
	}()

	for place := range chanPlaces {
		places = append(places, place)
	}
	total := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	return places, total, nil
}

func MarshalUnmarshal(wg *sync.WaitGroup, chanHit <-chan interface{}, chanPlaces chan<- Place) {
	defer wg.Done()

	// goroutineID := getGoroutineID()
	// fmt.Printf("Горутина %s \n", goroutineID)
	for hit := range chanHit {
		// fmt.Printf("Горутина %s \n", goroutineID)

		var place Place
		source := hit.(map[string]interface{})["_source"]

		sourceJSON, err := json.Marshal(source)
		if err != nil {
			fmt.Printf("ERROR! problem marshaling source: %v", source)
			continue
		}
		if err := json.Unmarshal(sourceJSON, &place); err != nil {
			fmt.Printf("ERROR! problem Unmarshal %v", source)
			continue
		}
		chanPlaces <- place
	}
}
