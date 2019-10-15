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

	content, err := ioutil.ReadFile(filepath.Join(buildDir, "cf-deployment-master", "cf-deployment.yml"))
	if err != nil {
		log.Fatalf("Failed to read master cf-deployment.yml: %s", err)
	}

	masterManifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		log.Fatalf("Failed to unmarshal master cf-deployment.yml: %s", err)
	}

	content, err = ioutil.ReadFile(filepath.Join(buildDir, "cf-deployment-release-candidate", "cf-deployment.yml"))
	if err != nil {
		log.Fatalf("Failed to read release-candidate cf-deployment.yml: %s", err)
	}

	releaseCandidateManifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		log.Fatalf("Failed to unmarshal release-candidate cf-deployment.yml: %s", err)
	}

	masterStemcell := masterManifest.Stemcells[0]
	releaseCandidateStemcell := releaseCandidateManifest.Stemcells[0]

	result, err := releaseCandidateStemcell.CompareVersion(masterStemcell)
	if err != nil {
		log.Fatalf("Failed to compare stemcell versions: %s", err)
	}

	if result == -1 {
		log.Fatalf("Release candidate stemcell version (%s) is smaller than the master stemcell version (%s). Aborting.",
			releaseCandidateStemcell.Version, masterStemcell.Version)
	}

	log.Printf("Release candidate stemcell version (%s) is greater or equal to the master stemcell version (%s). Proceeding.",
		releaseCandidateStemcell.Version, masterStemcell.Version)
}
