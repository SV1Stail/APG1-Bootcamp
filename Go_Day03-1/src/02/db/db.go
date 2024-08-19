package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
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
	var mu sync.Mutex

	go func() {
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			chanHit <- hit
		}
		close(chanHit)
	}()
	go func() {
		for place := range chanPlaces {
			mu.Lock()
			places = append(places, place)
			mu.Unlock()
		}
	}()

	for i := 0; i < kolvoGorutine; i++ {
		wg.Add(1)
		// fmt.Println(3)
		go MarshalUnmarshal(&wg, chanHit, chanPlaces)
	}

	wg.Wait()
	close(chanPlaces)

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
func getGoroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := strings.TrimSpace(string(buf[:n]))
	fields := strings.Split(stack, " ")
	if len(fields) > 1 {
		return fields[1]
	}
	return "unknown"
}
