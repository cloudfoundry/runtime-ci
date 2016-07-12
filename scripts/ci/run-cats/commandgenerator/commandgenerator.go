package commandgenerator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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
}

func GenerateCmd(env environment) (string, []string, error) {
	nodes := os.Getenv("NODES")
	var testBinPath string

	catsPath, keyExists := os.LookupEnv("CATS_PATH")
	if keyExists {
		testBinPath = filepath.Clean(catsPath + "/bin/test")
	} else {
		testBinPath = os.Getenv("GOPATH") + "/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"
	}

	skipPackages, err := generateSkipPackages(env)
	if err != nil {
		return "", nil, err
	}

	return testBinPath, []string{
		"-r",
		"-slowSpecThreshold=120",
		"-randomizeAllSpecs",
		"-nodes",
		fmt.Sprintf("%s", nodes),
		fmt.Sprintf("%s", skipPackages),
		fmt.Sprintf("%s", generateSkips()),
		"-keepGoing",
	}, nil
}

func generateSkips() string {
	skip := "-skip="

	if os.Getenv("SKIP_SSO") != "" {
		skip += "SSO|"
	}

	switch os.Getenv("BACKEND") {
	case "diego":
		skip += "NO_DIEGO_SUPPORT"

	case "dea":
		skip += "NO_DEA_SUPPORT"

	default:
		skip += "NO_DEA_SUPPORT|NO_DIEGO_SUPPORT"
	}

	return skip
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
