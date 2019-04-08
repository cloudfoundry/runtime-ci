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
const MIN_NUM_SUCCESSFUL_SEQUENTIAL_RESPONSES = 20
const MIN_PROPAGATION_DELAY = 240

func main() {
	client := createHttpsClient()
	url := os.Args[1]

	var (
		firstSuccessTime    time.Time
		numSuccessResponses int
	)

	for startTime := time.Now(); time.Since(startTime).Seconds() < MAX_TIMEOUT_IN_SEC; time.Sleep(5 * time.Second) {
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode != 200 {
			numSuccessResponses = 0
			fmt.Printf("Received error from the API, resetting...\n")
			continue
		}

		// Record time of the first success
		if numSuccessResponses == 0 {
			firstSuccessTime = time.Now()
		}

		numSuccessResponses += 1
		fmt.Printf(
			"Received %d successful responses from the API. %fs remain until propagation delay threshold is reached.\n",
			numSuccessResponses,
			MIN_PROPAGATION_DELAY-time.Since(firstSuccessTime).Seconds(),
		)

		// Api is healthy if both conditions are met:
		// 1. There was at least MIN_NUM_SUCCESSFUL_SEQUENTIAL_RESPONSES from the API server
		// 2. We've waited at least MIN_PROPAGATION_DELAY seconds
		if numSuccessResponses >= MIN_NUM_SUCCESSFUL_SEQUENTIAL_RESPONSES && time.Since(firstSuccessTime).Seconds() >= MIN_PROPAGATION_DELAY {
			fmt.Println("API is healthy!")
			os.Exit(0)
		}
	}

	fmt.Printf("API is unhealthy: could not get %d successful API responses in the row with %ds timeout\n", MIN_NUM_SUCCESSFUL_SEQUENTIAL_RESPONSES, MAX_TIMEOUT_IN_SEC)
	os.Exit(1)
}

func createHttpsClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{Transport: tr, Timeout: time.Duration(30 * time.Second)}
}
