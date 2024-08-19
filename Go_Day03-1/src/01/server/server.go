package server

import (
	"01/db"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<html><h1>400 Bad Request</h1><p>Invalid 'page' value</p></html>"))
		return
	}
	if page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		response := fmt.Sprintf("<html><h1>400 Bad Request</h1><p>Invalid 'page' value: %d</p></html>", page)
		w.Write([]byte(response))
		return
	}
	var places []db.Place
	limit := 10
	offset := (page - 1) * limit // number of skipped lines

	places, total, err := store.GetPlaces(limit, offset)

	if err != nil || page > total/limit {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<html><h1>400 Bad Request</h1><p>Invalid 'page' value</p></html>")
		return
	}

	totalPages := (total + limit - 1) / limit

	data := struct {
		Places     []db.Place
		Total      int
		Page       int
		TotalPages int
	}{
		Places:     places,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}
	tmplPath := filepath.Join("server", "template.html")
	t, err := template.New(filepath.Base(tmplPath)).Funcs(template.FuncMap{
		"minus": func(a, b int) int { return a - b },
		"plus":  func(a, b int) int { return a + b },
	}).ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
