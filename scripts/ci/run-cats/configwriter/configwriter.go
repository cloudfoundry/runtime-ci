package configwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type config struct {
	Api                    string `json:"api"`
	AdminUser              string `json:"admin_user"`
	AdminPassword          string `json:"admin_password"`
	AppsDomain             string `json:"apps_domain"`
	SkipSslValidation      bool   `json:"skip_ssl_validation"`
	UseHttp                bool   `json:"use_http"`
	ExistingUser           string `json:"existing_user"`
	UseExistingUser        bool   `json:"use_existing_user"`
	KeepUserAtSuiteEnd     bool   `json:"keep_user_at_suite_end"`
	ExistingUserPassword   string `json:"existing_user_password"`
	Backend                string `json:"backend,omitempty"`
	StaticBuildpackName    string `json:"staticfile_buildpack_name,omitempty"`
	JavaBuildpackName      string `json:"java_buildpack_name,omitempty"`
	RubyBuildpackName      string `json:"ruby_buildpack_name,omitempty"`
	NodeJsBuildpackName    string `json:"nodejs_buildpack_name,omitempty"`
	GoBuildpackName        string `json:"go_buildpack_name,omitempty"`
	PythonBuildpackName    string `json:"python_buildpack_name,omitempty"`
	PhpBuildpackName       string `json:"php_buildpack_name,omitempty"`
	BinaryBuildpackName    string `json:"binary_buildpack_name,omitempty"`
	PersistentAppHost      string `json:"persistent_app_host,omitempty"`
	PersistentAppSpace     string `json:"persistent_app_space,omitempty"`
	PersistentAppOrg       string `json:"persistent_app_org,omitempty"`
	PersistentAppQuotaName string `json:"persistent_app_quota_name,omitempty"`
	DefaultTimeout         *int   `json:"default_timeout,omitempty"`
	CfPushTimeout          *int   `json:"cf_push_timeout,omitempty"`
	LongCurlTimeout        *int   `json:"long_curl_timeout,omitempty"`
	BrokerStartTimeout     *int   `json:"broker_start_timeout,omitempty"`
}

type configFile struct {
	Config         config
	DestinationDir string
}

type Environment interface {
	GetBoolean(string) (bool, error)
}

func NewConfigFile(destinationDir string, env Environment) (configFile, error) {
	config, err := generateConfigFromEnv(env)
	return configFile{config, filepath.Clean(destinationDir)}, err
}

func getTimeoutIfPresent(envKey string) (*int, error) {
	if os.Getenv(envKey) == "" {
		return nil, nil
	}
	timeout, err := strconv.Atoi(os.Getenv(envKey))
	if err != nil || timeout <= 0 {
		return nil, fmt.Errorf("Invalid env var '%s' only allows positive integers", envKey)
	}
	return &timeout, err
}

func generateConfigFromEnv(env Environment) (config, error) {
	var (
		err                                                                error
		skipSslValidation, useHttp                                         bool
		defaultTimeout, cfPushTimeout, longCurlTimeout, brokerStartTimeout *int
	)

	skipSslValidation, err = env.GetBoolean("SKIP_SSL_VALIDATION")
	if err != nil {
		return config{}, err
	}
	useHttp, err = env.GetBoolean("USE_HTTP")
	if err != nil {
		return config{}, err
	}

	defaultTimeout, err = getTimeoutIfPresent("DEFAULT_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}
	cfPushTimeout, err = getTimeoutIfPresent("CF_PUSH_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}
	longCurlTimeout, err = getTimeoutIfPresent("LONG_CURL_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}
	brokerStartTimeout, err = getTimeoutIfPresent("BROKER_START_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}

	backend := os.Getenv("BACKEND")
	if backend != "" && backend != "dea" && backend != "diego" {
		return config{}, fmt.Errorf("Invalid env var 'BACKEND' only accepts 'dea' or 'diego'")
	}

	return config{
		Api:                  os.Getenv("CF_API"),
		AdminUser:            os.Getenv("CF_ADMIN_USER"),
		AdminPassword:        os.Getenv("CF_ADMIN_PASSWORD"),
		AppsDomain:           os.Getenv("CF_APPS_DOMAIN"),
		SkipSslValidation:    skipSslValidation,
		UseHttp:              useHttp,
		ExistingUser:         os.Getenv("EXISTING_USER"),
		UseExistingUser:      os.Getenv("EXISTING_USER") != "",
		KeepUserAtSuiteEnd:   os.Getenv("EXISTING_USER") != "",
		ExistingUserPassword: os.Getenv("EXISTING_USER_PASSWORD"),
		Backend:              backend,

		StaticBuildpackName: os.Getenv("STATICFILE_BUILDPACK_NAME"),
		JavaBuildpackName:   os.Getenv("JAVA_BUILDPACK_NAME"),
		RubyBuildpackName:   os.Getenv("RUBY_BUILDPACK_NAME"),
		NodeJsBuildpackName: os.Getenv("NODEJS_BUILDPACK_NAME"),
		GoBuildpackName:     os.Getenv("GO_BUILDPACK_NAME"),
		PythonBuildpackName: os.Getenv("PYTHON_BUILDPACK_NAME"),
		PhpBuildpackName:    os.Getenv("PHP_BUILDPACK_NAME"),
		BinaryBuildpackName: os.Getenv("BINARY_BUILDPACK_NAME"),

		PersistentAppHost:      os.Getenv("PERSISTENT_APP_HOST"),
		PersistentAppSpace:     os.Getenv("PERSISTENT_APP_SPACE"),
		PersistentAppOrg:       os.Getenv("PERSISTENT_APP_ORG"),
		PersistentAppQuotaName: os.Getenv("PERSISTENT_APP_QUOTA_NAME"),

		DefaultTimeout:     defaultTimeout,
		CfPushTimeout:      cfPushTimeout,
		LongCurlTimeout:    longCurlTimeout,
		BrokerStartTimeout: brokerStartTimeout,
	}, err
}

func (configFile *configFile) WriteConfigToFile() (*os.File, error) {
	integrationConfigFile, err := os.Create(configFile.DestinationDir + "/integration_config.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to write integration_config.json to %s", configFile.DestinationDir)
	}

	configJson, err := json.MarshalIndent(configFile.Config, "", "\t")
	if err != nil {
		return nil, err
	}

	contents := []byte(configJson)

	_, err = integrationConfigFile.Write(contents)
	if err != nil {
		return nil, err
	}

	return integrationConfigFile, nil
}

func (configFile *configFile) ExportConfigFilePath() {
	os.Setenv("CONFIG", fmt.Sprintf("%s/integration_config.json", configFile.DestinationDir))
}
