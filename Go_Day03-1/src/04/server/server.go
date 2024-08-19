package server

import (
	"04/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Token(w http.ResponseWriter, r *http.Request) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claim := &Claims{
		Username: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // создание ключа
	tokenStr, err := token.SignedString(jwtKey)               // подпись ключа с помощью jwtKey

	// tokenStr2 := token
	// tokenStr3 := tokenStr
	// fmt.Printf("2 %v\n", tokenStr2)
	// fmt.Printf("3 %v\n", tokenStr3)
	if err != nil {
		jsonError(w, "cant sign string", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + tokenStr + `"}`))
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			jsonError(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})

}

func Recommend(w http.ResponseWriter, r *http.Request) {
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
