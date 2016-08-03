package environment

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/validationerrors"
)

type environment struct{}

func New() *environment {
	return &environment{}
}

func (env *environment) Validate() validationerrors.Errors {
	var err error
	errs := validationerrors.Errors{}

	missingEnvKeys := buildMissingKeyList()

	if missingEnvKeys != "" {
		errs.Add(fmt.Errorf(`* Missing required environment variables:
%s`, missingEnvKeys))
	}

	if _, err = env.GetSkipDiegoSSH(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipV3(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipDiegoDocker(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipBackendCompatibility(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipSecurityGroups(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipOperator(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipInternetDependent(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipServices(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipRouteServices(); err != nil {
		errs.Add(err)
	}

	if _, err := env.GetNodes(); err != nil {
		errs.Add(err)
	}

	if _, err := env.GetSkipSSO(); err != nil {
		errs.Add(err)
	}

	if _, err := env.GetBackend(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetSkipSSLValidation(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetUseHTTP(); err != nil {
		errs.Add(err)
	}

	if _, err := env.GetIncludePrivilegedContainerSupport(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetDefaultTimeoutInSeconds(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetCFPushTimeoutInSeconds(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetLongCurlTimeoutInSeconds(); err != nil {
		errs.Add(err)
	}

	if _, err = env.GetBrokerStartTimeoutInSeconds(); err != nil {
		errs.Add(err)
	}

	return errs
}

func (e *environment) GetSkipSSLValidation() (bool, error) {
	return e.getBooleanDefaultToFalse("SKIP_SSL_VALIDATION")
}

func (e *environment) GetUseHTTP() (bool, error) {
	return e.getBooleanDefaultToFalse("USE_HTTP")
}

func (e *environment) GetIncludePrivilegedContainerSupport() (bool, error) {
	return e.getBooleanDefaultToFalse("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT")
}

func (e *environment) GetDefaultTimeoutInSeconds() (int, error) {
	return e.getInteger("DEFAULT_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetCFPushTimeoutInSeconds() (int, error) {
	return e.getInteger("CF_PUSH_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetLongCurlTimeoutInSeconds() (int, error) {
	return e.getInteger("LONG_CURL_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetBrokerStartTimeoutInSeconds() (int, error) {
	return e.getInteger("BROKER_START_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetCFAPI() string {
	return e.getString("CF_API")
}

func (e *environment) GetCFAdminUser() string {
	return e.getString("CF_ADMIN_USER")
}

func (e *environment) GetCFAdminPassword() string {
	return e.getString("CF_ADMIN_PASSWORD")
}

func (e *environment) GetCFAppsDomain() string {
	return e.getString("CF_APPS_DOMAIN")
}

func (e *environment) GetExistingUser() string {
	return e.getString("EXISTING_USER")
}

func (e *environment) UseExistingUser() bool {
	return e.getString("EXISTING_USER") != ""
}

func (e *environment) KeepUserAtSuiteEnd() bool {
	return e.getString("EXISTING_USER") != ""
}

func (e *environment) GetExistingUserPassword() string {
	return e.getString("EXISTING_USER_PASSWORD")
}

func (e *environment) GetStaticBuildpackName() string {
	return e.getString("STATICFILE_BUILDPACK_NAME")
}

func (e *environment) GetJavaBuildpackName() string {
	return e.getString("JAVA_BUILDPACK_NAME")
}

func (e *environment) GetRubyBuildpackName() string {
	return e.getString("RUBY_BUILDPACK_NAME")
}

func (e *environment) GetNodeJSBuildpackName() string {
	return e.getString("NODEJS_BUILDPACK_NAME")
}

func (e *environment) GetGoBuildpackName() string {
	return e.getString("GO_BUILDPACK_NAME")
}

func (e *environment) GetPythonBuildpackName() string {
	return e.getString("PYTHON_BUILDPACK_NAME")
}

func (e *environment) GetPHPBuildpackName() string {
	return e.getString("PHP_BUILDPACK_NAME")
}

func (e *environment) GetBinaryBuildpackName() string {
	return e.getString("BINARY_BUILDPACK_NAME")
}

func (e *environment) GetPersistentAppHost() string {
	return e.getString("PERSISTENT_APP_HOST")
}

func (e *environment) GetPersistentAppSpace() string {
	return e.getString("PERSISTENT_APP_SPACE")
}

func (e *environment) GetPersistentAppOrg() string {
	return e.getString("PERSISTENT_APP_ORG")
}

func (e *environment) GetPersistentAppQuotaName() string {
	return e.getString("PERSISTENT_APP_QUOTA_NAME")
}

func (e *environment) GetBackend() (string, error) {
	value := os.Getenv("BACKEND")
	switch value {
	case "dea", "diego", "":
		return value, nil
	default:
		return "", fmt.Errorf("* Invalid environment variable: 'BACKEND' must be 'diego', 'dea', or empty but was set to '%s'", value)
	}
}

func (e *environment) GetSkipDiegoSSH() (string, error) {
	return e.returnsSkipFlag("INCLUDE_DIEGO_SSH", "ssh", false)
}

func (e *environment) GetSkipV3() (string, error) {
	return e.returnsSkipFlag("INCLUDE_V3", "v3", false)
}

func (e *environment) GetSkipDiegoDocker() (string, error) {
	return e.returnsSkipFlag("INCLUDE_DIEGO_DOCKER", "docker", false)
}

func (e *environment) GetSkipSSO() (string, error) {
	return e.returnsSkipFlag("SKIP_SSO", "SSO", true)
}

func (e *environment) GetSkipSecurityGroups() (string, error) {
	return e.returnsSkipFlag("INCLUDE_SECURITY_GROUPS", "security_groups", false)
}

func (e *environment) GetSkipBackendCompatibility() (string, error) {
	return e.returnsSkipFlag("INCLUDE_BACKEND_COMPATIBILITY", "backend_compatibility", false)
}

func (e *environment) GetSkipOperator() (string, error) {
	return e.returnsSkipFlag("INCLUDE_OPERATOR", "operator", false)
}

func (e *environment) GetSkipInternetDependent() (string, error) {
	return e.returnsSkipFlag("INCLUDE_INTERNET_DEPENDENT", "internet_dependent", false)
}

func (e *environment) GetSkipServices() (string, error) {
	return e.returnsSkipFlag("INCLUDE_SERVICES", "services", false)
}

func (e *environment) GetSkipRouteServices() (string, error) {
	return e.returnsSkipFlag("INCLUDE_ROUTE_SERVICES", "route_services", false)
}

func (e *environment) GetCatsPath() string {
	catsPath := os.Getenv("CATS_PATH")
	if catsPath == "" {
		catsPath = e.GetGoPath() + "/src/github.com/cloudfoundry/cf-acceptance-tests"
	}
	return catsPath
}

func (e *environment) GetGoPath() string {
	return os.Getenv("GOPATH")
}

func (e *environment) GetNodes() (int, error) {
	return e.getInteger("NODES")
}

func (e *environment) getBooleanDefaultToFalse(varName string) (bool, error) {
	value := os.Getenv(varName)
	switch value {
	case "true":
		return true, nil
	case "false", "":
		return false, nil
	default:
		return false, fmt.Errorf("* Invalid environment variable: '%s' must be a boolean 'true' or 'false' but was set to '%s'", varName, value)
	}
}

func (e *environment) getBooleanDefaultToTrue(varName string) (bool, error) {
	value := os.Getenv(varName)
	switch value {
	case "true", "":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("* Invalid environment variable: '%s' must be a boolean 'true' or 'false' but was set to '%s'", varName, value)
	}
}

func (e *environment) getString(varName string) string {
	return os.Getenv(varName)
}

func (e *environment) getInteger(varName string) (int, error) {
	value := os.Getenv(varName)
	if value == "" {
		return 0, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil || intValue <= 0 {
		return 0, fmt.Errorf("* Invalid environment variable: '%s' must be an integer greater than 0 but was set to '%s'", varName, value)
	}

	return intValue, nil
}

func (e *environment) returnsSkipFlag(envVarName, skipFlag string, isSkipVariable bool) (string, error) {
	if isSkipVariable {
		value, err := e.getBooleanDefaultToTrue(envVarName)
		if err != nil {
			return "", err
		}
		if value {
			return skipFlag, nil
		}
	} else {
		value, err := e.getBooleanDefaultToFalse(envVarName)
		if err != nil {
			return "", err
		}
		if !value {
			return skipFlag, nil
		}
	}
	return "", nil
}

func buildMissingKeyList() string {
	var missingKeys string
	requiredEnvKeys := []string{
		"CF_API",
		"CF_ADMIN_USER",
		"CF_ADMIN_PASSWORD",
		"CF_APPS_DOMAIN",
	}

	for _, key := range requiredEnvKeys {
		if os.Getenv(key) == "" {
			missingKeys += "    " + key + "\n"
		}
	}

	return missingKeys
}
