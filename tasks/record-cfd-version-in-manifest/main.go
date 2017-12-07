package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
)

var majorPtr = flag.Bool("major", false, "Causes a major version bump in `manifest_version`")
var minorPtr = flag.Bool("minor", false, "Causes a minor version bump in `manifest_version`")
var patchPtr = flag.Bool("patch", false, "Causes a patch version bump in `manifest_version")

func main() {
	// TODO: Maybe make input and output file names variable?
	flag.Parse()

	// Grab the current value of the `manifest_version` using `bosh interpolate`
	manifestVersionBytes, err := boshExtract("cf-deployment.yml", "/manifest_version")
	if err != nil {
		log.Fatal(err, "\n[ Failed to extract manifest_version from cf-deployment.yml ]")
	}

	// Update the version based on flags passed to this utility
	manifestVersionString := strings.TrimSpace(string(manifestVersionBytes))
	newVersionString := updateVersion(manifestVersionString)

	// Write the new value out to a temporary ops file
	opsFile, err := ioutil.TempFile("", "manifest-version-ops")
	if err != nil {
		log.Fatal(err, "\n[ Failed to create manifest-version-ops ]")
	}

	opsFile.WriteString("---\n- type: replace\n  path: /manifest_version\n  value: " + newVersionString)
	defer os.Remove(opsFile.Name())

	// Insert the new value into the manifest and write it out using `bosh interpolate`
	newManifestBytes, err := boshIncept("cf-deployment.yml", opsFile.Name())
	if err != nil {
		log.Fatal(err, "\n[ Failed to update manifest_version in cf-deployment.yml ]")
	}

	ioutil.WriteFile("cf-deployment.yml", newManifestBytes, 0666)
}

func boshExtract(filePath, opPath string) ([]byte, error) {
	cmd := exec.Command("bosh", "int", filePath, fmt.Sprintf("--path=%s", opPath))

	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func boshIncept(manifestPath, opsFilePath string) ([]byte, error) {
	cmd := exec.Command("bosh", "int", manifestPath, "-o", opsFilePath)

	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Updates the `manifest_version` string from the cf-deployment manifest.
// `manifest_version` is expected to follow the semver spec.
// Prepends "v" to any semantic version it receives.
// Explicitly ignores any pre-release text in the semver: we expect this utility to be used only for releases.
func updateVersion(mvs string) string {
	manifestVersion, _ := version.NewVersion(mvs)
	versionSegments := manifestVersion.Segments()

	var newVersionStringSegments = make([]string, len(versionSegments))
	for i, v := range versionSegments {
		switch i {
		case 0:
			if *majorPtr {
				v++
			}
		case 1:
			if *minorPtr {
				v++
			}
		case 2:
			if *patchPtr {
				v++
			}
		}
		newVersionStringSegments[i] = strconv.Itoa(v)
	}

	return "v" + strings.Join(newVersionStringSegments, ".")
}
