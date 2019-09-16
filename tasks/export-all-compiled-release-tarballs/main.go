package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/command"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/deployment"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell"
)

func main() {
	boshCLI := new(command.BoshCLI)

	stemcells, err := stemcell.List(boshCLI)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deployments, err := deployment.List(boshCLI, stemcells)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var errOccured bool
	var wg sync.WaitGroup

	for _, boshDeployment := range deployments {
		boshStemcell := boshDeployment.Stemcell

		for _, boshRelease := range boshDeployment.Releases {
			wg.Add(1)

			go func(boshRelease deployment.Release, boshStemcell stemcell.Stemcell, boshDeployment deployment.Deployment, wg *sync.WaitGroup) {
				err := deployment.ExportRelease(boshCLI, boshRelease, boshStemcell, boshDeployment)
				if err != nil {
					errOccured = true
					fmt.Printf("Failed to export release: %s\n", err)
				}
				wg.Done()
			}(boshRelease, boshStemcell, boshDeployment, &wg)
		}
	}

	wg.Wait()

	if errOccured {
		os.Exit(1)
	}
}
