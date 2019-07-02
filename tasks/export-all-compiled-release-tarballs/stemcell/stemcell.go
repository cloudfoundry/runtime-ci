package stemcell

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Stemcell is a bosh json representation of a bosh stemcell
type Stemcell struct {
	Name    string
	OS      string
	Version string
}

func (s Stemcell) String() string {
	return fmt.Sprint(s.OS, "/", s.Version)
}

//go:generate counterfeiter BoshCLI

type BoshCLI interface {
	Cmd(name string, args ...string) (io.Reader, error)
}

func List(boshCLI BoshCLI) ([]Stemcell, error) {
	fmt.Println("Generating list of stemcells...")
	r, err := boshCLI.Cmd("stemcells", "--json")
	if err != nil {
		return nil, err
	}

	return parseStemcellsOutput(r)
}

func parseStemcellsOutput(r io.Reader) ([]Stemcell, error) {
	var output struct {
		Tables []struct {
			Rows []Stemcell
		}
	}

	err := json.NewDecoder(r).Decode(&output)
	if err != nil {
		return nil, err
	}

	var stemcells []Stemcell
	for _, table := range output.Tables {
		stemcells = append(stemcells, table.Rows...)
	}

	for i, stemcell := range stemcells {
		stemcells[i].Version = strings.Trim(stemcell.Version, "*")
	}

	return stemcells, nil
}
