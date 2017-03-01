package gatecrasher

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/config"
)

type EventLog struct {
	StatusCode int    `json:"statusCode"`
	Type       string `json:"type"`
	URL        string `json:"url"`
}

type SummaryEventLog struct {
	FinishTime     time.Time `json:"finishTime"`
	IntervalSize   int       `json:"intervalSize"`
	PercentSuccess float64   `json:"percentSuccess"`
	StartTime      time.Time `json:"startTime"`
	Type           string    `json:"type"`
	URL            string    `json:"url"`
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

	startTime := time.Now()
	successCounter := 0
	for i := 1; ; i++ {
		if i > config.TotalNumberOfRequests {
			return
		}
		event := makeRequest(config.Target, client, logger)

		if event.StatusCode == http.StatusOK {
			successCounter++
		}

		if config.SkipIndividualRequests != true {
			logJson(event, logger)
		}

		time.Sleep(time.Duration(config.PollIntervalInMs) * time.Millisecond)
		if i%config.ReportIntervalInRequests == 0 {
			summaryEventLog := SummaryEventLog{
				FinishTime:     time.Now(),
				IntervalSize:   config.ReportIntervalInRequests,
				PercentSuccess: float64(successCounter) / float64(config.ReportIntervalInRequests),
				StartTime:      startTime,
				Type:           "summary",
				URL:            config.Target,
			}
			logJson(summaryEventLog, logger)
			startTime = time.Now()
			successCounter = 0
		}
	}
}

func makeRequest(url string, client *http.Client, logger Logger) EventLog {
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}

	event := EventLog{
		StatusCode: resp.StatusCode,
		Type:       "request",
		URL:        url,
	}
	return event
}

func logJson(jsonObject interface{}, logger Logger) {
	marshalledJson, err := json.Marshal(jsonObject)
	if err != nil {
		panic(err)
	}
	logger.Printf("%s", marshalledJson)
}
