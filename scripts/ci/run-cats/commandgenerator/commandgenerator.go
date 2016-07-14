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
	"INCLUDE_LOGGING":               "logging",
	"INCLUDE_OPERATOR":              "operator",
	"INCLUDE_INTERNET_DEPENDENT":    "internet_dependent",
	"INCLUDE_SERVICES":              "services",
	"INCLUDE_ROUTE_SERVICES":        "route_services",
}

type environment interface {
	GetBoolean(string) (bool, error)
	GetString(string) string
	GetInteger(string) (int, error)
}

func GenerateCmd(env environment) (string, []string, error) {
	nodes, err := env.GetInteger("NODES")
	if err != nil {
		return "", nil, err
	}

	if nodes == 0 {
		nodes = DEFAULT_NODES
	}

	var testBinPath string
	catsPath := env.GetString("CATS_PATH")
	if catsPath != "" {
		testBinPath = filepath.Clean(catsPath + "/bin/test")
	} else {
		testBinPath = env.GetString("GOPATH") + "/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"
	}

	skipPackages, err := generateSkipPackages(env)
	if err != nil {
		return "", nil, err
	}
	skips, err := generateSkips(env)
	if err != nil {
		return "", nil, err
	}
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

func generateSkips(env environment) (string, error) {
	skip := "-skip="

	skipSso, err := env.GetBoolean("SKIP_SSO")
	if err != nil {
		return "", err
	}

	if skipSso {
		skip += "SSO|"
	}

	switch backend := env.GetString("BACKEND"); backend {
	case "diego":
		skip += "NO_DIEGO_SUPPORT"

	case "dea":
		skip += "NO_DEA_SUPPORT"

	case "":
		skip += "NO_DEA_SUPPORT|NO_DIEGO_SUPPORT"
	default:
		return "", fmt.Errorf("Invalid environment variable: 'BACKEND' was '%s', but must be 'diego', 'dea', or empty", backend)
	}

	return skip, nil
}

func generateSkipPackages(env environment) (string, error) {
	skipPackages := []string{"helpers"}
	for envVarName, packageName := range envVarToPackageMap {
		includePackage, err := env.GetBoolean(envVarName)
		if err != nil {
			return "", err
		}
		if !includePackage {
			skipPackages = append(skipPackages, packageName)
		}
	}
	sort.Strings(skipPackages)
	return "-skipPackage=" + strings.Join(skipPackages, ","), nil
}
