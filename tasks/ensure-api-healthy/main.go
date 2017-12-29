package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Constants
const MAX_TIMEOUT_IN_SEC = 900
const MIN_NUM_SUCCESSFUL_SEQUENTIAL_TRIES = 10

func main() {
	client := createHttpsClient()
	v2InfoURL := fmt.Sprintf("https://api.%s/v2/info", os.Getenv("SYSTEM_DOMAIN"))

	startTime := time.Now()
	numSuccessResponses := 0

	for time.Since(startTime).Seconds() < MAX_TIMEOUT_IN_SEC {
		resp, err := client.Get(v2InfoURL)
		if err != nil || resp.StatusCode != 200 {
			numSuccessResponses = 0
			continue
		}
		numSuccessResponses += 1
		if numSuccessResponses == MIN_NUM_SUCCESSFUL_SEQUENTIAL_TRIES {
			fmt.Println("API is healthy!")
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("API is unhealthy: could not get %d successful API responses in the row with %ds timeout\n", MIN_NUM_SUCCESSFUL_SEQUENTIAL_TRIES, MAX_TIMEOUT_IN_SEC)
	os.Exit(1)
}

func createHttpsClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{Transport: tr}
}
