package deployment

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell"
)

// Deployment is a bosh json representation of a bosh deployment
type Deployment struct {
	Name     string
	Releases []Release
	Stemcell stemcell.Stemcell
}

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

func List(boshCLI BoshCLI, stemcells []stemcell.Stemcell) ([]Deployment, error) {
	fmt.Println("Generating list of deployments...")
	r, err := boshCLI.Cmd("deployments", "--json")
	if err != nil {
		return nil, err
	}

	return parseDeploymentsOutput(r, stemcells)
}

func parseDeploymentsOutput(r io.Reader, stemcells []stemcell.Stemcell) ([]Deployment, error) {
	var output struct {
		Tables []struct {
			Rows []struct {
				Name      string
				Releases  string `json:"release_s"`
				Stemcells string `json:"stemcell_s"`
			}
		}
	}

	err := json.NewDecoder(r).Decode(&output)
	if err != nil {
		return nil, err
	}

	var outputDeployments []Deployment
	for _, table := range output.Tables {
		for _, row := range table.Rows {
			deploymentStemcells := strings.Split(row.Stemcells, "\n")

			if len(deploymentStemcells) != 1 {
				panic("only allows 1 stemcell")
			}

			stemcellInfo := strings.Split(deploymentStemcells[0], "/")

			// lookup stemcell OS from list
			os, err := getStemcellOS(stemcellInfo[0], stemcells)
			if err != nil {
				return nil, err
			}

			releases := strings.Split(row.Releases, "\n")
			deploymentReleases := []Release{}

			for _, release := range releases {
				releaseInfo := strings.Split(release, "/")

				if releaseInfo[0] != "bosh-dns" {
					deploymentReleases = append(deploymentReleases, Release{Name: releaseInfo[0], Version: releaseInfo[1]})
				}
			}

			outputDeployments = append(outputDeployments, Deployment{
				Name:     row.Name,
				Releases: deploymentReleases,
				Stemcell: stemcell.Stemcell{
					OS:      os,
					Version: stemcellInfo[1],
				},
			})
		}
	}

	return outputDeployments, nil
}

func getStemcellOS(stemcellName string, stemcells []stemcell.Stemcell) (string, error) {
	for _, stemcell := range stemcells {
		if stemcell.Name == stemcellName {
			return stemcell.OS, nil
		}
	}
	return "", fmt.Errorf("no matching stemcell name for %s", stemcellName)
}

func ExportRelease(boshCLI BoshCLI, release Release, stemcell stemcell.Stemcell, deployment Deployment) error {
	fmt.Printf("Exporting %s for %s from %s...\n", release.String(), stemcell.String(), deployment.Name)
	_, err := boshCLI.Cmd("export-release", "-d", deployment.Name, "--json", release.String(), stemcell.String())
	if err != nil {
		return err
	}

	return nil
}
