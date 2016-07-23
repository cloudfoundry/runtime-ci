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
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/environment"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	env := environment.New()
	errors := env.Validate()
	if !errors.Empty() {
		fmt.Fprintf(os.Stderr, "Your CATS input failed validation:\n%s\n", errors.Error())
		os.Exit(1)
	}

	configWriter, _ := configwriter.NewConfigFile(currentDir, env)
	path, arguments, _ := commandgenerator.GenerateCmd(env)

	fmt.Printf("path: %s\n", path)
	configWriter.WriteConfigToFile()
	configWriter.ExportConfigFilePath()
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
