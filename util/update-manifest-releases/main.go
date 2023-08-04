package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/compiledreleasesops"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
)

var cfDeploymentIgnoreDirs = []string{".git", ".github", "scripts", "example-vars-files", "iaas-support", "ci", "units"}
var cfDeploymentIgnoreFiles = []string{
	"cf-deployment.yml",
	".overcommit.yml",
	"use-compiled-releases.yml",
	"use-offline-windows2016fs.yml",
	"use-offline-windows1803fs.yml",
	"use-offline-windows2019fs.yml",
	"windows2016-cell.yml",
}

func getReleaseNames(buildDir string) ([]string, error) {
	files, err := os.ReadDir(buildDir)
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
	commitMessageFile := filepath.Join(buildDir, commitMessagePath)

	existingCommitMessage, err := os.ReadFile(commitMessageFile)

	if err != nil || strings.TrimSpace(string(existingCommitMessage)) == common.NoChangesCommitMessage {
		if err := os.WriteFile(commitMessageFile, []byte(commitMessage), 0666); err != nil {
			return err
		}
	}
	return nil
}

type updateFunc func([]string, string, []byte, common.MarshalFunc, common.UnmarshalFunc) ([]byte, string, error)

func isIgnored(name string, ignoreList []string) bool {
	for _, ignoreName := range ignoreList {
		if name == ignoreName {
			return true
		}
	}
	return false
}

func findOpsFiles(searchDir string, ignoreDirs, ignoreFiles []string) (map[string]string, error) {
	foundFiles := make(map[string]string)

	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && isIgnored(info.Name(), ignoreDirs) {
			return filepath.SkipDir
		} else if !info.IsDir() && filepath.Ext(info.Name()) == ".yml" && !isIgnored(info.Name(), ignoreFiles) {
			foundFiles[path] = strings.TrimPrefix(path, searchDir)
		}

		return nil
	})

	return foundFiles, err
}

func update(releases []string, inputPath, outputPath, inputDir, outputDir, buildDir, commitMessagePath string, f updateFunc) error {
	filesToUpdate := make(map[string]string)
	var err error
	ignoreNotFoundAndBadFormatErrors := false

	if inputPath == "" && outputPath == "" {
		filesToUpdate, err = findOpsFiles(filepath.Join(buildDir, inputDir), cfDeploymentIgnoreDirs, cfDeploymentIgnoreFiles)
		if err != nil {
			return err
		}
		ignoreNotFoundAndBadFormatErrors = true
	} else {
		filesToUpdate[filepath.Join(buildDir, inputDir, inputPath)] = outputPath
	}

	for inputPath, outputFileName := range filesToUpdate {
		fmt.Printf("Processing %s...\n", inputPath)
		originalFile, err := os.ReadFile(inputPath)
		if err != nil {
			return err
		}

		updatedFile, commitMessage, err := f(releases, buildDir, originalFile, yaml.Marshal, yaml.Unmarshal)
		if err != nil {
			isNotFoundError := strings.Contains(err.Error(), "Opsfile does not contain release named")
			isBadFormatError := err.Error() == opsfile.BadReleaseOpsFormatErrorMessage
			isNotFoundOrBadFormat := isNotFoundError || isBadFormatError

			if !(isNotFoundOrBadFormat && ignoreNotFoundAndBadFormatErrors) {
				return err
			}
		}

		if commitMessage != common.NoOpsFileChangesCommitMessage {
			if err := writeCommitMessage(buildDir, commitMessage, commitMessagePath); err != nil {
				return err
			}

			updatedOpsFilePath := filepath.Join(buildDir, outputDir, filepath.Dir(outputFileName))

			err := os.MkdirAll(updatedOpsFilePath, os.ModePerm)
			if err != nil {
				return err
			}

			fmt.Printf("Updating file: %s\n", inputPath)
			if err := os.WriteFile(filepath.Join(updatedOpsFilePath, filepath.Base(outputFileName)), updatedFile, 0666); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	var buildDir string
	flag.StringVar(&buildDir, "build-dir", "", "path to the build directory")

	var inputDir string
	flag.StringVar(&inputDir, "input-dir", "", "path to the input directory")

	var outputDir string
	flag.StringVar(&outputDir, "output-dir", "", "path to the output directory")

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
			inputDir,
			outputDir,
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
			inputDir,
			outputDir,
			buildDir,
			commitMessagePath,
			compiledreleasesops.UpdateCompiledReleases,
		); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else if target == "stemcell" {
		if err = update(
			releases,
			os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH"),
			os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH"),
			inputDir,
			outputDir,
			buildDir,
			commitMessagePath,
			manifest.UpdateStemcell,
		); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		if err = update(
			releases,
			os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH"),
			os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH"),
			inputDir,
			outputDir,
			buildDir,
			commitMessagePath,
			manifest.UpdateReleases,
		); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
