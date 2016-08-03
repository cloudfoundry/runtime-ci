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
	"INCLUDE_OPERATOR":              "operator",
	"INCLUDE_INTERNET_DEPENDENT":    "internet_dependent",
	"INCLUDE_SERVICES":              "services",
	"INCLUDE_ROUTE_SERVICES":        "route_services",
}

type Environment interface {
	GetSkipDiegoSSH() (string, error)
	GetSkipV3() (string, error)
	GetSkipSSO() (string, error)
	GetSkipDiegoDocker() (string, error)
	GetSkipBackendCompatibility() (string, error)
	GetSkipSecurityGroups() (string, error)
	GetSkipOperator() (string, error)
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
	skips, _ := generateSkips(env)

	return testBinPath, []string{
		"-r",
		"-slowSpecThreshold=120",
		"-randomizeAllSpecs",
		fmt.Sprintf("-nodes=%d", nodes),
		fmt.Sprintf("%s", skipPackages),
		fmt.Sprintf("%s", skips),
		"-keepGoing",
	}, nil
}

func generateSkips(env Environment) (string, error) {
	skips := []string{}

	skipSso, _ := env.GetSkipSSO()

	if skipSso != "" {
		skips = append(skips, skipSso)
	}

	backend, _ := env.GetBackend()

	switch backend {
	case "diego":
		skips = append(skips, "NO_DIEGO_SUPPORT")
	case "dea":
		skips = append(skips, "NO_DEA_SUPPORT")
	case "":
		skips = append(skips, "NO_DEA_SUPPORT|NO_DIEGO_SUPPORT")
	}

	return "-skip=" + strings.Join(skips, "|"), nil
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

	skipPackage, _ = env.GetSkipOperator()
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
