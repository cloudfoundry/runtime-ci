package bosh

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name           string
	Update         block
	InstanceGroups []yaml.Node `yaml:"instance_groups"`
	Releases       []Release
	Stemcells      []Stemcell
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
	tempFile, err := ioutil.TempFile(currentDir, "manifest*.yml")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	err = ioutil.WriteFile(tempFile.Name(), manifestFile, 0644)
	if err != nil {
		return err
	}

	r, err := boshCLI.Cmd("deploy", tempFile.Name(), "-d", "cf-compilation", "-n")
	if err != nil {
		io.Copy(os.Stdout, r)
		return err
	}

	return nil
}
