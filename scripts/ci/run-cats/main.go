package main

import (
	"bufio"
	"fmt"
	"io"
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

	configWriter, err := configwriter.NewConfigFile(currentDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	configWriter.WriteConfigToFile()
	configWriter.ExportConfigFilePath()

	path, arguments := commandgenerator.GenerateCmd()
	command := exec.Command(path, arguments...)

	stdOut, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stdErr, err := command.StderrPipe()
	if err != nil {
		panic(err)
	}

	go func(stdout io.ReadCloser) {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error with the Stdout scanner in attached container", err)
		}
	}(stdOut)

	go func(stderr io.ReadCloser) {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error with the Stderr scanner in attached container", err)
		}
	}(stdErr)

	err = command.Start()
	if err != nil {
		panic(err)
	}
	err = command.Wait()

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Fprintf(os.Stderr, "ERR:"+err.Error())
				os.Exit(status.ExitStatus())
			}
		} else {
			panic(err)
		}
	}

	stdOut.Close()
	stdErr.Close()
}

func buildMissingKeyList() string {
	var missingKeys string
	requiredEnvKeys := []string{
		"CF_API",
		"CF_ADMIN_USER",
		"CF_ADMIN_PASSWORD",
		"CF_APPS_DOMAIN",
	}

	for _, key := range requiredEnvKeys {
		if os.Getenv(key) == "" {
			missingKeys += key + "\n"
		}
	}

	return missingKeys
}
