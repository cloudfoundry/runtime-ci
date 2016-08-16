package commandgenerator

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const DEFAULT_NODES = 2

var envVarToPackageMap = map[string]string{
	"INCLUDE_DIEGO_SSH":             "ssh",
	"INCLUDE_V3":                    "v3",
	"INCLUDE_DIEGO_DOCKER":          "docker",
	"INCLUDE_BACKEND_COMPATIBILITY": "backend_compatibility",
	"INCLUDE_SECURITY_GROUPS":       "security_groups",
	"INCLUDE_INTERNET_DEPENDENT":    "internet_dependent",
	"INCLUDE_SERVICES":              "services",
	"INCLUDE_ROUTE_SERVICES":        "route_services",
}

type Environment interface {
	GetSkipDiegoSSH() (string, error)
	GetSkipV3() (string, error)
	GetSkipDiegoDocker() (string, error)
	GetSkipBackendCompatibility() (string, error)
	GetSkipSecurityGroups() (string, error)
	GetSkipInternetDependent() (string, error)
	GetSkipServices() (string, error)
	GetSkipRouteServices() (string, error)
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

	skipPackages, _ := generateSkipPackages(env)

	return testBinPath, []string{
		"-r",
		"-slowSpecThreshold=120",
		"-randomizeAllSpecs",
		fmt.Sprintf("-nodes=%d", nodes),
		fmt.Sprintf("%s", skipPackages),
		"-keepGoing",
	}, nil
}

func generateSkipPackages(env Environment) (string, error) {
	skipPackages := []string{"helpers"}
	var skipPackage string

	skipPackage, _ = env.GetSkipDiegoSSH()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	skipPackage, _ = env.GetSkipV3()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	skipPackage, _ = env.GetSkipDiegoDocker()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	skipPackage, _ = env.GetSkipBackendCompatibility()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	skipPackage, _ = env.GetSkipSecurityGroups()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	skipPackage, _ = env.GetSkipInternetDependent()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	skipPackage, _ = env.GetSkipServices()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	skipPackage, _ = env.GetSkipRouteServices()
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	sort.Strings(skipPackages)
	return "-skipPackage=" + strings.Join(skipPackages, ","), nil
}
