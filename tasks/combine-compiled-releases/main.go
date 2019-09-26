package main

import (
	"fmt"
	"os"
)

func main() {
	buildDir := os.Args[1]

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
}
