package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/compiledreleasesops"
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

func writeCommitMessage(buildDir, commitMessage, commitMessagePath string) error {
	commitMessageFile := filepath.Join(buildDir, "commit-message", commitMessagePath)

	// We are ignoring error here
	// because we don't care if commit message file does not exist
	existingCommitMessage, err := ioutil.ReadFile(commitMessageFile)

	if err != nil || string(existingCommitMessage) == common.NoChangesCommitMessage {
		if err := ioutil.WriteFile(commitMessageFile, []byte(commitMessage), 0666); err != nil {
			return err
		}
	}
	return nil
}

type updateFunc func([]string, string, []byte, common.MarshalFunc, common.UnmarshalFunc) ([]byte, string, error)

func update(releases []string, inputPath, outputPath, inputDir, outputDir, buildDir, commitMessagePath string, f updateFunc) error {
	originalFile, err := ioutil.ReadFile(filepath.Join(buildDir, inputDir, inputPath))
	if err != nil {
		return err
	}

	updatedFile, commitMessage, err := f(releases, buildDir, originalFile, yaml.Marshal, yaml.Unmarshal)
	if err != nil {
		return err
	}

	if err := writeCommitMessage(buildDir, commitMessage, commitMessagePath); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(buildDir, outputDir, outputPath), updatedFile, 0666); err != nil {
		return err
	}

	return nil
}

func main() {
	var buildDir string
	flag.StringVar(&buildDir, "build-dir", "", "path to the build directory")

	var release string
	flag.StringVar(&release, "release", "", "name of release, without -release suffix")

	var target string
	flag.StringVar(&target, "target", "manifest", "choose whether to update releases in manifest or opsfile")
	flag.Parse()

	var err error
	releases := []string{release}
	if release == "" {
		releases, err = getReleaseNames(buildDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	commitMessagePath := os.Getenv("COMMIT_MESSAGE_PATH")

	if target == "opsfile" {
		if err = update(
			releases,
			os.Getenv("ORIGINAL_OPS_FILE_PATH"),
			os.Getenv("UPDATED_OPS_FILE_PATH"),
			"original-ops-file",
			"updated-ops-file",
			buildDir,
			commitMessagePath,
			opsfile.UpdateReleases,
		); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else if target == "compiledReleasesOpsfile" {
		if err = update(
			releases,
			os.Getenv("ORIGINAL_OPS_FILE_PATH"),
			os.Getenv("UPDATED_OPS_FILE_PATH"),
			"original-ops-file",
			"updated-ops-file",
			buildDir,
			commitMessagePath,
			compiledreleasesops.UpdateCompiledReleases,
		); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		if err = update(
			releases,
			os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH"),
			os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH"),
			"deployment-configuration",
			"updated-deployment-manifest",
			buildDir,
			commitMessagePath,
			manifest.UpdateReleasesAndStemcells,
		); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
