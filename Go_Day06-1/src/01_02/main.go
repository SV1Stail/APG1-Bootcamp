package main

import (
	"01/server"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(100, 100)
var mu sync.Mutex

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if !limiter.Allow() {
			log.Println("Too many requests")
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	creds, err := server.Readadmin_credentials() // Предположим, что ReadAdminCredentials тоже в server
	if err != nil {
		log.Fatalf("Failed to read credentials: %v", err)
	}

	if err := server.ConnectToDB(&creds); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := server.DB.DB()

	if err == nil {
		sqlDB.SetMaxOpenConns(25)           // Максимальное количество открытых соединений
		sqlDB.SetMaxIdleConns(25)           // Максимальное количество простаивающих соединений
		sqlDB.SetConnMaxLifetime(time.Hour) // Время жизни соединений
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./images"))
	mux.Handle("/images/", http.StripPrefix("/images/", fs))
	// Главная страница с отображением статей
	mux.HandleFunc("/", server.MainPage)

	// Страница отдельной статьи с рендерингом markdown
	mux.HandleFunc("/articles/", server.ArticlesPageHandler)
	mux.HandleFunc("/article/", server.ArticlePageHandler)

	// Панель администратора для создания новых статей
	mux.HandleFunc("/admin", server.Login)
	mux.Handle("/admin/insert", server.RequireAuth(http.HandlerFunc(server.WritePost)))
	mux.HandleFunc("/admin/wrong_auth", server.WrongAuthPage)

	handler := rateLimitMiddleware(mux)

	log.Println("Сервер запущен на http://localhost:8888")
	http.ListenAndServe(":8888", handler)
}
