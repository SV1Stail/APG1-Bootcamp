package main

import (
	"02/server"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", server.HandlePlaces)
	fmt.Println("Server is listening on port 8888...")
	http.ListenAndServe(":8888", nil)
}
