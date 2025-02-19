package main

import (
	"fmt"
	"log"
	"os"

	"stemcell-version-bump/cmd/out/runner"
	"stemcell-version-bump/resource"
)

func main() {
	buildDir := os.Args[1]

	err := os.Chdir(buildDir)
	if err != nil {
		log.Fatalf("Failed to move into build dir: %s", err)
	}

	request, err := resource.NewOutRequest(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to load request from stdin: %s", err)
	}

	client, err := resource.NewGCSClient(request.Source.JSONKey)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %s", err)
	}

	version, err := runner.NewVersion(request)
	if err != nil {
		log.Fatalf("Failed to read output file: %s", err)
	}

	err = runner.UploadVersion(request, client, version)
	if err != nil {
		log.Fatalf("Failed to upload version info to GCS: %s", err)
	}

	resourceVersion, err := runner.GenerateResourceOutput(version)
	if err != nil {
		log.Fatalf("Failed to upload version info to GCS: %s", err)
	}

	fmt.Println(resourceVersion)
}
