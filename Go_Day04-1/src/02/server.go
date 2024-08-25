package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *ask_cow(const char *phrase) {
    int phrase_len = strlen(phrase);
    size_t buf_size = 160 + (phrase_len + 2) * 3;
    char *buf = (char *)malloc(buf_size);

    if (!buf) {
        return NULL; // Handle allocation failure
    }

    int offset = 0;

    offset += snprintf(buf + offset, buf_size - offset, " ");
    offset += snprintf(buf + offset, buf_size - offset, "%.*s", phrase_len + 2, "______________________________________________________________");
    offset += snprintf(buf + offset, buf_size - offset, "\n< %s >\n ", phrase);
    offset += snprintf(buf + offset, buf_size - offset, "%.*s", phrase_len + 2, "--------------------------------------------------------------");
    offset += snprintf(buf + offset, buf_size - offset, "\n        \\   ^__^\n");
    offset += snprintf(buf + offset, buf_size - offset, "         \\  (oo)\\_______\n");
    offset += snprintf(buf + offset, buf_size - offset, "            (__)\\       )\\/\\\n");
    offset += snprintf(buf + offset, buf_size - offset, "                ||----w |\n");
    offset += snprintf(buf + offset, buf_size - offset, "                ||     ||\n");

    return buf;
}

*/
import "C"
import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"unsafe"
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
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Problem with decoding json", http.StatusBadRequest)
		return
	}

	candies := map[string]int{
		"CE": 10,
		"AA": 15,
		"NT": 17,
		"DE": 21,
		"YR": 23,
	}

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
		phrase := C.CString("Thank you!")
		defer C.free(unsafe.Pointer(phrase)) // Освобождение памяти для строки

		cowSay := C.ask_cow(phrase)
		defer C.free(unsafe.Pointer(cowSay))

		respond := Respond{Thanks: C.GoString(cowSay), Change: order.Money - order.CandyCount*price}
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

	serverCerts, err := tls.LoadX509KeyPair("certs/candy.tld/cert.pem", "certs/candy.tld/key.pem")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCerts},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
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
