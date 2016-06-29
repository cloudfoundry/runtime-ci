package configwriter

import (
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
	StaticBuildpackName    string `json:"staticfile_buildpack_name"`
	JavaBuildpackName      string `json:"java_buildpack_name"`
	RubyBuildpackName      string `json:"ruby_buildpack_name"`
	NodeJsBuildpackName    string `json:"nodejs_buildpack_name"`
	GoBuildpackName        string `json:"go_buildpack_name"`
	PythonBuildpackName    string `json:"python_buildpack_name"`
	PhpBuildpackName       string `json:"php_buildpack_name"`
	BinaryBuildpackName    string `json:"binary_buildpack_name"`
	PersistentAppHost      string `json:"persistent_app_host"`
	PersistentAppSpace     string `json:"persistent_app_space"`
	PersistentAppOrg       string `json:"persistent_app_org"`
	PersistentAppQuotaName string `json:"persistent_app_quota_name"`
}

func GenerateConfigFromEnv() config {
	skipSslValidation, _ := strconv.ParseBool(os.Getenv("SKIP_SSL_VALIDATION"))
	useHttp, _ := strconv.ParseBool(os.Getenv("USE_HTTP"))
	return config{
		Api:                    os.Getenv("CF_API"),
		AdminUser:              os.Getenv("CF_ADMIN_USER"),
		AdminPassword:          os.Getenv("CF_ADMIN_PASSWORD"),
		AppsDomain:             os.Getenv("CF_APPS_DOMAIN"),
		SkipSslValidation:      skipSslValidation,
		UseHttp:                useHttp,
		ExistingUser:           os.Getenv("EXISTING_USER"),
		UseExistingUser:        os.Getenv("EXISTING_USER") != "",
		KeepUserAtSuiteEnd:     os.Getenv("EXISTING_USER") != "",
		ExistingUserPassword:   os.Getenv("EXISTING_USER_PASSWORD"),
		StaticBuildpackName:    os.Getenv("STATICFILE_BUILDPACK_NAME"),
		JavaBuildpackName:      os.Getenv("JAVA_BUILDPACK_NAME"),
		RubyBuildpackName:      os.Getenv("RUBY_BUILDPACK_NAME"),
		NodeJsBuildpackName:    os.Getenv("NODEJS_BUILDPACK_NAME"),
		GoBuildpackName:        os.Getenv("GO_BUILDPACK_NAME"),
		PythonBuildpackName:    os.Getenv("PYTHON_BUILDPACK_NAME"),
		PhpBuildpackName:       os.Getenv("PHP_BUILDPACK_NAME"),
		BinaryBuildpackName:    os.Getenv("BINARY_BUILDPACK_NAME"),
		PersistentAppHost:      os.Getenv("PERSISTENT_APP_HOST"),
		PersistentAppSpace:     os.Getenv("PERSISTENT_APP_SPACE"),
		PersistentAppOrg:       os.Getenv("PERSISTENT_APP_ORG"),
		PersistentAppQuotaName: os.Getenv("PERSISTENT_APP_QUOTA_NAME"),
	}
}
