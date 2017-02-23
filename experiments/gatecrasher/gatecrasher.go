package gatecrasher

import "net/http"

//go:generate counterfeiter . Logger

type Logger interface {
	Printf(format string, v ...interface{})
}

func Run(url string, logger Logger) int {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	logger.Printf("hey")
	return resp.StatusCode
}
