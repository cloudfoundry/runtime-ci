package manifest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
)

type Stemcell struct {
	Alias   string `yaml:"alias"`
	OS      string `yaml:"os"`
	Version string `yaml:"version"`
}

type Manifest struct {
	Releases  []common.Release `yaml:"releases"`
	Stemcells []Stemcell       `yaml:"stemcells"`
}

func mergeReleases(manifestReleases []common.Release, updatingReleases []string) []common.Release {
	manifestReleaseMap := map[string]bool{}
	for _, r := range manifestReleases {
		manifestReleaseMap[r.Name] = true
	}

	allReleases := manifestReleases
	for _, release := range updatingReleases {
		if _, found := manifestReleaseMap[release]; !found {
			allReleases = append(allReleases, common.Release{Name: release})
		}
	}
	return allReleases
}

func UpdateReleasesAndStemcells(releases []string, buildDir string, cfDeploymentManifest []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	r := regexp.MustCompile(`(?m:^releases:$)`)

	submatches := r.FindSubmatchIndex([]byte(cfDeploymentManifest))

	if len(submatches) == 0 {
		return nil, "", fmt.Errorf("releases was not found at the bottom of the manifest")
	}

	cfDeploymentManifestReleasesIndex := submatches[0]

	cfDeploymentPreamble := make([]byte, len(cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]))
	copy(cfDeploymentPreamble, cfDeploymentManifest[:cfDeploymentManifestReleasesIndex])

	var deserializedManifestSuffix map[string]interface{}
	if err := unmarshalFunc(cfDeploymentManifest[cfDeploymentManifestReleasesIndex:], &deserializedManifestSuffix); err != nil {
		return nil, "", err
	}

	if len(deserializedManifestSuffix) > 2 {
		return nil, "", fmt.Errorf(`found keys other than "releases" and "stemcells" at the bottom of the manifest`)
	}

	if _, ok := deserializedManifestSuffix["stemcells"]; !ok {
		return nil, "", fmt.Errorf("stemcells was not found at the bottom of the manifest")
	}

	var releasesAndStemcells Manifest
	if err := unmarshalFunc(cfDeploymentManifest[cfDeploymentManifestReleasesIndex:], &releasesAndStemcells); err != nil {
		return nil, "", err
	}

	releasesByName := make(map[string]common.Release)
	for _, release := range releasesAndStemcells.Releases {
		releasesByName[release.Name] = release
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

	var changes []string
	cfDeploymentReleasesAndStemcells := Manifest{}
	for _, release := range releasesAndStemcells.Releases {
		var newRelease common.Release

		if _, found := releaseMap[release.Name]; found {
			var err error

			newRelease, err = common.GetReleaseFromFile(buildDir, release.Name)
			if err != nil {
				return nil, "", err
			}

			if releasesByName[newRelease.Name] != newRelease {
				changes = append(changes, fmt.Sprintf("%s-release %s", newRelease.Name, newRelease.Version))
			}
		} else {
			newRelease = release
		}
		cfDeploymentReleasesAndStemcells.Releases = append(cfDeploymentReleasesAndStemcells.Releases, newRelease)
	}

	stemcellVersion, err := ioutil.ReadFile(filepath.Join(buildDir, "stemcell", "version"))
	if err != nil {
		return nil, "", err
	}
	trimmedStemcellVersion := strings.TrimSpace(string(stemcellVersion))

	cfDeploymentReleasesAndStemcells.Stemcells = []Stemcell{
		{
			Alias:   "default",
			OS:      "ubuntu-trusty",
			Version: trimmedStemcellVersion,
		},
	}

	if stemcellsVersions["default"] != trimmedStemcellVersion {
		changes = append(changes, fmt.Sprintf("ubuntu-trusty stemcell %s", trimmedStemcellVersion))
	}

	cfDeploymentReleasesAndStemcellsYaml, err := marshalFunc(cfDeploymentReleasesAndStemcells)
	if err != nil {
		return nil, "", err
	}

	changeMessage := common.NoChangesCommitMessage
	if len(changes) > 0 {
		changeMessage = fmt.Sprintf("Updated manifest with %s", strings.Join(changes, ", "))
	}

	return append(cfDeploymentPreamble, cfDeploymentReleasesAndStemcellsYaml...), changeMessage, err
}
