package commandgenerator

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/validationerrors"
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

type Environment interface {
	GetSkipDiegoSSH() (string, error)
	GetSkipV3() (string, error)
	GetSkipSSO() (string, error)
	GetSkipDiegoDocker() (string, error)
	GetSkipBackendCompatibility() (string, error)
	GetSkipSecurityGroups() (string, error)
	GetSkipLogging() (string, error)
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
	var errs validationerrors.Errors

	nodes, err := env.GetNodes()
	if err != nil {
		errs.Add(err)
	}

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

	skipPackages, err := generateSkipPackages(env)
	if err != nil {
		errs.Add(err)
	}
	skips, err := generateSkips(env)
	if err != nil {
		errs.Add(err)
	}

	if !errs.Empty() {
		return "", nil, errs
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

func generateSkips(env Environment) (string, error) {
	skips := []string{}

	skipSso, err := env.GetSkipSSO()
	if err != nil {
		return "", err
	}

	if skipSso != "" {
		skips = append(skips, skipSso)
	}

	backend, err := env.GetBackend()
	if err != nil {
		return "", err
	}

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
	var err error
	errs := validationerrors.Errors{}

	if skipPackage, err = env.GetSkipDiegoSSH(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipV3(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipDiegoDocker(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipBackendCompatibility(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipSecurityGroups(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipLogging(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipOperator(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipInternetDependent(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipServices(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if skipPackage, err = env.GetSkipRouteServices(); err != nil {
		errs.Add(err)
	}
	if skipPackage != "" {
		skipPackages = append(skipPackages, skipPackage)
	}

	if !errs.Empty() {
		return "", errs
	}

	sort.Strings(skipPackages)
	return "-skipPackage=" + strings.Join(skipPackages, ","), nil
}
