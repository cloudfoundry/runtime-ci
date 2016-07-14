package configwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	DefaultTimeout         int    `json:"default_timeout,omitempty"`
	CfPushTimeout          int    `json:"cf_push_timeout,omitempty"`
	LongCurlTimeout        int    `json:"long_curl_timeout,omitempty"`
	BrokerStartTimeout     int    `json:"broker_start_timeout,omitempty"`
}

type configFile struct {
	Config         config
	DestinationDir string
}

type Environment interface {
	GetBoolean(string) (bool, error)
	GetString(string) string
	GetInteger(string) (int, error)
	GetBackend() (string, error)
}

func NewConfigFile(destinationDir string, env Environment) (configFile, error) {
	config, err := generateConfigFromEnv(env)
	return configFile{config, filepath.Clean(destinationDir)}, err
}

func generateConfigFromEnv(env Environment) (config, error) {
	var (
		err                                                                error
		skipSslValidation, useHttp                                         bool
		defaultTimeout, cfPushTimeout, longCurlTimeout, brokerStartTimeout int
	)

	skipSslValidation, err = env.GetBoolean("SKIP_SSL_VALIDATION")
	if err != nil {
		return config{}, err
	}
	useHttp, err = env.GetBoolean("USE_HTTP")
	if err != nil {
		return config{}, err
	}

	defaultTimeout, err = env.GetInteger("DEFAULT_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}
	cfPushTimeout, err = env.GetInteger("CF_PUSH_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}
	longCurlTimeout, err = env.GetInteger("LONG_CURL_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}
	brokerStartTimeout, err = env.GetInteger("BROKER_START_TIMEOUT_IN_SECONDS")
	if err != nil {
		return config{}, err
	}
	backend, err := env.GetBackend()
	if err != nil {
		return config{}, err
	}

	return config{
		Api:                  env.GetString("CF_API"),
		AdminUser:            env.GetString("CF_ADMIN_USER"),
		AdminPassword:        env.GetString("CF_ADMIN_PASSWORD"),
		AppsDomain:           env.GetString("CF_APPS_DOMAIN"),
		SkipSslValidation:    skipSslValidation,
		UseHttp:              useHttp,
		ExistingUser:         env.GetString("EXISTING_USER"),
		UseExistingUser:      env.GetString("EXISTING_USER") != "",
		KeepUserAtSuiteEnd:   env.GetString("EXISTING_USER") != "",
		ExistingUserPassword: env.GetString("EXISTING_USER_PASSWORD"),
		Backend:              backend,

		StaticBuildpackName: env.GetString("STATICFILE_BUILDPACK_NAME"),
		JavaBuildpackName:   env.GetString("JAVA_BUILDPACK_NAME"),
		RubyBuildpackName:   env.GetString("RUBY_BUILDPACK_NAME"),
		NodeJsBuildpackName: env.GetString("NODEJS_BUILDPACK_NAME"),
		GoBuildpackName:     env.GetString("GO_BUILDPACK_NAME"),
		PythonBuildpackName: env.GetString("PYTHON_BUILDPACK_NAME"),
		PhpBuildpackName:    env.GetString("PHP_BUILDPACK_NAME"),
		BinaryBuildpackName: env.GetString("BINARY_BUILDPACK_NAME"),

		PersistentAppHost:      env.GetString("PERSISTENT_APP_HOST"),
		PersistentAppSpace:     env.GetString("PERSISTENT_APP_SPACE"),
		PersistentAppOrg:       env.GetString("PERSISTENT_APP_ORG"),
		PersistentAppQuotaName: env.GetString("PERSISTENT_APP_QUOTA_NAME"),

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
