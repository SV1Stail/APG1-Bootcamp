package server

import (
	"03/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || lat < 0 {
		jsonError(w, fmt.Sprintf(`{Error: invald lat %s}`, latStr), http.StatusBadRequest)
		return
	}
	lonStr := r.URL.Query().Get("lon")
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil || lon < 0 {
		jsonError(w, fmt.Sprintf(`{Error: invald lon %s}`, lonStr), http.StatusBadRequest)
		return
	}
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}
	indexName := "places"
	store := db.NewElasticsearchStore(client, indexName)
	var bestPlaces []db.Place
	bestPlaces, err = store.GetNearesRestaurant(lat, lon, 3)
	if err != nil {
		jsonError(w, `{"error": "Failed Get Neares Restaurant"}`, http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

	data := map[string]interface{}{
		"name":   indexName,
		"Places": bestPlaces,
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		jsonError(w, `{"error": "Failed to encode response"}`, http.StatusBadRequest)
	}

}
func jsonError(w http.ResponseWriter, str string, errCode int) {

	http.Error(w, str, errCode)
}
