package gatecrasher

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type EventLog struct {
	URL        string `json:"url"`
	StatusCode int    `json:"statusCode"`
}

type Logger interface {
	Printf(format string, v ...interface{})
	SetFlags(flag int)
	SetOutput(writer io.Writer)
}

func Run(url string, logger Logger) int {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}
	// Ensure our logging contains a timestamp
	logger.SetOutput(os.Stdout)
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
