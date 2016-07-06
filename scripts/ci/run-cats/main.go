package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/configwriter"
)

func main() {
	currentDir, _ := os.Getwd()

	missingEnvKeys := buildMissingKeyList()

	if missingEnvKeys != "" {
		fmt.Printf(`Missing required environment variables:
%s`, missingEnvKeys)
		os.Exit(1)
	}

	configWriter := configwriter.NewConfigFile(currentDir)
	configWriter.WriteConfigToFile()
	configWriter.ExportConfigFilePath()

	path, arguments := commandgenerator.GenerateCmd()
	command := exec.Command(path, arguments...)

	output, err := command.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("ERR:" + err.Error())
				os.Exit(status.ExitStatus())
			}
		} else {
			panic(err)
		}
	} else {
		fmt.Println(string(output))
	}
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
