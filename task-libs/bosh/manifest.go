package bosh

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name           string      `yaml:",omitempty"`
	Update         block       `yaml:",omitempty"`
	InstanceGroups []yaml.Node `yaml:"instance_groups,omitempty"`
	Releases       []Release   `yaml:"releases,omitempty"`
	Stemcells      []Stemcell  `yaml:",omitempty"`
}

type block map[string]interface{}

// NewManifestFromFile creates a manifest from a yaml file.
func NewManifestFromFile(file []byte) (Manifest, error) {
	var manifest Manifest
	err := yaml.Unmarshal(file, &manifest)
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}

//go:generate counterfeiter boshCLI
type boshCLI interface {
	Cmd(name string, args ...string) (io.Reader, error)
}

func (m Manifest) Deploy(boshCLI boshCLI) error {
	m.Update = block{
		"canaries":          1,
		"max_in_flight":     1,
		"canary_watch_time": 1,
		"update_watch_time": 1,
	}
	m.Stemcells[0].Alias = "default"
	manifestFile, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	tempFile, err := os.CreateTemp(currentDir, "manifest*.yml")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name()) //nolint:errcheck

	err = os.WriteFile(tempFile.Name(), manifestFile, 0644)
	if err != nil {
		return err
	}

	_, err = boshCLI.Cmd("deploy", tempFile.Name(), "-d", m.Name, "-n", "--json")
	if err != nil {
		return err
	}

	return nil
}
