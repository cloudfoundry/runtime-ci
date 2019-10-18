package main

import (
	"fmt"
	"log"
	"os"

	"stemcell-version-bump/cmd/check/runner"
	"stemcell-version-bump/resource"
)

func main() {
	config, err := resource.NewConfig(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to load config from stdin: %s", err)
	}

	err = resource.SetupEnvironment(config.Source.JSONKey)
	if err != nil {
		log.Fatalf("Failed to setup environment: %s", err)
	}

	client, err := resource.NewGCSClient()
	if err != nil {
		log.Fatalf("Failed to create GCS client: %s", err)
	}

	versions, err := runner.Check(config, client)
	if err != nil {
		log.Fatalf("Failed checking for new versions: %s", err)
	}

	fmt.Println(versions)
}
