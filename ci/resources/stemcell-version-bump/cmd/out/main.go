package main

import (
	"log"
	"os"

	"stemcell-version-bump/cmd/out/runner"
	"stemcell-version-bump/resource"
)

func main() {
	config, err := resource.NewOutRequest(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to load config from stdin: %s", err)
	}

	client, err := resource.NewGCSClient(config.Source.JSONKey)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %s", err)
	}

	output, err := runner.ReadVersionBump(config)
	if err != nil {
		log.Fatalf("Failed to read output file: %s", err)
	}

	err = runner.Out(config, client, output)
	if err != nil {
		log.Fatalf("Failed to upload version info to GCS: %s", err)
	}
}
