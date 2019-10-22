package runner

import (
	"fmt"
	"io/ioutil"
	"stemcell-version-bump/resource"
)

//go:generate counterfeiter . Putter
type Putter interface {
	Put(bucketName string, objectPath string, contents []byte) error
}

func ReadVersionFile(outRequest resource.OutRequest) ([]byte, error) {
	return ioutil.ReadFile(outRequest.Params.VersionFile)
}

func Out(config resource.OutRequest, putter Putter, contents []byte) error {
	err := putter.Put(config.Source.BucketName, config.Source.FileName, contents)
	if err != nil {
		return fmt.Errorf("updating version info in bucket/file (%s, %s): %w", config.Source.BucketName, config.Source.FileName, err)
	}
	return nil
}
