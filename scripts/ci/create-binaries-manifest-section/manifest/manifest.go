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

func UpdateReleasesAndStemcells(releases []string, buildDir string, cfDeploymentManifest []byte) ([]byte, error) {
	r := regexp.MustCompile(`(?m:^releases:$)`)
	cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
	cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]

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
