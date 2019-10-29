package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/runtime-ci/tasks/detect-stemcell-bump/concourseio"
)

func main() {
	buildDir := os.Args[1]

	runner, err := concourseio.NewRunner(buildDir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.ReadCFDeploymentStemcell()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.ReadStemcell()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.DetectStemcellBump()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.WriteStemcellBumpTypeToFile()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
