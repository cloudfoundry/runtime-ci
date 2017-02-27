package main

import (
	"log"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher"
)

func main() {
	logger := new(log.Logger)

	gatecrasher.Run("https://api.hermione.cf-app.com/v2/info", logger)
}
