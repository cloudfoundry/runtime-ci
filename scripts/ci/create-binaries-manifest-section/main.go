package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/create-binaries-manifest-section/manifest"
)

func main() {
	var releases = []string{
		"capi",
		"consul",
		"diego",
		"etcd",
		"loggregator",
		"nats",
		"cf-mysql",
		"routing",
		"uaa",
		"garden-runc",
		"cflinuxfs2-rootfs",
		"binary-buildpack",
		"dotnet-core-buildpack",
		"go-buildpack",
		"java-buildpack",
		"nodejs-buildpack",
		"php-buildpack",
		"python-buildpack",
		"ruby-buildpack",
		"staticfile-buildpack",
	}

	deploymentConfigurationPath := os.Getenv("DEPLOYMENT_CONFIGURATION_PATH")
	deploymentManifestPath := os.Getenv("DEPLOYMENT_MANIFEST_PATH")

	var buildDir string
	flag.StringVar(&buildDir, "build-dir", "", "path to the build directory")
	flag.Parse()

	cfDeploymentManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "deployment-configuration", deploymentConfigurationPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	updatedDeploymentManifest, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filepath.Join(buildDir, "deployment-manifest", deploymentManifestPath), updatedDeploymentManifest, 0666); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
