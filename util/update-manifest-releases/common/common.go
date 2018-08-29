package common

import (
	"path/filepath"
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	"errors"
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

	sha1, shaErr := ioutil.ReadFile(filepath.Join(releasePath, "sha1"))
	url, urlErr := ioutil.ReadFile(filepath.Join(releasePath, "url"))
	version, verErr := ioutil.ReadFile(filepath.Join(releasePath, "version"))

	isShaErr := shaErr != nil
	isUrlErr := urlErr != nil

	// We accept neither or both of "sha1" and "url".  If we error out on only one or the other, something is wrong.
	if isShaErr != isUrlErr {
		if isShaErr {
			return Release{}, shaErr
		}
		return Release{}, urlErr
	}

	if verErr != nil {
		return Release{}, verErr
	}

	newRelease.URL = strings.TrimSpace(string(url))
	newRelease.SHA1 = strings.TrimSpace(string(sha1))
	newRelease.Version = strings.TrimSpace(string(version))

	return newRelease, nil
}

func StemcellInfoFromTarballName(tarballName string, releaseName string, releaseVersion string) (string, string, error) {
	// a valid tarball name is e.g. package-name-1.0-stemcell-name-2.0-45-23-44.tgz
	// ^package-name '-' package-version '-'  <stemcell name + version >-(\d+) -\d+-\d+-\d+$

	versionRegexString := fmt.Sprintf(`%s-%s-(.*)-([\d.]+)-\d+-\d+-\d+.tgz`,
		regexp.QuoteMeta(releaseName), regexp.QuoteMeta(releaseVersion))
	versionRegex := regexp.MustCompile(versionRegexString)

	allMatches := versionRegex.FindAllStringSubmatch(tarballName, 1)

	if len(allMatches) != 1 {
		return "", "", fmt.Errorf("invalid tarball name syntax: %s, %s", tarballName, versionRegexString)
	}

	if len(allMatches[0]) != 3 {
		return "", "", errors.New("internal error: len allMatches[0] should be 2, but it is not.")
	}

	versionMatch := allMatches[0][2]
	osMatch := allMatches[0][1]

	return versionMatch, osMatch, nil
}