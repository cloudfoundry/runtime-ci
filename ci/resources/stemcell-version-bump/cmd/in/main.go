package main

import (
	"fmt"
	"log"
	"os"

	"stemcell-version-bump/cmd/in/runner"
	"stemcell-version-bump/resource"
)

func main() {
	config, err := resource.NewConfig(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to load config from stdin: %s", err)
	}

	client, err := resource.NewGCSClient(config.Source.JSONKey)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %s", err)
	}

	currentVersion, err := runner.In(config, client)
	if err != nil {
		log.Fatalf("Failed to fetch resource: %s", err)
	}

	fmt.Println(currentVersion)
}
