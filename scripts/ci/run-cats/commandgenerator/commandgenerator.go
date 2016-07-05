package commandgenerator

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateCmd() (string, []string) {
	nodes := os.Getenv("NODES")
	var testBinPath string

	catsPath, keyExists := os.LookupEnv("CATS_PATH")
	if keyExists {
		testBinPath = filepath.Clean(catsPath + "/bin/test")
	} else {
		testBinPath = "gopath/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"
	}

	return testBinPath, []string{
		"-r",
		"-slowSpecThreshold=120",
		"-randomizeAllSpecs",
		"-nodes",
		fmt.Sprintf("%s", nodes),
		fmt.Sprintf("%s", generateSkipPackages()),
		fmt.Sprintf("%s", generateSkips()),
		"-keepGoing",
	}
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

func generateSkipPackages() string {
	type envVarStruct struct {
		envKey   string
		envValue string
	}

	envVarMap := []envVarStruct{
		{"INCLUDE_DIEGO_SSH", "ssh"},
		{"INCLUDE_V3", "v3"},
		{"INCLUDE_DIEGO_DOCKER", "docker"},
		{"INCLUDE_BACKEND_COMPATIBILITY", "backend_compatibility"},
		{"INCLUDE_SECURITY_GROUPS", "security_groups"},
		{"INCLUDE_LOGGING", "logging"},
		{"INCLUDE_OPERATOR", "operator"},
		{"INCLUDE_INTERNET_DEPENDENT", "internet_dependent"},
		{"INCLUDE_SERVICES", "services"},
		{"INCLUDE_ROUTE_SERVICES", "route_services"},
	}

	skipPackages := "-skipPackage=helpers"

	for _, envVar := range envVarMap {
		envVarValue, envVarExists := os.LookupEnv(envVar.envKey)
		if !envVarExists || envVarValue != "true" {
			skipPackages += "," + envVar.envValue
		}
	}
	return skipPackages
}
