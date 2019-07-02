package release

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell"
)

// Release is a bosh json representation of a bosh release
type Release struct {
	Name    string
	Version string
}

func (r Release) String() string {
	return fmt.Sprint(r.Name, "/", r.Version)
}

//go:generate counterfeiter BoshCLI

type BoshCLI interface {
	Cmd(name string, args ...string) (io.Reader, error)
}

func List(boshCLI BoshCLI) ([]Release, error) {
	fmt.Println("Generating list of releases to export...")
	r, err := boshCLI.Cmd("releases", "--json")
	if err != nil {
		return nil, err
	}

	return parseReleasesOutput(r)
}

func parseReleasesOutput(r io.Reader) ([]Release, error) {
	var output struct {
		Tables []struct {
			Rows []Release
		}
	}

	err := json.NewDecoder(r).Decode(&output)
	if err != nil {
		return nil, err
	}

	var releases []Release
	for _, table := range output.Tables {
		for _, row := range table.Rows {
			if row.Name != "bosh-dns" {
				releases = append(releases, Release{Name: row.Name, Version: strings.Trim(row.Version, "*")})
			}
		}
	}

	return releases, nil
}

func Export(boshCLI BoshCLI, release Release, stemcell stemcell.Stemcell) error {
	fmt.Printf("Exporting %s for %s...\n", release.String(), stemcell.String())
	_, err := boshCLI.Cmd("export-release", "--json", release.String(), stemcell.String())
	if err != nil {
		return err
	}

	return nil
}
