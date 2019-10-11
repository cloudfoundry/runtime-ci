package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudfoundry/runtime-ci/tasks/cf-deployment-minor-stemcell-bump-release-notes/concourseio"
)

func main() {
	buildDir := os.Args[1]
	runner, err := concourseio.NewRunner(buildDir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println("Reading default stemcell info from manifest...")
	oldStemcell, err := runner.ReadStemcellInfoFromManifest("default")
	if err != nil {
		log.Fatalf("Failed to read default stemcell info from manifest: %s", err)
	}

	fmt.Println("Reading new stemcell version from resource...")
	newStemcell, err := runner.ReadStemcellFromResource()
	if err != nil {
		log.Fatalf("Failed to read new stemcell version from resource: %s", err)
	}

	fmt.Println("Validating stemcell bump...")
	err = runner.ValidateStemcellBump(oldStemcell, newStemcell)
	if err != nil {
		log.Fatalf("Failed to validate stemcell bump: %s", err)
	}

	fmt.Println("Generating release notes...")
	err = runner.GenerateReleaseNotes(oldStemcell, newStemcell)
	if err != nil {
		log.Fatalf("Failed to generate release notes: %s", err)
	}

	fmt.Println("Generating release name...")
	err = runner.GenerateReleaseName()
	if err != nil {
		log.Fatalf("Failed to generate release name: %s", err)
	}
}
