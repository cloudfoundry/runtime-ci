package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
)

func main() {
	buildDir := os.Args[1]

	content, err := ioutil.ReadFile(filepath.Join(buildDir, "cf-deployment-main", "cf-deployment.yml"))
	if err != nil {
		log.Fatalf("Failed to read main cf-deployment.yml: %s", err)
	}

	mainManifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		log.Fatalf("Failed to unmarshal main cf-deployment.yml: %s", err)
	}

	content, err = ioutil.ReadFile(filepath.Join(buildDir, "cf-deployment-release-candidate", "cf-deployment.yml"))
	if err != nil {
		log.Fatalf("Failed to read release-candidate cf-deployment.yml: %s", err)
	}

	releaseCandidateManifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		log.Fatalf("Failed to unmarshal release-candidate cf-deployment.yml: %s", err)
	}

	mainStemcell := mainManifest.Stemcells[0]
	releaseCandidateStemcell := releaseCandidateManifest.Stemcells[0]

	result, err := releaseCandidateStemcell.CompareVersion(mainStemcell)
	if err != nil {
		log.Fatalf("Failed to compare stemcell versions: %s", err)
	}

	if result == -1 {
		log.Fatalf("Release candidate stemcell version (%s) is behind the main stemcell version (%s). Aborting.",
			releaseCandidateStemcell.Version, mainStemcell.Version)
	}

	log.Printf("Release candidate stemcell version (%s) is ahead of, or equal to, the main stemcell version (%s). Proceeding.",
		releaseCandidateStemcell.Version, mainStemcell.Version)
}
