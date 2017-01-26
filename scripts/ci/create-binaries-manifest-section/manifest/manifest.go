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
var yamlUnmarshal func([]byte, interface{}) error = yaml.Unmarshal

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

func UpdateReleasesAndStemcells(releases []string, buildDir string, cfDeploymentManifest []byte) ([]byte, error) {
	r := regexp.MustCompile(`(?m:^releases:$)`)

	submatches := r.FindSubmatchIndex([]byte(cfDeploymentManifest))

	if len(submatches) == 0 {
		return nil, fmt.Errorf("releases was not found at the bottom of the manifest")
	}

	cfDeploymentManifestReleasesIndex := submatches[0]

	cfDeploymentPreamble := make([]byte, len(cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]))
	copy(cfDeploymentPreamble, cfDeploymentManifest[:cfDeploymentManifestReleasesIndex])

	var deserializedManifestSuffix map[string]interface{}
	if err := yamlUnmarshal(cfDeploymentManifest[cfDeploymentManifestReleasesIndex:], &deserializedManifestSuffix); err != nil {
		return nil, err
	}

	if len(deserializedManifestSuffix) > 2 {
		return nil, fmt.Errorf(`found keys other than "releases" and "stemcells" at the bottom of the manifest`)
	}

	if _, ok := deserializedManifestSuffix["stemcells"]; !ok {
		return nil, fmt.Errorf("stemcells was not found at the bottom of the manifest")
	}

	cfDeploymentReleasesAndStemcells := Manifest{}

	for _, release := range releases {
		releasePath := filepath.Join(buildDir, fmt.Sprintf("%s-release", release))

		sha, err := ioutil.ReadFile(filepath.Join(releasePath, "sha1"))
		if err != nil {
			return nil, err
		}

		url, err := ioutil.ReadFile(filepath.Join(releasePath, "url"))
		if err != nil {
			return nil, err
		}

		version, err := ioutil.ReadFile(filepath.Join(releasePath, "version"))
		if err != nil {
			return nil, err
		}

		cfDeploymentReleasesAndStemcells.Releases = append(cfDeploymentReleasesAndStemcells.Releases, Release{
			Name:    release,
			SHA1:    strings.TrimSpace(string(sha)),
			Version: strings.TrimSpace(string(version)),
			URL:     strings.TrimSpace(string(url)),
		})
	}

	stemcellVersion, err := ioutil.ReadFile(filepath.Join(buildDir, "stemcell", "version"))
	if err != nil {
		return nil, err
	}

	cfDeploymentReleasesAndStemcells.Stemcells = []Stemcell{
		{
			Alias:   "default",
			OS:      "ubuntu-trusty",
			Version: strings.TrimSpace(string(stemcellVersion)),
		},
	}

	cfDeploymentReleasesAndStemcellsYaml, err := yamlMarshal(cfDeploymentReleasesAndStemcells)
	if err != nil {
		return nil, err
	}

	return append(cfDeploymentPreamble, cfDeploymentReleasesAndStemcellsYaml...), err
}
