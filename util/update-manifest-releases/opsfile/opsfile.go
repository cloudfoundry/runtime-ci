package opsfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
)

type Op struct {
	TypeField string      `yaml:"type"`
	Path      string      `yaml:"path"`
	Value     interface{} `yaml:"value,omitempty"`
}

func UpdateReleases(releaseNames []string, buildDir string, opsFile []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	if len(releaseNames) > 1 {
		err := errors.New("Releases array should have only one release")
		return nil, "", err
	}

	if releaseNames == nil || len(releaseNames) == 0 {
		err := errors.New("releaseNames provided to UpdateReleases must contain at least one release name")
		return nil, "", err
	}

	var deserializedOpsFile []Op
	if err := unmarshalFunc(opsFile, &deserializedOpsFile); err != nil {
		return nil, "", err
	}

	var changes []string
	var releaseFound bool

	for _, op := range deserializedOpsFile {
		if op.TypeField == "replace" && strings.HasPrefix(op.Path, "/releases/") {
			valueMap := op.Value.(map[interface{}]interface{})
			for _, releaseName := range releaseNames {
				if valueMap["version"] == nil {
					err := errors.New(fmt.Sprintf("No version for %s in ops-file", releaseName))
					return nil, "", err
				}
				if valueMap["name"] == releaseName {
					releaseFound = true
					oldRelease := common.Release{
						Name:    strings.TrimSpace(valueMap["name"].(string)),
						Version: strings.TrimSpace(valueMap["version"].(string)),
					}

					newRelease, err := getReleaseFromFile(buildDir, releaseName)
					if err != nil {
						return nil, "", err
					}

					if sha, ok := valueMap["sha1"]; ok {
						oldRelease.SHA1 = strings.TrimSpace(sha.(string))
						valueMap["sha1"] = newRelease.SHA1
					}

					if url, ok := valueMap["url"]; ok {
						oldRelease.URL = strings.TrimSpace(url.(string))
						valueMap["url"] = newRelease.URL
					}

					valueMap["version"] = newRelease.Version

					if newRelease != oldRelease {
						changes = append(changes, fmt.Sprintf("%s-release %s", newRelease.Name, newRelease.Version))
					}
				}
			}
		}
	}

	if releaseFound == false {
		err := errors.New(fmt.Sprintf("Opsfile does not contain release named %s", releaseNames[0]))
		return nil, "", err
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

func getReleaseFromFile(buildDir, releaseName string) (common.Release, error) {
	newRelease := common.Release{
		Name: releaseName,
	}
	releasePath := filepath.Join(buildDir, fmt.Sprintf("%s-release", releaseName))

	sha1, shaErr := ioutil.ReadFile(filepath.Join(releasePath, "sha1"))
	url, urlErr := ioutil.ReadFile(filepath.Join(releasePath, "url"))
	version, verErr := ioutil.ReadFile(filepath.Join(releasePath, "version"))

	isShaErr := shaErr != nil
	isUrlErr := urlErr != nil

	// We accept neither or both of "sha1" and "url".  If we error out on only one or the other, something is wrong.
	if isShaErr != isUrlErr {
		if isShaErr {
			return common.Release{}, shaErr
		}
		return common.Release{}, urlErr
	}

	if verErr != nil {
		return common.Release{}, verErr
	}

	newRelease.URL = strings.TrimSpace(string(url))
	newRelease.SHA1 = strings.TrimSpace(string(sha1))
	newRelease.Version = strings.TrimSpace(string(version))

	return newRelease, nil
}
