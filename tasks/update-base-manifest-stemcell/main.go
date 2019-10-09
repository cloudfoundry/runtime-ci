package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	"github.com/cloudfoundry/runtime-ci/tasks/update-base-manifest-stemcell/concourseio"
)

func main() {
	buildDir := os.Args[1]
	runner, err := concourseio.NewRunner(buildDir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.ReadStemcell()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.UpdateManifest(bosh.UpdateStemcellSection)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	commitMessagePath := filepath.Join(buildDir, "commit-message.txt")

	err = runner.WriteCommitMessage(commitMessagePath)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
