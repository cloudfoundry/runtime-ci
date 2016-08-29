package configwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type config struct {
	Api                               string `json:"api"`
	AdminUser                         string `json:"admin_user"`
	AdminPassword                     string `json:"admin_password"`
	AppsDomain                        string `json:"apps_domain"`
	SkipSslValidation                 bool   `json:"skip_ssl_validation"`
	IncludeSSO                        bool   `json:"include_sso"`
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
	AsyncServiceOperationTimeout      int    `json:"async_service_operation_timeout,omitempty"`
	IncludePrivilegedContainerSupport bool   `json:"include_privileged_container_support,omitempty"`
	IncludeApps                       bool   `json:"include_apps"`
	IncludeDiegoSSH                   bool   `json:"include_ssh"`
	IncludeV3                         bool   `json:"include_v3"`
	IncludeDiegoDocker                bool   `json:"include_docker"`
	IncludeSecurityGroups             bool   `json:"include_security_groups"`
	IncludeBackendCompatibility       bool   `json:"include_backend_compatibility"`
	IncludeInternetDependent          bool   `json:"include_internet_dependent"`
	IncludeServices                   bool   `json:"include_services"`
	IncludeRouteServices              bool   `json:"include_route_services"`
	IncludeRouting                    bool   `json:"include_routing"`
	IncludeTasks                      bool   `json:"include_tasks"`
	IncludeDetect                     bool   `json:"include_detect"`
}

type configFile struct {
	Config         config
	DestinationDir string
}

type Environment interface {
	GetSkipSSLValidation() (bool, error)
	GetIncludeSSO() (bool, error)
	GetUseHTTP() (bool, error)
	GetIncludePrivilegedContainerSupport() (bool, error)
	GetDefaultTimeoutInSeconds() (int, error)
	GetCFPushTimeoutInSeconds() (int, error)
	GetLongCurlTimeoutInSeconds() (int, error)
	GetBrokerStartTimeoutInSeconds() (int, error)
	GetAsyncServiceOperationTimeoutInSeconds() (int, error)
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
	GetIncludeApps() (bool, error)
	GetIncludeDiegoSSH() (bool, error)
	GetIncludeV3() (bool, error)
	GetIncludeDiegoDocker() (bool, error)
	GetIncludeSecurityGroups() (bool, error)
	GetIncludeBackendCompatibility() (bool, error)
	GetIncludeInternetDependent() (bool, error)
	GetIncludeServices() (bool, error)
	GetIncludeRouteServices() (bool, error)
	GetIncludeRouting() (bool, error)
	GetIncludeTasks() (bool, error)
	GetIncludeDetect() (bool, error)
}

func NewConfigFile(destinationDir string, env Environment) (configFile, error) {
	config, err := generateConfigFromEnv(env)
	return configFile{config, filepath.Clean(destinationDir)}, err
}

func generateConfigFromEnv(env Environment) (config, error) {
	var (
		skipSslValidation, includeSSO, useHttp                                         bool
		includeDiegoSSH, includeV3, includeDiegoDocker, includeSecurityGroups          bool
		includeBackendCompatibility, includeInternetDependent, includeServices         bool
		includeRouteServices, includeRouting, includeDetect, includeApps, includeTasks bool
		defaultTimeout, cfPushTimeout, longCurlTimeout, brokerStartTimeout             int
		asyncServiceOperationTimeout                                                   int
	)

	skipSslValidation, _ = env.GetSkipSSLValidation()

	includeSSO, _ = env.GetIncludeSSO()

	useHttp, _ = env.GetUseHTTP()

	includePrivilegedContainerSupport, _ := env.GetIncludePrivilegedContainerSupport()

	defaultTimeout, _ = env.GetDefaultTimeoutInSeconds()

	cfPushTimeout, _ = env.GetCFPushTimeoutInSeconds()

	longCurlTimeout, _ = env.GetLongCurlTimeoutInSeconds()

	brokerStartTimeout, _ = env.GetBrokerStartTimeoutInSeconds()

	asyncServiceOperationTimeout, _ = env.GetAsyncServiceOperationTimeoutInSeconds()

	includeApps, _ = env.GetIncludeApps()
	includeDiegoSSH, _ = env.GetIncludeDiegoSSH()
	includeV3, _ = env.GetIncludeV3()
	includeDiegoDocker, _ = env.GetIncludeDiegoDocker()
	includeSecurityGroups, _ = env.GetIncludeSecurityGroups()
	includeBackendCompatibility, _ = env.GetIncludeBackendCompatibility()
	includeInternetDependent, _ = env.GetIncludeInternetDependent()
	includeServices, _ = env.GetIncludeServices()
	includeRouteServices, _ = env.GetIncludeRouteServices()
	includeRouting, _ = env.GetIncludeRouting()
	includeTasks, _ = env.GetIncludeTasks()
	includeDetect, _ = env.GetIncludeDetect()

	backend, _ := env.GetBackend()

	return config{
		Api:                  env.GetCFAPI(),
		AdminUser:            env.GetCFAdminUser(),
		AdminPassword:        env.GetCFAdminPassword(),
		AppsDomain:           env.GetCFAppsDomain(),
		SkipSslValidation:    skipSslValidation,
		IncludeSSO:           includeSSO,
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

		DefaultTimeout:               defaultTimeout,
		CfPushTimeout:                cfPushTimeout,
		LongCurlTimeout:              longCurlTimeout,
		BrokerStartTimeout:           brokerStartTimeout,
		AsyncServiceOperationTimeout: asyncServiceOperationTimeout,

		IncludePrivilegedContainerSupport: includePrivilegedContainerSupport,
		IncludeApps:                       includeApps,
		IncludeDiegoSSH:                   includeDiegoSSH,
		IncludeV3:                         includeV3,
		IncludeDiegoDocker:                includeDiegoDocker,
		IncludeSecurityGroups:             includeSecurityGroups,
		IncludeBackendCompatibility:       includeBackendCompatibility,
		IncludeInternetDependent:          includeInternetDependent,
		IncludeServices:                   includeServices,
		IncludeRouteServices:              includeRouteServices,
		IncludeRouting:                    includeRouting,
		IncludeTasks:                      includeTasks,
		IncludeDetect:                     includeDetect,
	}, nil
}

func (configFile *configFile) WriteConfigToFile() (*os.File, error) {
	integrationFilePath := configFile.DestinationDir + "/integration_config.json"
	integrationConfigFile, err := os.Create(integrationFilePath)
	if err != nil {
		panic("Unable to create file:" + integrationFilePath)
	}

	configJson, err := json.MarshalIndent(configFile.Config, "", "\t")
	if err != nil {
		panic("Unable to marshal json to " + integrationFilePath)
	}

	contents := []byte(configJson)

	_, err = integrationConfigFile.Write(contents)
	if err != nil {
		panic("Unable to write contents to file" + integrationFilePath)
	}

	return integrationConfigFile, nil
}

func (configFile *configFile) ExportConfigFilePath() {
	os.Setenv("CONFIG", fmt.Sprintf("%s/integration_config.json", configFile.DestinationDir))
}
