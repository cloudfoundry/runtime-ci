package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/runtime-ci/tasks/combine-compiled-releases/taskio"
)

func main() {
	buildDir := os.Args[1]

	//go:generate concourseTaskRunner task.ym
	runner, err := taskio.NewRunner(buildDir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.Run(tarballcombiner.NewTarballCombiner())
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// readBaseManifestReleases
	// loadPreviousTarballs
	// loadNextTarballs

	// populateCombinedReleases
	// for each release in manifest return file path
	// checkNextFirst
	// checkPrev
	// fail

	// write to outDir
}
