package main

import (
	"fmt"
	"log"
	"os"

	"stemcell-version-bump/cmd/in/runner"
	"stemcell-version-bump/resource"
)

func main() {
	request, err := resource.NewCheckInRequest(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to load request from stdin: %s", err)
	}

	client, err := resource.NewGCSClient(request.Source.JSONKey)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %s", err)
	}

	currentVersion, err := runner.In(request, client)
	if err != nil {
		log.Fatalf("Failed to fetch resource: %s", err)
	}

	fmt.Println(currentVersion)
}
