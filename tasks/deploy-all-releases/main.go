package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/command"
	"github.com/spf13/pflag"
)

const (
	cfDeployment = "cf-deployment"
)

var (
	cfDeploymentDir string
	stemcellDir     string
)

func init() {
	pflag.Parse()
	rootDir := pflag.Arg(0)

	cfDeploymentDir = filepath.Join(rootDir, cfDeployment)
}

func main() {
	boshCLI := new(command.BoshCLI)

	content, _ := ioutil.ReadFile(filepath.Join(cfDeploymentDir, "cf-deployment.yml"))
	manifest, _ := bosh.NewManifestFromFile(content)
	releases := manifest.Releases

	stemcell, _ := bosh.NewStemcellFromInput(stemcellDir)

	newManifest := bosh.Manifest{
		Releases:  releases,
		Stemcells: []bosh.Stemcell{stemcell},
		Name:      "releases",
	}

	newManifest.Deploy(boshCLI)
}
