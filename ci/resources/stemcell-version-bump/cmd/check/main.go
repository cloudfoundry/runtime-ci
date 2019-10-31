package main

import (
	"fmt"
	"log"
	"os"

	"stemcell-version-bump/cmd/check/runner"
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

	versions, err := runner.Check(request, client)
	if err != nil {
		log.Fatalf("Failed checking for new versions: %s", err)
	}

	fmt.Println(versions)
}
