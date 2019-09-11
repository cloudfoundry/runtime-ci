package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/command"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/release"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell"
)

func main() {
	boshCLI := new(command.BoshCLI)
	releases, err := release.List(boshCLI)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stemcells, err := stemcell.List(boshCLI)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(stemcells) != 1 {
		fmt.Printf("Unable to determine which stemcell to use. Found %d.\n", len(stemcells))
		os.Exit(1)
	}

	var errOccured bool
	var wg sync.WaitGroup

	for _, boshRelease := range releases {
		wg.Add(1)
		go func(boshRelease release.Release, wg sync.WaitGroup) {
			err := release.Export(boshCLI, boshRelease, stemcells[0])
			if err != nil {
				errOccured = true
				fmt.Printf("Failed to export release: %s\n", err)
			}
			wg.Done()
		}(boshRelease, wg)
	}

	wg.Wait()

	if errOccured {
		os.Exit(1)
	}
}
