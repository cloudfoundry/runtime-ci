package commandgenerator

import (
	"fmt"
	"path/filepath"
)

const DEFAULT_NODES = 2

type Environment interface {
	GetBackend() (string, error)
	GetCatsPath() string
	GetNodes() (int, error)
	GetGoPath() string
}

func GenerateCmd(env Environment) (string, []string, error) {
	nodes, _ := env.GetNodes()

	if nodes == 0 {
		nodes = DEFAULT_NODES
	}

	var testBinPath string
	catsPath := env.GetCatsPath()
	if catsPath != "" {
		testBinPath = filepath.Clean(catsPath + "/bin/test")
	} else {
		testBinPath = env.GetGoPath() + "/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"
	}

	skipPackages := "-skipPackage=helpers"

	return testBinPath, []string{
		"-r",
		"-slowSpecThreshold=120",
		"-randomizeAllSpecs",
		fmt.Sprintf("-nodes=%d", nodes),
		fmt.Sprintf("%s", skipPackages),
		"-keepGoing",
	}, nil
}
