package gatecrasher

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/config"
)

type EventLog struct {
	URL        string `json:"url"`
	StatusCode int    `json:"statusCode"`
}

type Logger interface {
	Printf(format string, v ...interface{})
	SetFlags(flag int)
}

func Run(config config.Config, logger Logger) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if config.Total_number_of_requests <= 0 {
		for {
			makeRequest(config.Target, client, logger)
			time.Sleep(time.Duration(config.Poll_interval_in_seconds) * time.Second)
		}
	}
	for i := 0; i < config.Total_number_of_requests; i++ {
		makeRequest(config.Target, client, logger)
		// This is to avoid sleeping a second for every test case
		if i+1 < config.Total_number_of_requests {
			time.Sleep(time.Duration(config.Poll_interval_in_seconds) * time.Second)
		}
	}
}
func makeRequest(url string, client *http.Client, logger Logger) {
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}

	event := EventLog{
		URL:        url,
		StatusCode: resp.StatusCode,
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	logger.Printf("%s", jsonEvent)
}
