package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var buildDir string
	flag.StringVar(&buildDir, "build-dir", "", "path to the build directory")

	var inputDir string
	flag.StringVar(&inputDir, "input-dir", "", "path to the original cf-deployment")

	var outputDir string
	flag.StringVar(&outputDir, "output-dir", "", "path to the updated cf-deployment")

	flag.Parse()

	if buildDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: build-dir")
		os.Exit(1)
	}

	if inputDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: input-dir")
		os.Exit(1)
	}

	if outputDir == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: output-dir")
		os.Exit(1)
	}
}
