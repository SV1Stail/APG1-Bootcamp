package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	kind := flag.String("k", "", "accepts two-letter abbreviation for the candy type")
	count := flag.Int("c", 0, "count of candy to buy")

	money := flag.Int("m", 0, "amount of money you \"gave to machine\"")
	flag.Parse()
	if *kind == "" {
		log.Fatal("no candy type")
		os.Exit(1)
	}

	caCert, err := os.ReadFile("certs/minica.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair("certs/clientcandy.tld/cert.pem", "certs/clientcandy.tld/key.pem")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("https://servercandy.tld:3333/buy_candy?k=%s&c=%d&m=%d", *kind, *count, *money)
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

}
