package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/compiledrelease"
	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/concourseio"
	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/manifest"
	"github.com/spf13/pflag"
)

func main() {
	buildDir := pflag.Arg(0)
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

	err = runner.UpdateManifest(manifest.Update)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = runner.UpdateStemcell(
		compiledrelease.NewOpsfileUpdater(
			runner.In.CompiledReleasesDir,
			filepath.Join(runner.Out.UpdatedCFDeploymentDir, "operations", "use-compiled-releases.yml"),
		),
	)
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
