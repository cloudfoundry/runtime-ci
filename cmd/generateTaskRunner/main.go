package generateTaskRunner

import (
	"log"

	"github.com/cloudfoundry/runtime-ci/cmd/generateTaskRunner/generate"
	"github.com/spf13/pflag"
)

var (
	outputDir     string
	outputPackage string
	taskYAMLPath  string
)

func init() {
	pflag.StringVarP(&outputDir, "output-directory", "d", "", "the directory of the generated content. if not specified, use the package name relative to the current directory")
	pflag.StringVarP(&outputPackage, "output-package", "o", "runner", "the package name of the generated content")
	pflag.StringVarP(&taskYAMLPath, "task-yaml-path", "p", "task.yml", "the path to the task YAML file")

	pflag.Parse()
}

func main() {
	// taskYAML object
	// //Inputs []strings
	// //Outputs []strings
	taskPackage, err := generate.NewPackage(outputPackage, taskYAMLPath)
	if err != nil {
		log.Fatal(err)
	}

	if outputDir != "" {
		taskPackage.SetOutputDir(outputDir)
	}

	// Generate Content
	err := taskPackage.Write()
	if err != nil {
		log.Fatal(err)
	}
	// Write
	//  // Create ginkgo suite test
	// 	// Write runner.go
	// 	// Write runner_test.go
}
