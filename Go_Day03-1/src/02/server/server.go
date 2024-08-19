package server

import (
	"02/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

func HandlePlaces(w http.ResponseWriter, r *http.Request) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}
	indexName := "places"
	store := db.NewElasticsearchStore(client, indexName)

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		jsonError(w, fmt.Sprintf(`{"error": "Invalid 'page' value (value is not number): '%s'"}`, pageStr), http.StatusBadRequest)
		return
	}
	if page < 1 {
		jsonError(w, fmt.Sprintf(`{"error": "Invalid 'page' value: '%s'"}`, pageStr), http.StatusBadRequest)
		return
	}

	limit := 10
	offset := (page - 1) * limit
	places, total, err := store.GetPlaces(limit, offset)
	if err != nil {
		jsonError(w, `{"error": "Failed to retrieve places"}`, http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit
	if page > totalPages {
		jsonError(w, fmt.Sprintf(`{"error": "Invalid 'page' value: '%s' too much"}`, pageStr), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"name":       indexName,
		"Rows in db": total,
		"Page":       page,
		"TotalPages": totalPages,
		"Places":     places,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		jsonError(w, `{"error": "Failed to encode response"}`, http.StatusBadRequest)
	}
}

func jsonError(w http.ResponseWriter, str string, errCode int) {

	http.Error(w, str, errCode)
}
