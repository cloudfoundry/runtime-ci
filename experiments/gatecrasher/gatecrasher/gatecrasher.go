package gatecrasher

import (
	"crypto/tls"
	"encoding/json"
	"net/http"

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
	for i := 0; i < config.Total_number_of_requests; i++ {
		resp, err := client.Get(config.Target)
		if err != nil {
			panic(err)
		}

		err = resp.Body.Close()
		if err != nil {
			panic(err)
		}

		event := EventLog{
			URL:        config.Target,
			StatusCode: resp.StatusCode,
		}

		jsonEvent, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		logger.Printf("%s", jsonEvent)
	}
}
