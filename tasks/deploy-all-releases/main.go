package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/command"
	"github.com/spf13/pflag"
)

const (
	cfDeployment = "cf-deployment"
	releaseList  = "release-list"
	stemcell     = "stemcell"
)

var (
	cfDeploymentDir string
	releaseListDir  string
	stemcellDir     string
)

func init() {
	pflag.Parse()
	rootDir := pflag.Arg(0)

	cfDeploymentDir = filepath.Join(rootDir, cfDeployment)
	releaseListDir = filepath.Join(rootDir, releaseList)
	stemcellDir = filepath.Join(rootDir, stemcell)
}

func main() {
	boshCLI := new(command.BoshCLI)

	releaseListPath := filepath.Join(releaseListDir, "releases.yml")
	var fileToRead string

	_, err := os.Stat(releaseListPath)
	if err == nil {
		fmt.Println("Reading releases.yml...")
		fileToRead = releaseListPath
	} else if os.IsNotExist(err) {
		fmt.Println("Reading cf-deployment.yml...")
		fileToRead = filepath.Join(cfDeploymentDir, "cf-deployment.yml")
	} else {
		fmt.Println(err)
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(fileToRead)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	manifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	releases := manifest.Releases

	fmt.Println("Reading stemcell information...")
	stemcell, err := bosh.NewStemcellFromInput(stemcellDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, release := range releases {
		newManifest := bosh.Manifest{
			Releases:  []bosh.Release{release},
			Stemcells: []bosh.Stemcell{stemcell},
			Name:      fmt.Sprintf("%s-compilation", release.Name),
		}

		fmt.Printf("Deploying %s...\n", release.Name)

		if err := newManifest.Deploy(boshCLI); err != nil {
			fmt.Println(err)
		}
	}
}
