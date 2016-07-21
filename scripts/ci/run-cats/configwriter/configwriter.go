package configwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/validationerrors"
)

type config struct {
	Api                               string `json:"api"`
	AdminUser                         string `json:"admin_user"`
	AdminPassword                     string `json:"admin_password"`
	AppsDomain                        string `json:"apps_domain"`
	SkipSslValidation                 bool   `json:"skip_ssl_validation"`
	UseHttp                           bool   `json:"use_http"`
	ExistingUser                      string `json:"existing_user,omitempty"`
	ExistingUserPassword              string `json:"existing_user_password,omitempty"`
	UseExistingUser                   bool   `json:"use_existing_user"`
	KeepUserAtSuiteEnd                bool   `json:"keep_user_at_suite_end"`
	Backend                           string `json:"backend,omitempty"`
	StaticBuildpackName               string `json:"staticfile_buildpack_name,omitempty"`
	JavaBuildpackName                 string `json:"java_buildpack_name,omitempty"`
	RubyBuildpackName                 string `json:"ruby_buildpack_name,omitempty"`
	NodeJsBuildpackName               string `json:"nodejs_buildpack_name,omitempty"`
	GoBuildpackName                   string `json:"go_buildpack_name,omitempty"`
	PythonBuildpackName               string `json:"python_buildpack_name,omitempty"`
	PhpBuildpackName                  string `json:"php_buildpack_name,omitempty"`
	BinaryBuildpackName               string `json:"binary_buildpack_name,omitempty"`
	PersistentAppHost                 string `json:"persistent_app_host,omitempty"`
	PersistentAppSpace                string `json:"persistent_app_space,omitempty"`
	PersistentAppOrg                  string `json:"persistent_app_org,omitempty"`
	PersistentAppQuotaName            string `json:"persistent_app_quota_name,omitempty"`
	DefaultTimeout                    int    `json:"default_timeout,omitempty"`
	CfPushTimeout                     int    `json:"cf_push_timeout,omitempty"`
	LongCurlTimeout                   int    `json:"long_curl_timeout,omitempty"`
	BrokerStartTimeout                int    `json:"broker_start_timeout,omitempty"`
	IncludePrivilegedContainerSupport bool   `json:"include_privileged_container_support,omitempty"`
}

type configFile struct {
	Config         config
	DestinationDir string
}

type Environment interface {
	GetSkipSSLValidation() (bool, error)
	GetUseHTTP() (bool, error)
	GetIncludePrivilegedContainerSupport() (bool, error)
	GetDefaultTimeoutInSeconds() (int, error)
	GetCFPushTimeoutInSeconds() (int, error)
	GetLongCurlTimeoutInSeconds() (int, error)
	GetBrokerStartTimeoutInSeconds() (int, error)
	GetCFAPI() string
	GetCFAdminUser() string
	GetCFAdminPassword() string
	GetCFAppsDomain() string
	GetExistingUser() string
	UseExistingUser() bool
	KeepUserAtSuiteEnd() bool
	GetExistingUserPassword() string
	GetStaticBuildpackName() string
	GetJavaBuildpackName() string
	GetRubyBuildpackName() string
	GetNodeJSBuildpackName() string
	GetGoBuildpackName() string
	GetPythonBuildpackName() string
	GetPHPBuildpackName() string
	GetBinaryBuildpackName() string
	GetPersistentAppHost() string
	GetPersistentAppSpace() string
	GetPersistentAppOrg() string
	GetPersistentAppQuotaName() string
	GetBackend() (string, error)
}

func NewConfigFile(destinationDir string, env Environment) (configFile, error) {
	config, err := generateConfigFromEnv(env)
	return configFile{config, filepath.Clean(destinationDir)}, err
}

func generateConfigFromEnv(env Environment) (config, error) {
	var (
		err                                                                error
		errs                                                               validationerrors.Errors
		skipSslValidation, useHttp                                         bool
		defaultTimeout, cfPushTimeout, longCurlTimeout, brokerStartTimeout int
	)
	errs = validationerrors.Errors{}

	skipSslValidation, err = env.GetSkipSSLValidation()
	if err != nil {
		errs.Add(err)
	}

	useHttp, err = env.GetUseHTTP()
	if err != nil {
		errs.Add(err)
	}

	includePrivilegedContainerSupport, err := env.GetIncludePrivilegedContainerSupport()
	if err != nil {
		errs.Add(err)
	}

	defaultTimeout, err = env.GetDefaultTimeoutInSeconds()
	if err != nil {
		errs.Add(err)
	}

	cfPushTimeout, err = env.GetCFPushTimeoutInSeconds()
	if err != nil {
		errs.Add(err)
	}

	longCurlTimeout, err = env.GetLongCurlTimeoutInSeconds()
	if err != nil {
		errs.Add(err)
	}

	brokerStartTimeout, err = env.GetBrokerStartTimeoutInSeconds()
	if err != nil {
		errs.Add(err)
	}

	backend, err := env.GetBackend()
	if err != nil {
		errs.Add(err)
	}

	if !errs.Empty() {
		return config{}, errs
	}

	return config{
		Api:                  env.GetCFAPI(),
		AdminUser:            env.GetCFAdminUser(),
		AdminPassword:        env.GetCFAdminPassword(),
		AppsDomain:           env.GetCFAppsDomain(),
		SkipSslValidation:    skipSslValidation,
		UseHttp:              useHttp,
		ExistingUser:         env.GetExistingUser(),
		UseExistingUser:      env.UseExistingUser(),
		KeepUserAtSuiteEnd:   env.KeepUserAtSuiteEnd(),
		ExistingUserPassword: env.GetExistingUserPassword(),
		Backend:              backend,

		StaticBuildpackName: env.GetStaticBuildpackName(),
		JavaBuildpackName:   env.GetJavaBuildpackName(),
		RubyBuildpackName:   env.GetRubyBuildpackName(),
		NodeJsBuildpackName: env.GetNodeJSBuildpackName(),
		GoBuildpackName:     env.GetGoBuildpackName(),
		PythonBuildpackName: env.GetPythonBuildpackName(),
		PhpBuildpackName:    env.GetPHPBuildpackName(),
		BinaryBuildpackName: env.GetBinaryBuildpackName(),

		PersistentAppHost:      env.GetPersistentAppHost(),
		PersistentAppSpace:     env.GetPersistentAppSpace(),
		PersistentAppOrg:       env.GetPersistentAppOrg(),
		PersistentAppQuotaName: env.GetPersistentAppQuotaName(),

		DefaultTimeout:     defaultTimeout,
		CfPushTimeout:      cfPushTimeout,
		LongCurlTimeout:    longCurlTimeout,
		BrokerStartTimeout: brokerStartTimeout,

		IncludePrivilegedContainerSupport: includePrivilegedContainerSupport,
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
