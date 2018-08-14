package opsfile

import (
	"errors"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"fmt"
	"strings"
	"path/filepath"
	"regexp"
	"io/ioutil"
	"crypto/sha1"
)

const compiledReleasesURLPrefix = "https://storage.googleapis.com/cf-deployment-compiled-releases"

func UpdateCompiledReleases(releaseNames []string, buildDir string, opsFile []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	if releaseNames == nil || len(releaseNames) == 0 {
		err := errors.New("releaseNames provided to UpdateReleases must contain at least one release name")
		return nil, "", err
	}

	var deserializedOpsFile []Op
	if err := unmarshalFunc(opsFile, &deserializedOpsFile); err != nil {
		return nil, "", err
	}

	for i, op := range deserializedOpsFile {
		for _, releaseName := range releaseNames {
			newRelease, err := getCompiledReleaseForBuild(buildDir, releaseName)
			if err != nil {
				return nil, "", err
			}

			matchingURLPath := fmt.Sprintf("/releases/name=%s/url", releaseName)
			matchingVersionPath := fmt.Sprintf("/releases/name=%s/version", releaseName)
			matchingSha1Path := fmt.Sprintf("/releases/name=%s/sha1", releaseName)
			matchingStemcellPath := fmt.Sprintf("/releases/name=%s/stemcell?", releaseName)

			if strings.Contains(op.Path, matchingURLPath) {
				deserializedOpsFile[i].Value = newRelease.URL
			} else if strings.Contains(op.Path, matchingVersionPath) {
				deserializedOpsFile[i].Value = newRelease.Version
			} else if strings.Contains(op.Path, matchingSha1Path) {
				deserializedOpsFile[i].Value = newRelease.SHA1
			} else if strings.Contains(op.Path, matchingStemcellPath) {
				deserializedOpsFile[i].Value = newRelease.Stemcell
			}
		}
	}

	updatedOpsFile, err := marshalFunc(&deserializedOpsFile)
	if err != nil {
		return nil, "", err
	}

	return updatedOpsFile, "", nil
}

func getCompiledReleaseForBuild(buildDir, releaseName string) (common.Release, error) {
	releaseTarballGlob := filepath.Join(buildDir, fmt.Sprintf("%s-compiled-release-tarball", releaseName), "*.tgz")

	matches, err := filepath.Glob(releaseTarballGlob)
	if err != nil {
		return common.Release{}, err
	}

	if len(matches) != 1 {
		return common.Release{}, errors.New("expected to find exactly 1 compiled release tarball")
	}

	return getCompiledReleaseInfo(matches[0])
}

func getCompiledReleaseInfo(releaseTarballPath string) (common.Release, error) {
	var compiledRelease common.Release

	releaseTarballName := filepath.Base(releaseTarballPath)
	compiledRelease.URL = fmt.Sprintf("%s/%s", compiledReleasesURLPrefix, releaseTarballName)

	re := regexp.MustCompile(`([a-z]+-*[a-z]*)-(\d+\.\d+\.\d+|\d+\.\d+)`)
	matches := re.FindAllStringSubmatch(releaseTarballName, -1)

	if len(matches) != 2 {
		return common.Release{},
			errors.New("release tarball should be of the form <release-name>-<release-version>-<stemcell-os>-<stemcell-version>-<timestamp>.tgz")
	}

	compiledRelease.Name = matches[0][1]
	compiledRelease.Version = matches[0][2]
	compiledRelease.Stemcell.OS = matches[1][1]
	compiledRelease.Stemcell.Version = matches[1][2]
	releaseSha1Sum, err := computeSha1Sum(releaseTarballPath)
	if err != nil {
		return common.Release{}, err
	}
	compiledRelease.SHA1 = releaseSha1Sum

	return compiledRelease, nil
}

func computeSha1Sum(filepath string) (string, error) {
	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(fileContents)), nil
}