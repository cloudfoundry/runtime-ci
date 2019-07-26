package main

import (
	"fmt"
	"os"

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

	// runner.UpdateCompiledReleases()
	// if err != nil {
	// 	fmt.Print(err)
	// 	os.Exit(1)
	// }

	// runner.WriteCommitMessage()
	// if err != nil {
	// 	fmt.Print(err)
	// 	os.Exit(1)
	// }
}
