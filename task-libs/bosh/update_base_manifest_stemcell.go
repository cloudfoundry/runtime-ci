package bosh

import (
	"fmt"
	"regexp"
)

func UpdateStemcellSection(manifestContent []byte, stemcell Stemcell) ([]byte, error) {
	if manifestContent == nil {
		return manifestContent, fmt.Errorf("manifest file has no content")
	}
	stemcellPattern := regexp.MustCompile(`(?s)stemcells:.*- alias: ([\w\-]*).*os: .* version: .*`)

	stemcellsTemplate := `stemcells:
- alias: $1
  os: %s
  version: "%s"
`
	updatedManifestContent := stemcellPattern.ReplaceAll(manifestContent,
		[]byte(fmt.Sprintf(stemcellsTemplate, stemcell.OS, stemcell.Version)))

	return updatedManifestContent, nil
}
