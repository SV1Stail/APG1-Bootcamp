package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Order struct {
	Money      int
	CandyType  string
	CandyCount int
}

type Respond struct {
	Thanks string `json:"thanks,omitempty"`
	Change int    `json:"change,omitempty"`
	Error  string `json:"error,omitempty"`
}

func BuyCandy(w http.ResponseWriter, r *http.Request) {
	// if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
	// 	http.Error(w, "Problem with decoding json", http.StatusBadRequest)
	// 	return
	// }
	var order Order
	k := r.URL.Query().Get("k")
	order.CandyType = k
	c := r.URL.Query().Get("c")

	cInt, err := strconv.Atoi(c)
	if err != nil {
		respond := Respond{Error: "Problem with decoding -c flag value"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respond)
		return
	}
	order.CandyCount = cInt

	m := r.URL.Query().Get("m")
	mInt, err := strconv.Atoi(m)
	if err != nil {
		respond := Respond{Error: "Problem with decoding -m flag value"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respond)
		return
	}
	order.Money = mInt
	candies := map[string]int{
		"CE": 10,
		"AA": 15,
		"NT": 17,
		"DE": 21,
		"YR": 23,
	}
	// w.Header().Set("Content-Type", "application/json")
	price, ok := candies[order.CandyType]
	if !ok {
		respond := Respond{Error: "Problem with decoding json"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respond)
		return
	}

	if order.CandyCount < 0 {
		respond := Respond{Error: "invalid candy count"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respond)
		return
	}

	if order.Money < order.CandyCount*price {
		respond := Respond{Error: fmt.Sprintf("You need %d more money", order.CandyCount*price-order.Money)}
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(respond)
		return
	}

	if order.CandyCount*price <= order.Money {
		respond := Respond{Thanks: "Thank you!", Change: order.Money - order.CandyCount*price}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(respond)

	}

}

func main() {
	caCert, err := os.ReadFile("certs/minica.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	serverCerts, err := tls.LoadX509KeyPair("certs/servercandy.tld/cert.pem", "certs/servercandy.tld/key.pem")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCerts},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/buy_candy", BuyCandy)
	server := http.Server{
		Addr:      ":3333",
		Handler:   rootMux,
		TLSConfig: tlsConfig,
	}
	log.Println("Starting server on https://localhost:3333...")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
