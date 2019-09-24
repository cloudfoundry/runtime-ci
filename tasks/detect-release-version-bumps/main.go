package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/runtime-ci/tasks/detect-release-version-bumps/concourseio"
	"github.com/cloudfoundry/runtime-ci/tasks/detect-release-version-bumps/releasedetector"
)

func main() {
	buildDir := os.Args[1]

	runner, err := concourseio.NewRunner(buildDir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.Run(releasedetector.NewReleaseDetector())
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
