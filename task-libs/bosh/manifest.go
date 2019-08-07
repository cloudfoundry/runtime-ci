package bosh

import "io"

type Manifest struct {
	Name      string
	Releases  []Release
	Stemcells []Stemcell
}

// NewManifestFromFile creates a manifest from a yaml file.
func NewManifestFromFile(file []byte) (Manifest, error) {
	return Manifest{}, nil
}

//go:generate counterfeiter BoshCLI
type BoshCLI interface {
	Cmd(name string, args ...string) (io.Reader, error)
}

func Deploy(boshCLI boshCLI) error {
	return nil
}
