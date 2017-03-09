package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"
)

func main() {
	deploymentConfigurationPath := os.Getenv("DEPLOYMENT_CONFIGURATION_PATH")
	deploymentManifestPath := os.Getenv("DEPLOYMENT_MANIFEST_PATH")
	commitMessagePath := os.Getenv("COMMIT_MESSAGE_PATH")

	var buildDir string
	flag.StringVar(&buildDir, "build-dir", "", "path to the build directory")
	flag.Parse()

	files, err := ioutil.ReadDir(buildDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	releases := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "-release") {
			releases = append(releases, strings.TrimSuffix(file.Name(), "-release"))
		}
	}

	cfDeploymentManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "deployment-configuration", deploymentConfigurationPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	updatedDeploymentManifest, commitMessage, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filepath.Join(buildDir, "commit-message", commitMessagePath), []byte(commitMessage), 0666); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filepath.Join(buildDir, "deployment-manifest", deploymentManifestPath), updatedDeploymentManifest, 0666); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
