package opsfile

import (
	"errors"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"fmt"
	"strings"
	"path/filepath"
)

func UpdateCompiledReleases(releaseNames []string, buildDir string, opsFile []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	if releaseNames == nil || len(releaseNames) == 0 {
		err := errors.New("releaseNames provided to UpdateReleases must contain at least one release name")
		return nil, "", err
	}

	var deserializedOpsFile []Op
	if err := unmarshalFunc(opsFile, &deserializedOpsFile); err != nil {
		return nil, "", err
	}

	for i, op := range deserializedOpsFile {
		for _, releaseName := range releaseNames {
			newRelease, err := getCompiledReleaseForBuild(buildDir, releaseName)
			if err != nil {
				return nil, "", err
			}

			matchingURLPath := fmt.Sprintf("/path:/releases/name=%s/url", releaseName)
			matchingVersionPath := fmt.Sprintf("/path:/releases/name=%s/version", releaseName)
			matchingSha1Path := fmt.Sprintf("/path:/releases/name=%s/sha1", releaseName)
			matchingStemcellPath := fmt.Sprintf("/path:/releases/name=%s/stemcell?", releaseName)

			if strings.Contains(op.Path, matchingURLPath) {
				deserializedOpsFile[i].Value = newRelease.URL
			} else if strings.Contains(op.Path, matchingVersionPath) {

			} else if strings.Contains(op.Path, matchingSha1Path) {

			} else if strings.Contains(op.Path, matchingStemcellPath) {

			} else {

			}

		}
	}

	updatedOpsFile, err := marshalFunc(&deserializedOpsFile)
	if err != nil {
		return nil, "", err
	}

	return updatedOpsFile, "", nil
}

func getCompiledReleaseForBuild(buildDir, releaseName string)  (common.Release, error) {
	newRelease := common.Release{
		Name: releaseName,
	}

	releaseTarballGlob := filepath.Join(buildDir, fmt.Sprintf("%s-compiled-release-tarball", releaseName), "*.tgz")

	matches, err := filepath.Glob(releaseTarballGlob)
	if err != nil {
		return common.Release{}, err
	}

	if len(matches) != 1 {
		return common.Release{}, errors.New("expected to find exactly 1 compiled release tarball")
	}

	matchedFile := matches[0]

	newRelease.URL = matchedFile
	fmt.Println(matchedFile)

	return newRelease, nil

}