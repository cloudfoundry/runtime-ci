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
	stemcell     = "stemcell"
)

var (
	cfDeploymentDir string
	stemcellDir     string
)

func init() {
	pflag.Parse()
	rootDir := pflag.Arg(0)

	cfDeploymentDir = filepath.Join(rootDir, cfDeployment)

	stemcellDir = filepath.Join(rootDir, stemcell)
}

func main() {
	boshCLI := new(command.BoshCLI)

	fmt.Println("Reading cf-deployment")
	content, err := ioutil.ReadFile(filepath.Join(cfDeploymentDir, "cf-deployment.yml"))
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

	fmt.Println("Reading stemcell")
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

		fmt.Println("Deploying manifest")

		if err := newManifest.Deploy(boshCLI, fmt.Sprintf("%s-compilation", release.Name)); err != nil {
			fmt.Println(err)
		}
	}
}
