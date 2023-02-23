package manifest

import (
	"fmt"
	"os"
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

func stemcellOSfromURL(url string) (string, error) {
	// url example:
	// https://s3.amazonaws.com/bosh-gce-light-stemcells/light-bosh-stemcell-0.1-google-kvm-ubuntu-foo-go_agent.tgz
	urlSplit := strings.Split(url, "/")
	tarballName := urlSplit[len(urlSplit)-1]

	versionRegex := regexp.MustCompile(`(ubuntu-\w+)`)

	allMatches := versionRegex.FindAllStringSubmatch(tarballName, 1)

	if len(allMatches) != 1 {
		return "", fmt.Errorf("Stemcell URL does not contain 'ubuntu': %s", strings.Trim(url, "\n"))
	}

	osMatch := allMatches[0][1]

	return osMatch, nil
}

func updateReleasesOrStemcell(releases []string, buildDir string, cfDeploymentManifest []byte, stemcellBump bool, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
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

	releasesAndStemcells.Releases = mergeReleases(releasesAndStemcells.Releases, releases)

	var changes []string
	cfDeploymentReleasesAndStemcells := Manifest{}

	if stemcellBump {
		stemcellsVersions := map[string]string{}
		for _, stemcell := range releasesAndStemcells.Stemcells {
			stemcellsVersions[stemcell.Alias] = stemcell.Version
		}

		stemcellVersion, err := os.ReadFile(filepath.Join(buildDir, "stemcell", "version"))
		if err != nil {
			return nil, "", err
		}
		trimmedStemcellVersion := strings.TrimSpace(string(stemcellVersion))

		stemcellURL, err := os.ReadFile(filepath.Join(buildDir, "stemcell", "url"))
		if err != nil {
			return nil, "", err
		}

		stemcellOS, err := stemcellOSfromURL(string(stemcellURL))
		if err != nil {
			return nil, "", err
		}

		cfDeploymentReleasesAndStemcells.Releases = releasesAndStemcells.Releases
		cfDeploymentReleasesAndStemcells.Stemcells = []Stemcell{
			{
				Alias:   "default",
				OS:      stemcellOS,
				Version: trimmedStemcellVersion,
			},
		}

		if stemcellsVersions["default"] != trimmedStemcellVersion {
			changes = append(changes, fmt.Sprintf("%s stemcell %s", stemcellOS, trimmedStemcellVersion))
		}
	} else {
		releasesByName := make(map[string]common.Release)
		for _, release := range releasesAndStemcells.Releases {
			releasesByName[release.Name] = release
		}

		releaseMap := map[string]bool{}
		for _, r := range releases {
			releaseMap[r] = true
		}

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

		cfDeploymentReleasesAndStemcells.Stemcells = releasesAndStemcells.Stemcells
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

func UpdateReleases(releases []string, buildDir string, cfDeploymentManifest []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	return updateReleasesOrStemcell(releases, buildDir, cfDeploymentManifest, false, marshalFunc, unmarshalFunc)
}

func UpdateStemcell(releases []string, buildDir string, cfDeploymentManifest []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	return updateReleasesOrStemcell(releases, buildDir, cfDeploymentManifest, true, marshalFunc, unmarshalFunc)
}
