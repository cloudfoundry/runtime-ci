package main

import (
	"fmt"
	"os"
)

func main() {

	missingEnvKeys := buildMissingKeyList()
	fmt.Printf(`Missing required environment variables:
%s`, missingEnvKeys)

}

func buildMissingKeyList() string {
	var missingKeys string
	requiredEnvKeys := []string{
		"CF_API",
		"CF_ADMIN_USER",
		"CF_ADMIN_PASSWORD",
		"CF_APPS_DOMAIN",
		"EXISTING_USER",
		"EXISTING_USER_PASSWORD",
	}

	for _, key := range requiredEnvKeys {
		if os.Getenv(key) == "" {
			missingKeys += key + "\n"
		}
	}

	return missingKeys
}
