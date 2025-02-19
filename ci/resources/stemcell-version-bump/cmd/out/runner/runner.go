package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"stemcell-version-bump/resource"
)

type OutResponse struct {
	Version resource.Version `json:"version"`
}

func NewVersion(request resource.OutRequest) (resource.Version, error) {
	versionContent, err := os.ReadFile(request.Params.VersionFile)
	if err != nil {
		dir, _ := os.Getwd()
		log.Printf("Current working directory: %s\n", dir)
		return resource.Version{}, fmt.Errorf("reading version file: %w", err)
	}

	bumpTypeContent, err := os.ReadFile(request.Params.TypeFile)
	if err != nil {
		return resource.Version{}, fmt.Errorf("reading bump type file: %w", err)
	}

	stemcellVersionBump := resource.Version{
		Type:    string(bumpTypeContent),
		Version: string(versionContent),
	}

	err = validateBumpType(stemcellVersionBump.Type)
	if err != nil {
		return resource.Version{}, err
	}

	return stemcellVersionBump, nil
}

//go:generate counterfeiter . Putter
type Putter interface {
	Put(bucketName string, objectPath string, contents []byte) error
}

func UploadVersion(request resource.OutRequest, putter Putter, version resource.Version) error {
	contents, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("failed to marshal version bump info for put: %w", err)
	}

	err = putter.Put(request.Source.BucketName, request.Source.FileName, contents)
	if err != nil {
		return fmt.Errorf("updating version info in bucket/file (%s, %s): %w", request.Source.BucketName, request.Source.FileName, err)
	}

	return nil
}

func GenerateResourceOutput(version resource.Version) (string, error) {
	output, err := json.Marshal(OutResponse{Version: version})
	if err != nil {
		return "", fmt.Errorf("failed to marshal concourse output version info: %w", err)
	}

	return string(output), nil
}

func validateBumpType(bumpType string) error {
	if bumpType != "minor" && bumpType != "major" {
		return fmt.Errorf("invalid bump type: %q", bumpType)
	}

	return nil
}
