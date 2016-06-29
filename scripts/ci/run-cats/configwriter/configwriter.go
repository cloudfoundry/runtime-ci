package configwriter

import (
	"os"
	"strconv"
)

type config struct {
	Api                  string `json:"api"`
	AdminUser            string `json:"admin_user"`
	AdminPassword        string `json:"admin_password"`
	AppsDomain           string `json:"apps_domain"`
	SkipSslValidation    bool   `json:"skip_ssl_validation"`
	UseHttp              bool   `json:"use_http"`
	ExistingUser         string `json:"existing_user"`
	UseExistingUser      bool   `json:"use_existing_user"`
	KeepUserAtSuiteEnd   bool   `json:"keep_user_at_suite_end"`
	ExistingUserPassword string `json:"existing_user_password"`
}

func GenerateConfigFromEnv() config {
	skipSslValidation, _ := strconv.ParseBool(os.Getenv("SKIP_SSL_VALIDATION"))
	useHttp, _ := strconv.ParseBool(os.Getenv("USE_HTTP"))
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
	}
}
