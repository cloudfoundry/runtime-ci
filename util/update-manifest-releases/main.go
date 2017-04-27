package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
)

func getReleaseNames(buildDir string) ([]string, error) {
	files, err := ioutil.ReadDir(buildDir)
	if err != nil {
		return nil, err

	}

	releases := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "-release") {
			releases = append(releases, strings.TrimSuffix(file.Name(), "-release"))
		}
	}
	return releases, nil
}

func main() {
	var inputPath, outputPath string
	commitMessagePath := os.Getenv("COMMIT_MESSAGE_PATH")

	var buildDir string
	flag.StringVar(&buildDir, "build-dir", "", "path to the build directory")

	var target string
	flag.StringVar(&target, "target", "manifest", "choose whether to update releases in manifest or opsfile")
	flag.Parse()

	releases, err := getReleaseNames(buildDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if target == "opsfile" {
		inputPath = os.Getenv("ORIGINAL_OPS_FILE_PATH")
		outputPath = os.Getenv("UPDATED_OPS_FILE_PATH")

		originalOpsFile, err := ioutil.ReadFile(filepath.Join(buildDir, "original-ops-file", inputPath))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		updatedOpsFile, commitMessage, err := opsfile.UpdateReleases(releases, buildDir, originalOpsFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		if err := ioutil.WriteFile(filepath.Join(buildDir, "commit-message", commitMessagePath), []byte(commitMessage), 0666); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		if err := ioutil.WriteFile(filepath.Join(buildDir, "updated-ops-file", outputPath), updatedOpsFile, 0666); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		inputPath = os.Getenv("DEPLOYMENT_CONFIGURATION_PATH")
		outputPath = os.Getenv("DEPLOYMENT_MANIFEST_PATH")

		cfDeploymentManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "deployment-configuration", inputPath))
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

		if err := ioutil.WriteFile(filepath.Join(buildDir, "deployment-manifest", outputPath), updatedDeploymentManifest, 0666); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
