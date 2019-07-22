package manifest

import (
	"fmt"
	"regexp"
)

type Stemcell struct {
	OS      string
	Version string
}

func Update(manifestContent []byte, stemcell Stemcell) ([]byte, error) {
	releasePattern := regexp.MustCompile(`stemcells:
- alias: (.*)
  os: .*
	version: .*`)

	updatedManifestContent := releasePattern.ReplaceAll(manifestContent, []byte(fmt.Sprintf(`stemcells:
- alias: $1
  os: %s
	version: %s`, stemcell.OS, stemcell.Version)))

	fmt.Printf("Updated: %q\n", string(updatedManifestContent))
	return updatedManifestContent, nil
}
