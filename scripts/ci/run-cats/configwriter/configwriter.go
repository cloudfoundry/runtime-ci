package configwriter

import (
	"encoding/json"
	"os"
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
	DefaultTimeout         int    `json:"default_timeout,omitempty"`
	CfPushTimeout          int    `json:"cf_push_timeout,omitempty"`
	LongCurlTimeout        int    `json:"long_curl_timeout,omitempty"`
	BrokerStartTimeout     int    `json:"broker_start_timeout,omitempty"`
}

func GenerateConfigFromEnv() config {
	skipSslValidation, _ := strconv.ParseBool(os.Getenv("SKIP_SSL_VALIDATION"))
	useHttp, _ := strconv.ParseBool(os.Getenv("USE_HTTP"))
	defaultTimeout, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT"))
	cfPushTimeout, _ := strconv.Atoi(os.Getenv("CF_PUSH_TIMEOUT"))
	longCurlTimeout, _ := strconv.Atoi(os.Getenv("LONG_CURL_TIMEOUT"))
	brokerStartTimeout, _ := strconv.Atoi(os.Getenv("BROKER_START_TIMEOUT"))

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
		Backend:              os.Getenv("BACKEND"),

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
	}
}

func WriteConfigToFile(destination string, config config) {
	configFile, _ := os.Create(destination + "integration_config.json")
	configJson, _ := json.Marshal(config)
	contents := []byte(configJson)

	configFile.Write(contents)
}
