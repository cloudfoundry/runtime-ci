package runner

import (
	"encoding/json"
	"fmt"

	"stemcell-version-bump/resource"
)

type InResponse struct {
	Version resource.Version `json:"version"`
}

//go:generate counterfeiter . Getter
type Getter interface {
	Get(bucketName string, objectPath string) ([]byte, error)
}

func In(request resource.CheckInRequest, getter Getter) (string, error) {
	content, err := getter.Get(request.Source.BucketName, request.Source.FileName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version info from bucket/file (%s, %s): %w", request.Source.BucketName, request.Source.FileName, err)
	}

	var currentVersion resource.Version

	err = json.Unmarshal(content, &currentVersion)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal version info file: %w", err)
	}

	if currentVersion != request.Version {
		return "", fmt.Errorf("failed to retrieve specified version: requested %s, found %s", request.Version, currentVersion)
	}

	response := InResponse{
		Version: currentVersion,
	}

	output, err := json.Marshal(response)

	return string(output), err
}
