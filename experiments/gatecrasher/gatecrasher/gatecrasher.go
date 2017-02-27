package gatecrasher

import (
	"encoding/json"
	"log"
	"net/http"
)

type EventLog struct {
	URL        string `json:"url"`
	StatusCode int    `json:"statusCode"`
}

type Logger interface {
	Printf(format string, v ...interface{})
	SetFlags(flag int)
}

func Run(url string, logger Logger) int {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}
	// Ensure our logging contains a timestamp
	logger.SetFlags(log.LstdFlags)
	event := EventLog{
		URL:        url,
		StatusCode: resp.StatusCode,
	}
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	logger.Printf("%s", jsonEvent)
	return resp.StatusCode
}
