package manifest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var yamlMarshal func(interface{}) ([]byte, error) = yaml.Marshal

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

	cfDeploymentReleasesAndStemcells := Manifest{}
	for _, release := range releases {
		releasePath := filepath.Join(buildDir, fmt.Sprintf("%s-release", release))

		sha1, err := ioutil.ReadFile(filepath.Join(releasePath, "sha1"))
		if err != nil {
			return nil, "", err
		}

		url, err := ioutil.ReadFile(filepath.Join(releasePath, "url"))
		if err != nil {
			return nil, "", err
		}

		version, err := ioutil.ReadFile(filepath.Join(releasePath, "version"))
		if err != nil {
			return nil, "", err
		}

		if releasesSHA1s[release] != strings.TrimSpace(string(sha1)) {
			changes = append(changes, fmt.Sprintf("%s-release", release))
		}

		cfDeploymentReleasesAndStemcells.Releases = append(cfDeploymentReleasesAndStemcells.Releases, Release{
			Name:    release,
			SHA1:    strings.TrimSpace(string(sha1)),
			Version: strings.TrimSpace(string(version)),
			URL:     strings.TrimSpace(string(url)),
		})
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

	cfDeploymentReleasesAndStemcellsYaml, err := yamlMarshal(cfDeploymentReleasesAndStemcells)
	if err != nil {
		return nil, "", err
	}

	changeMessage := "No release or stemcell version updates"
	if len(changes) > 0 {
		changeMessage = fmt.Sprintf("Updated %s", strings.Join(changes, ", "))
	}

	return append(cfDeploymentPreamble, cfDeploymentReleasesAndStemcellsYaml...), changeMessage, err
}
