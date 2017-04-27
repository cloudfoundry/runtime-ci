package opsfile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
)

type Op struct {
	TypeField string      `yaml:"type"`
	Path      string      `yaml:"path"`
	Value     interface{} `yaml:"value"`
}

func UpdateReleases(releaseNames []string, buildDir string, opsFile []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	var changes []string

	var deserializedOpsFile []Op
	if err := unmarshalFunc(opsFile, &deserializedOpsFile); err != nil {
		return nil, "", err
	}

	for _, op := range deserializedOpsFile {
		if op.Path == "/releases/-" {
			valueMap := op.Value.(map[interface{}]interface{})
			for _, releaseName := range releaseNames {
				if valueMap["name"] == releaseName {
					newRelease, err := getReleaseFromFile(buildDir, releaseName)
					if err != nil {
						return nil, "", err
					}

					if newRelease.SHA1 != strings.TrimSpace(valueMap["sha1"].(string)) {
						changes = append(changes, fmt.Sprintf("%s-release", newRelease.Name))
					}

					valueMap["sha1"] = newRelease.SHA1
					valueMap["url"] = newRelease.URL
					valueMap["version"] = newRelease.Version
				}
			}
		}
	}

	updatedOpsFile, err := marshalFunc(&deserializedOpsFile)
	if err != nil {
		return nil, "", err
	}

	changeMessage := "No opsfile release updates"
	if len(changes) > 0 {
		changeMessage = fmt.Sprintf("Updated opsfile with %s", strings.Join(changes, ", "))
	}

	return updatedOpsFile, changeMessage, nil
}

func getReleaseFromFile(buildDir, releaseName string) (*common.Release, error) {
	newRelease := &common.Release{
		Name: releaseName,
	}
	releasePath := filepath.Join(buildDir, fmt.Sprintf("%s-release", releaseName))

	sha1, err := ioutil.ReadFile(filepath.Join(releasePath, "sha1"))
	newRelease.SHA1 = strings.TrimSpace(string(sha1))
	if err != nil {
		return nil, err
	}

	url, err := ioutil.ReadFile(filepath.Join(releasePath, "url"))
	if err != nil {
		return nil, err
	}
	newRelease.URL = strings.TrimSpace(string(url))

	version, err := ioutil.ReadFile(filepath.Join(releasePath, "version"))
	if err != nil {
		return nil, err
	}
	newRelease.Version = strings.TrimSpace(string(version))

	return newRelease, nil
}
