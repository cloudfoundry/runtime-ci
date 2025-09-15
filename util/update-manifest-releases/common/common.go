package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const NoChangesCommitMessage = "No manifest release or stemcell version updates"
const NoOpsFileChangesCommitMessage = "No opsfile release updates"

type MarshalFunc func(interface{}) ([]byte, error)
type UnmarshalFunc func([]byte, interface{}) error

type Release struct {
	Name     string             `yaml:"name"`
	URL      string             `yaml:"url"`
	Version  string             `yaml:"version"`
	SHA1     string             `yaml:"sha1"`
	Stemcell StemcellForRelease `yaml:"stemcell,omitempty"`
}

type StemcellForRelease struct {
	OS      string `yaml:"os,omitempty"`
	Version string `yaml:"version,omitempty"`
}

func GetReleaseFromFile(buildDir, releaseName string) (Release, error) {
	newRelease := Release{
		Name: releaseName,
	}
	releasePath := filepath.Join(buildDir, fmt.Sprintf("%s-release", releaseName))

	version, verErr := os.ReadFile(filepath.Join(releasePath, "version"))

	if verErr != nil {
		return Release{}, verErr
	}

	newRelease.Version = strings.TrimSpace(string(version))

	url, urlErr := os.ReadFile(filepath.Join(releasePath, "url"))

	if urlErr != nil {
		return Release{}, urlErr
	}

	_, commitShaErr := os.ReadFile(filepath.Join(releasePath, "commit_sha"))

	if commitShaErr != nil {
		// Bosh.io release
		fmt.Println("Missing commit_sha file. Assuming bosh.io release...")
		sha1, shaErr := os.ReadFile(filepath.Join(releasePath, "sha256"))

		if shaErr != nil {
			return Release{}, shaErr
		}

		newRelease.SHA1 = strings.TrimSpace("sha256:" + string(sha1))
		newRelease.URL = strings.TrimSpace(string(url))
	} else {
		// Github release
		fmt.Println("Found commit_sha file. Assuming github release...")
	}

	return newRelease, nil
}

func InfoFromTarballName(tarballName string, releaseName string) (string, string, string, error) {
	// a valid tarball name is e.g. package-name-1.0-stemcell-name-2.0-45-23-44.tgz
	// ^package-name '-' package-version '-'  <stemcell name + version >-(\d+) -\d+-\d+-\d+$

	versionRegexString := fmt.Sprintf(`%s-([\d.]+)-(.*)-([\d.]+)-\d+-\d+-\d+.tgz`,
		regexp.QuoteMeta(releaseName))
	versionRegex := regexp.MustCompile(versionRegexString)

	allMatches := versionRegex.FindAllStringSubmatch(tarballName, 1)

	if len(allMatches) != 1 {
		return "", "", "", fmt.Errorf("invalid tarball name syntax: %s, %s", tarballName, versionRegexString)
	}

	if len(allMatches[0]) != 4 {
		return "", "", "", errors.New("internal error: len allMatches[0] should be 4, but it is not")
	}

	releaseVersionMatch := allMatches[0][1]
	stemcellOsMatch := allMatches[0][2]
	stemcellVersionMatch := allMatches[0][3]

	return releaseVersionMatch, stemcellVersionMatch, stemcellOsMatch, nil
}
