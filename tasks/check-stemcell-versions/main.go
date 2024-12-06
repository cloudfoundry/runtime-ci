package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <buildDir> <branchToCompare>", os.Args[0])
	}

	buildDir := os.Args[1]
	branchToCompare := os.Args[2]

	content, err := os.ReadFile(filepath.Join(buildDir, "cf-deployment-main", "cf-deployment.yml"))
	if err != nil {
		log.Fatalf("Failed to read main branch cf-deployment.yml: %s", err)
	}

	mainManifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		log.Fatalf("Failed to unmarshal main branch cf-deployment.yml: %s", err)
	}

	content, err = os.ReadFile(filepath.Join(buildDir, "cf-deployment-"+branchToCompare, "cf-deployment.yml"))
	if err != nil {
		log.Fatalf("Failed to read %s branch cf-deployment.yml: %s", branchToCompare, err)
	}

	branchToCompareManifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		log.Fatalf("Failed to unmarshal %s branch cf-deployment.yml: %s", branchToCompare, err)
	}

	mainStemcell := mainManifest.Stemcells[0]
	branchToCompareStemcell := branchToCompareManifest.Stemcells[0]

	if mainStemcell.OS != branchToCompareStemcell.OS {
		log.Printf("%s branch stemcell OS (%s) is different to the main branch stemcell OS (%s). Proceeding.",
			branchToCompare, branchToCompareStemcell.OS, mainStemcell.OS)
		os.Exit(0)
	}

	result, err := branchToCompareStemcell.CompareVersion(mainStemcell)
	if err != nil {
		log.Fatalf("Failed to compare stemcell versions: %s", err)
	}

	if result == -1 {
		log.Fatalf("%s branch stemcell version (%s) is behind the main branch stemcell version (%s). Aborting.",
			branchToCompare, branchToCompareStemcell.Version, mainStemcell.Version)
	}

	log.Printf("%s branch stemcell version (%s) is ahead of, or equal to, the main branch stemcell version (%s). Proceeding.",
		branchToCompare, branchToCompareStemcell.Version, mainStemcell.Version)
}
