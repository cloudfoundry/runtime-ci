package runner

import (
	"encoding/json"
	"fmt"
	"stemcell-version-bump/resource"
)

type CheckResponse []resource.Version

//go:generate counterfeiter . Getter
type Getter interface {
	Get(bucketName string, objectPath string) ([]byte, error)
}

func Check(config resource.Config, getter Getter) (string, error) {
	content, err := getter.Get(config.Source.BucketName, config.Source.FileName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version info from bucket/file (%s, %s): %w", config.Source.BucketName, config.Source.FileName, err)
	}

	var currentVersion resource.Version
	err = json.Unmarshal(content, &currentVersion)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal version info file: %w", err)
	}

	if currentVersion.Type != config.Source.TypeFilter || currentVersion.Version == config.Version.Version {
		return "[]", nil
	}

	response := CheckResponse{currentVersion}

	output, err := json.Marshal(response)

	return string(output), err
}
