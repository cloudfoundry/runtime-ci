package environment

import (
	"fmt"
	"os"
	"strconv"
)

type environment struct{}

func New() *environment {
	return &environment{}
}

func (e *environment) GetSkipSSLValidation() (bool, error) {
	return e.GetBoolean("SKIP_SSL_VALIDATION")
}

func (e *environment) GetUseHTTP() (bool, error) {
	return e.GetBoolean("USE_HTTP")
}

func (e *environment) GetIncludePrivilegedContainerSupport() (bool, error) {
	return e.GetBoolean("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT")
}

func (e *environment) GetDefaultTimeoutInSeconds() (int, error) {
	return e.GetInteger("DEFAULT_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetCFPushTimeoutInSeconds() (int, error) {
	return e.GetInteger("CF_PUSH_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetLongCurlTimeoutInSeconds() (int, error) {
	return e.GetInteger("LONG_CURL_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetBrokerStartTimeoutInSeconds() (int, error) {
	return e.GetInteger("BROKER_START_TIMEOUT_IN_SECONDS")
}

func (e *environment) GetCFAPI() string {
	return e.GetString("CF_API")
}

func (e *environment) GetCFAdminUser() string {
	return e.GetString("CF_ADMIN_USER")
}

func (e *environment) GetCFAdminPassword() string {
	return e.GetString("CF_ADMIN_PASSWORD")
}

func (e *environment) GetCFAppsDomain() string {
	return e.GetString("CF_APPS_DOMAIN")
}

func (e *environment) GetExistingUser() string {
	return e.GetString("EXISTING_USER")
}

func (e *environment) UseExistingUser() bool {
	return e.GetString("EXISTING_USER") != ""
}

func (e *environment) KeepUserAtSuiteEnd() bool {
	return e.GetString("EXISTING_USER") != ""
}

func (e *environment) GetExistingUserPassword() string {
	return e.GetString("EXISTING_USER_PASSWORD")
}

func (e *environment) GetStaticBuildpackName() string {
	return e.GetString("STATICFILE_BUILDPACK_NAME")
}

func (e *environment) GetJavaBuildpackName() string {
	return e.GetString("JAVA_BUILDPACK_NAME")
}

func (e *environment) GetRubyBuildpackName() string {
	return e.GetString("RUBY_BUILDPACK_NAME")
}

func (e *environment) GetNodeJSBuildpackName() string {
	return e.GetString("NODEJS_BUILDPACK_NAME")
}

func (e *environment) GetGoBuildpackName() string {
	return e.GetString("GO_BUILDPACK_NAME")
}

func (e *environment) GetPythonBuildpackName() string {
	return e.GetString("PYTHON_BUILDPACK_NAME")
}

func (e *environment) GetPHPBuildpackName() string {
	return e.GetString("PHP_BUILDPACK_NAME")
}

func (e *environment) GetBinaryBuildpackName() string {
	return e.GetString("BINARY_BUILDPACK_NAME")
}

func (e *environment) GetPersistentAppHost() string {
	return e.GetString("PERSISTENT_APP_HOST")
}

func (e *environment) GetPersistentAppSpace() string {
	return e.GetString("PERSISTENT_APP_SPACE")
}

func (e *environment) GetPersistentAppOrg() string {
	return e.GetString("PERSISTENT_APP_ORG")
}

func (e *environment) GetPersistentAppQuotaName() string {
	return e.GetString("PERSISTENT_APP_QUOTA_NAME")
}

func (e *environment) GetBoolean(varName string) (bool, error) {
	switch os.Getenv(varName) {
	case "true":
		return true, nil
	case "false", "":
		return false, nil
	default:
		return false, fmt.Errorf("Invalid environment variable: '%s' must be a boolean 'true' or 'false'", varName)
	}
}

func (e *environment) GetBooleanDefaultToTrue(varName string) (bool, error) {
	switch os.Getenv(varName) {
	case "true", "":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("Invalid environment variable: '%s' must be a boolean 'true' or 'false'", varName)
	}
}

func (e *environment) GetBackend() (string, error) {
	value := os.Getenv("BACKEND")
	switch value {
	case "dea", "diego", "":
		return value, nil
	default:
		return "", fmt.Errorf("Invalid environment variable: 'BACKEND' was '%s', but must be 'diego', 'dea', or empty", value)
	}
}

func (e *environment) GetString(varName string) string {
	return os.Getenv(varName)
}

func (e *environment) GetInteger(varName string) (int, error) {
	value := os.Getenv(varName)
	if value == "" {
		return 0, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil || intValue <= 0 {
		return 0, fmt.Errorf("Invalid environment variable: '%s' must be an integer greater than 0", varName)
	}

	return intValue, nil
}
