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
	if manifestContent == nil {
		return manifestContent, fmt.Errorf("manifest file has no content")
	}
	releasePattern := regexp.MustCompile(`(?s)stemcells:.*- alias: ([\w\-]*).*os: .* version: .*`)

	stemcellsTemplate := `stemcells:
- alias: $1
  os: %s
  version: %s
`
	updatedManifestContent := releasePattern.ReplaceAll(manifestContent,
		[]byte(fmt.Sprintf(stemcellsTemplate, stemcell.OS, stemcell.Version)))

	return updatedManifestContent, nil
}
