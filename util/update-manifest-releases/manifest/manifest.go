package manifest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var YamlMarshal func(interface{}) ([]byte, error) = yaml.Marshal

type Stemcell struct {
	Alias   string `yaml:"alias"`
	OS      string `yaml:"os"`
	Version string `yaml:"version"`
}

type Release struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Version string `yaml:"version"`
	SHA1    string `yaml:"sha1"`
}

type Manifest struct {
	Releases  []Release  `yaml:"releases"`
	Stemcells []Stemcell `yaml:"stemcells"`
}

func mergeReleases(manifestReleases []Release, updatingReleases []string) []Release {
	manifestReleaseMap := map[string]bool{}
	for _, r := range manifestReleases {
		manifestReleaseMap[r.Name] = true
	}

	allReleases := manifestReleases
	for _, release := range updatingReleases {
		if _, found := manifestReleaseMap[release]; !found {
			allReleases = append(allReleases, Release{Name: release})
		}
	}
	return allReleases
}

func UpdateReleasesAndStemcells(releases []string, buildDir string, cfDeploymentManifest []byte) ([]byte, string, error) {
	changes := []string{}
	r := regexp.MustCompile(`(?m:^releases:$)`)

	submatches := r.FindSubmatchIndex([]byte(cfDeploymentManifest))

	if len(submatches) == 0 {
		return nil, "", fmt.Errorf("releases was not found at the bottom of the manifest")
	}

	cfDeploymentManifestReleasesIndex := submatches[0]

	cfDeploymentPreamble := make([]byte, len(cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]))
	copy(cfDeploymentPreamble, cfDeploymentManifest[:cfDeploymentManifestReleasesIndex])

	var deserializedManifestSuffix map[string]interface{}
	if err := yaml.Unmarshal(cfDeploymentManifest[cfDeploymentManifestReleasesIndex:], &deserializedManifestSuffix); err != nil {
		return nil, "", err
	}

	if len(deserializedManifestSuffix) > 2 {
		return nil, "", fmt.Errorf(`found keys other than "releases" and "stemcells" at the bottom of the manifest`)
	}

	if _, ok := deserializedManifestSuffix["stemcells"]; !ok {
		return nil, "", fmt.Errorf("stemcells was not found at the bottom of the manifest")
	}

	var releasesAndStemcells Manifest
	if err := yaml.Unmarshal(cfDeploymentManifest[cfDeploymentManifestReleasesIndex:], &releasesAndStemcells); err != nil {
		return nil, "", err
	}

	releasesSHA1s := map[string]string{}
	for _, release := range releasesAndStemcells.Releases {
		releasesSHA1s[release.Name] = release.SHA1
	}

	stemcellsVersions := map[string]string{}
	for _, stemcell := range releasesAndStemcells.Stemcells {
		stemcellsVersions[stemcell.Alias] = stemcell.Version
	}

	releaseMap := map[string]bool{}
	for _, r := range releases {
		releaseMap[r] = true
	}

	releasesAndStemcells.Releases = mergeReleases(releasesAndStemcells.Releases, releases)

	newRelease := Release{}
	cfDeploymentReleasesAndStemcells := Manifest{}
	for _, release := range releasesAndStemcells.Releases {
		if _, found := releaseMap[release.Name]; found {
			releasePath := filepath.Join(buildDir, fmt.Sprintf("%s-release", release.Name))

			sha1, err := ioutil.ReadFile(filepath.Join(releasePath, "sha1"))
			newRelease.SHA1 = strings.TrimSpace(string(sha1))
			if err != nil {
				return nil, "", err
			}

			newRelease.Name = release.Name
			if releasesSHA1s[newRelease.Name] != strings.TrimSpace(string(sha1)) {
				changes = append(changes, fmt.Sprintf("%s-release", newRelease.Name))
			}

			url, err := ioutil.ReadFile(filepath.Join(releasePath, "url"))
			if err != nil {
				return nil, "", err
			}
			newRelease.URL = strings.TrimSpace(string(url))

			version, err := ioutil.ReadFile(filepath.Join(releasePath, "version"))
			if err != nil {
				return nil, "", err
			}
			newRelease.Version = strings.TrimSpace(string(version))

		} else {
			newRelease = release
		}
		cfDeploymentReleasesAndStemcells.Releases = append(cfDeploymentReleasesAndStemcells.Releases, newRelease)
	}

	stemcellVersion, err := ioutil.ReadFile(filepath.Join(buildDir, "stemcell", "version"))
	if err != nil {
		return nil, "", err
	}

	cfDeploymentReleasesAndStemcells.Stemcells = []Stemcell{
		{
			Alias:   "default",
			OS:      "ubuntu-trusty",
			Version: strings.TrimSpace(string(stemcellVersion)),
		},
	}

	if stemcellsVersions["default"] != strings.TrimSpace(string(stemcellVersion)) {
		changes = append(changes, "ubuntu-trusty stemcell")
	}

	cfDeploymentReleasesAndStemcellsYaml, err := YamlMarshal(cfDeploymentReleasesAndStemcells)
	if err != nil {
		return nil, "", err
	}

	changeMessage := "No release or stemcell version updates"
	if len(changes) > 0 {
		changeMessage = fmt.Sprintf("Updated %s", strings.Join(changes, ", "))
	}

	return append(cfDeploymentPreamble, cfDeploymentReleasesAndStemcellsYaml...), changeMessage, err
}
