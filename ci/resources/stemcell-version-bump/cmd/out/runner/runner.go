package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"stemcell-version-bump/resource"
)

func ReadVersionBump(outRequest resource.OutRequest) ([]byte, error) {
	versionContent, err := ioutil.ReadFile(outRequest.Params.VersionFile)
	if err != nil {
		dir, _ := os.Getwd()
		log.Printf("Current working directory: %s\n", dir)
		return nil, fmt.Errorf("reading version file: %w", err)
	}

	bumpTypeContent, err := ioutil.ReadFile(outRequest.Params.TypeFile)
	if err != nil {
		return nil, fmt.Errorf("reading bump type file: %w", err)
	}

	stemcellVersionBump := resource.Version{
		Type:    string(bumpTypeContent),
		Version: string(versionContent),
	}

	err = validateBumpType(stemcellVersionBump.Type)
	if err != nil {
		return nil, err
	}

	return json.Marshal(stemcellVersionBump)
}

//go:generate counterfeiter . Putter
type Putter interface {
	Put(bucketName string, objectPath string, contents []byte) error
}

func Out(config resource.OutRequest, putter Putter, contents []byte) error {
	err := putter.Put(config.Source.BucketName, config.Source.FileName, contents)
	if err != nil {
		return fmt.Errorf("updating version info in bucket/file (%s, %s): %w", config.Source.BucketName, config.Source.FileName, err)
	}
	return nil
}

func validateBumpType(bumpType string) error {
	if bumpType != "minor" && bumpType != "major" {
		return fmt.Errorf("invalid bump type: %q", bumpType)
	}

	return nil
}
