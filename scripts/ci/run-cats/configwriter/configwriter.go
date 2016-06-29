package configwriter

import "os"

type config struct {
	Api       string
	AdminUser string
}

func GenerateConfigFromEnv() config {
	return config{
		Api:       "api.example.com",
		AdminUser: os.Getenv("CF_ADMIN_USER"),
	}
}

func GenerateConfig(api string, adminUser string) config {

	return config{
		Api:       api,
		AdminUser: adminUser,
	}
}
