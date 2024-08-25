package main

import (
	"00/server"
	"net/http"
)

func main() {
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/buy_candy", server.BuyCandy)
	http.ListenAndServe(":3333", rootMux)
}
