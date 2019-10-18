package runner

import (
	"encoding/json"
	"fmt"
	"stemcell-version-bump/resource"
)

//go:generate counterfeiter . Getter
type Getter interface {
	Get(bucketName string, objectPath string) ([]byte, error)
}

func Check(config resource.Config, getter Getter) (string, error) {
	content, err := getter.Get(config.Source.BucketName, config.Source.FileName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version info from bucket/file (%s, %s): %w", config.Source.BucketName, config.Source.FileName, err)
	}

	var versionInfo resource.Version
	err = json.Unmarshal(content, &versionInfo)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal version info file: %w", err)
	}

	if versionInfo.Type != config.Source.TypeFilter || versionInfo.Version == config.Version.Version {
		return "[]", nil
	}

	return fmt.Sprintf("[%s]", string(content)), nil
}
