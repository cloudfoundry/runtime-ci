package main

import (
	"log"
	"os"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/config"
	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher"
)

func main() {
	config := config.Load()
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// gatecrasher.Run("https://api.luna.cf-app.com/v2/info", logger)
	gatecrasher.Run(config, logger)
}
