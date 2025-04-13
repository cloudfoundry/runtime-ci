package opsfile

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
)

type Op struct {
	TypeField string      `yaml:"type"`
	Path      string      `yaml:"path"`
	Value     interface{} `yaml:"value,omitempty"`
}

const BadReleaseOpsFormatErrorMessage = "cannot update ops file: make sure all release information is updated with one operation in the ops file"

func UpdateReleases(releaseNames []string, buildDir string, opsFile []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	if len(releaseNames) == 0 {
		err := errors.New("releaseNames provided to UpdateReleases must contain at least one release name")
		return nil, common.NoOpsFileChangesCommitMessage, err
	}

	var deserializedOpsFile []Op
	if err := unmarshalFunc(opsFile, &deserializedOpsFile); err != nil {
		return nil, common.NoOpsFileChangesCommitMessage, err
	}

	var changes []string
	var releaseFound bool

	for _, op := range deserializedOpsFile {
		if op.TypeField == "replace" && strings.HasPrefix(op.Path, "/releases/") {
			valueMap, ok := op.Value.(map[interface{}]interface{})
			if !ok {
				return nil, common.NoOpsFileChangesCommitMessage, errors.New(BadReleaseOpsFormatErrorMessage)
			}

			for _, releaseName := range releaseNames {
				if valueMap["name"] == releaseName {
					releaseFound = true
					oldRelease := common.Release{
						Name:    strings.TrimSpace(valueMap["name"].(string)),
						Version: strings.TrimSpace(valueMap["version"].(string)),
					}

					newRelease, err := common.GetReleaseFromFile(buildDir, releaseName)
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

	if !releaseFound {
		err := fmt.Errorf("opsfile does not contain release named %s", releaseNames[0])
		return nil, common.NoOpsFileChangesCommitMessage, err
	}

	updatedOpsFile, err := marshalFunc(&deserializedOpsFile)
	if err != nil {
		return nil, common.NoOpsFileChangesCommitMessage, err
	}

	changeMessage := "No opsfile release updates"
	if len(changes) > 0 {
		changeMessage = fmt.Sprintf("Updated ops file(s) with %s", strings.Join(changes, ", "))
	}

	return updatedOpsFile, changeMessage, nil
}
