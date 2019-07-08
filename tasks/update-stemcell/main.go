package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/compiledreleasesops"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
	"gopkg.in/yaml.v2"
)

type updateFunc func([]string, string, []byte, common.MarshalFunc, common.UnmarshalFunc) ([]byte, string, error)

func writeCommitMessage(buildDir, commitMessage, commitMessagePath string) error {
	commitMessageFile := filepath.Join(buildDir, commitMessagePath)

	existingCommitMessage, err := ioutil.ReadFile(commitMessageFile)

	if err != nil || strings.TrimSpace(string(existingCommitMessage)) == common.NoChangesCommitMessage {
		if err := ioutil.WriteFile(commitMessageFile, []byte(commitMessage), 0666); err != nil {
			return err
		}
	}
	return nil
}

func getReleaseNames(inputDeploymentmanifestPath string) ([]string, error) {
	fileContents, err := ioutil.ReadFile(inputDeploymentmanifestPath)
	if err != nil {
		return nil, err
	}

	var manifest manifest.Manifest
	err = yaml.Unmarshal(fileContents, &manifest)
	if err != nil {
		return nil, err
	}

	releases := []string{}
	for _, release := range manifest.Releases {
		releases = append(releases, release.Name)
	}

	return releases, nil
}

func update(releases []string, inputPath, outputPath, inputDir, outputDir, buildDir, commitMessagePath string, f updateFunc) error {
	inputFilePath := filepath.Join(buildDir, inputDir, inputPath)

	var err error
	fmt.Printf("Processing %s...\n", inputFilePath)
	originalFile, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	updatedFile, commitMessage, err := f(releases, buildDir, originalFile, yaml.Marshal, yaml.Unmarshal)
	if err != nil {
		isNotFoundError := strings.Contains(err.Error(), "Opsfile does not contain release named")
		isBadFormatError := err.Error() == opsfile.BadReleaseOpsFormatErrorMessage
		isNotFoundOrBadFormat := isNotFoundError || isBadFormatError

		if !isNotFoundOrBadFormat {
			return err
		}
	}

	if commitMessage != common.NoOpsFileChangesCommitMessage {
		if err := writeCommitMessage(buildDir, commitMessage, commitMessagePath); err != nil {
			return err
		}

		updatedOpsFilePath := filepath.Join(buildDir, outputDir, filepath.Dir(outputPath))

		err := os.MkdirAll(updatedOpsFilePath, os.ModePerm)
		if err != nil {
			return err
		}

		fmt.Printf("Updating file: %s\n", inputFilePath)
		if err := ioutil.WriteFile(filepath.Join(updatedOpsFilePath, filepath.Base(outputPath)), updatedFile, 0666); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// TODO: remove pseudo-code comments:
	// parse flags
	var buildDir string
	flag.StringVar(&buildDir, "build-dir", "", "path to the build directory")

	var commitMessageFile string
	flag.StringVar(&commitMessageFile, "commit-message-file", "", "path to the commit message file")

	var compiledReleaseTarballDir string
	flag.StringVar(&compiledReleaseTarballDir, "compiled-release-tarball-dir", "", "path to the compiled releases tarballs directory")

	var stemcellDir string
	flag.StringVar(&stemcellDir, "stemcell-dir", "", "path to the stemcell directory")

	var originalCfdDir string
	flag.StringVar(&originalCfdDir, "original-cf-d-dir", "", "path to the original cf-deployment")

	var updatedCfdDir string
	flag.StringVar(&updatedCfdDir, "updated-cf-d-dir", "", "path to the updated cf-deployment")

	var target string
	flag.StringVar(&target, "target", "manifest", "choose whether to update releases in manifest or opsfile")

	flag.Parse()

	if buildDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: build-dir")
		os.Exit(1)
	}

	if compiledReleaseTarballDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: compiled-release-tarball-dir")
		os.Exit(1)
	}

	if stemcellDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: stemcell-dir")
		os.Exit(1)
	}

	if originalCfdDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: original-cf-d-dir")
		os.Exit(1)
	}

	if updatedCfdDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: updated-cf-d-dir")
		os.Exit(1)
	}

	inputDeploymentmanifestPath := os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH")
	outputDeploymentmanifestPath := os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH")

	inputCompiledReleasesOpsFilePath := os.Getenv("ORIGINAL_OPS_FILE_PATH")
	outputCompiledReleasesOpsFilePath := os.Getenv("UPDATED_OPS_FILE_PATH")

	commitMessagePath := os.Getenv("COMMIT_MESSAGE_PATH")

	if inputDeploymentmanifestPath == "" {
		fmt.Fprintln(os.Stderr, "missing path to input deployment manifest")
		os.Exit(1)
	}

	if outputDeploymentmanifestPath == "" {
		fmt.Fprintln(os.Stderr, "missing path to output deployment manifest")
		os.Exit(1)
	}

	if inputCompiledReleasesOpsFilePath == "" {
		fmt.Fprintln(os.Stderr, "missing path to input compiled release ops-file")
		os.Exit(1)
	}

	if outputCompiledReleasesOpsFilePath == "" {
		fmt.Fprintln(os.Stderr, "missing path to output compiled release ops-file")
		os.Exit(1)
	}

	// update stemcell in base manifest
	if err := update(
		[]string{},
		inputDeploymentmanifestPath,
		outputDeploymentmanifestPath,
		originalCfdDir,
		updatedCfdDir,
		buildDir,
		commitMessagePath,
		UpdateStemcell,
	); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// generate list of compiled releases
	releases, err := getReleaseNames(filepath.Join(buildDir, originalCfdDir, inputDeploymentmanifestPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// for each compiled release:
	//	 retrieve name/version from tarball (bosh inspect-local-release)
	//	 append to list of ops-file entries
	// marshal and write ops-file
	if err := update(
		releases,
		inputCompiledReleasesOpsFilePath,
		outputCompiledReleasesOpsFilePath,
		originalCfdDir,
		updatedCfdDir,
		buildDir,
		commitMessagePath,
		compiledreleasesops.UpdateCompiledReleases,
	); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func UpdateStemcell(releases []string, buildDir string, cfDeploymentManifest []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	return manifest.UpdateReleasesOrStemcell([]string{}, buildDir, cfDeploymentManifest, true, marshalFunc, unmarshalFunc)
}
