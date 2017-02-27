package main

import (
	"log"
	"os"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	gatecrasher.Run("https://api.hermione.cf-app.com/v2/info", logger)
}
