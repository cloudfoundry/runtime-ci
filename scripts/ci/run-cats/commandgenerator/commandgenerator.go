package commandgenerator

import (
	"fmt"
	"os"
)

func GenerateCmd() string {
	nodes := os.Getenv("NODES")

	return fmt.Sprintf(
		"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes %s %s %s -keepGoing",
		nodes,
		generateSkipPackages(),
		generateSkips(),
	)
}

func generateSkips() string {
	skip := "-skip="

	if os.Getenv("SKIP_SSO") != "" {
		skip += "SSO|"
	}

	switch os.Getenv("BACKEND") {
	case "diego":
		skip += "NO_DEA_SUPPORT"

	case "dea":
		skip += "NO_DIEGO_SUPPORT"

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
		if os.Getenv(envVar.envKey) != "" {
			skipPackages += "," + envVar.envValue
		}
	}
	return skipPackages
}
