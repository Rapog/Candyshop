package main

import (
	"cli/operations"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	candyType := flag.String("k", "", "Candy type")
	candyCount := flag.Int64("c", 0, "Quantity of candyes")
	money := flag.Int64("m", 0, "Money")

	flag.Parse()

	var reqCandyBody = operations.BuyCandyBody{
		CandyCount: candyCount,
		CandyType:  candyType,
		Money:      money,
	}

	var request operations.BuyCandyParams = operations.BuyCandyParams{}
	request.SetDefaults()
	request.SetOrder(reqCandyBody)

	certFile := "localhost/cert.pem"
	keyFile := "localhost/key.pem"
	caFile := "../minica.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate and key: %v", err)
	}
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	rootCAPool := x509.NewCertPool()
	if ok := rootCAPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("Failed to append CA certificate to pool")
	}

	transport := httptransport.New("localhost:3333", "/", []string{"https"})
	transport.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			Certificates:       []tls.Certificate{cert},
			RootCAs:            rootCAPool,
		},
	}

	client := operations.New(transport, strfmt.Default)
	res, err := client.BuyCandy(&request)
	if err != nil {
		if strings.Contains(err.Error(), "Post ") {
			log.Fatalf("Error requesting: %s", err)
		}
		var errReq reqError

		sliceErr := strings.Split(err.Error(), "{")
		var errResp = "{" + sliceErr[len(sliceErr)-1]

		unErr := json.Unmarshal([]byte(errResp), &errReq)
		if unErr != nil {
			log.Printf("Unmarshal err: %s", unErr)
		}
		fmt.Println(errReq.Error)
	} else {
		fmt.Printf("%s Your change is %d", res.Payload.Thanks, res.Payload.Change)
	}
}

type reqError struct {
	Error string `json:"error"`
}
