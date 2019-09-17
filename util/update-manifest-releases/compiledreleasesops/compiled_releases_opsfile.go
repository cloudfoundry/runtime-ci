package compiledreleasesops

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
)

const compiledReleasesURLPrefix = "https://storage.googleapis.com/cf-deployment-compiled-releases"

func UpdateCompiledReleases(releaseNames []string, buildDir string, opsFile []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	if releaseNames == nil || len(releaseNames) == 0 {
		err := errors.New("releaseNames provided to UpdateReleases must contain at least one release name")
		return nil, "", err
	}

	var deserializedOpsFile []opsfile.Op
	if err := unmarshalFunc(opsFile, &deserializedOpsFile); err != nil {
		return nil, "", err
	}

	var commitMessage string

	for _, releaseName := range releaseNames {
		fmt.Printf("Updating release %s...\n", releaseName)
		var newRelease common.Release
		var err error

		foundRelease := false

		matchingReleasePath := fmt.Sprintf("/releases/name=%s", releaseName)

		for i, op := range deserializedOpsFile {

			if strings.Contains(op.Path, matchingReleasePath) {
				newRelease, err = getCompiledReleaseForBuild(buildDir, releaseName)
				if err != nil {
					return nil, "", err
				}
				foundRelease = true
				deserializedOpsFile[i].Value = newRelease
				commitMessage = fmt.Sprintf("Updated compiled releases with %s %s", newRelease.Name, newRelease.Version)
			}
		}

		if !foundRelease {
			newRelease, err = getCompiledReleaseForBuild(buildDir, releaseName)
			if err != nil {
				return nil, "", err
			}
			deserializedOpsFile = appendNewRelease(newRelease, deserializedOpsFile)
			commitMessage = fmt.Sprintf("Updated compiled releases with %s %s", newRelease.Name, newRelease.Version)
		}
	}

	updatedOpsFile, err := marshalFunc(&deserializedOpsFile)
	if err != nil {
		return nil, "", err
	}

	return updatedOpsFile, commitMessage, nil
}

func newStemcellOp(newRelease common.Release) opsfile.Op {
	return opsfile.Op{
		TypeField: "replace",
		Path:      fmt.Sprintf("/releases/name=%s/stemcell?", newRelease.Name),
		Value:     newRelease.Stemcell,
	}
}

func appendNewRelease(newRelease common.Release, opsFile []opsfile.Op) []opsfile.Op {
	newReleaseOps := opsfile.Op{
		TypeField: "replace",
		Path:      fmt.Sprintf("/releases/name=%s", newRelease.Name),
		Value:     newRelease,
	}

	return append(opsFile, newReleaseOps)
}

func getCompiledReleaseForBuild(buildDir, releaseName string) (common.Release, error) {
	release, err := common.GetReleaseFromFile(buildDir, releaseName)
	if err != nil {
		return common.Release{}, fmt.Errorf("could not find necessary release info: %s", err.Error())
	}

	releaseTarballGlob := filepath.Join(buildDir, fmt.Sprintf("%s-compiled-release-tarball", releaseName), "*.tgz")

	matches, err := filepath.Glob(releaseTarballGlob)
	if err != nil {
		return common.Release{}, err
	}
	if len(matches) != 1 {
		return common.Release{}, errors.New("expected to find exactly 1 compiled release tarball")
	}

	releaseTarballPath := matches[0]
	releaseTarballName := filepath.Base(releaseTarballPath)

	release.Stemcell.Version, release.Stemcell.OS, err = common.StemcellInfoFromTarballName(releaseTarballName, release.Name, release.Version)
	if err != nil {
		return common.Release{}, err
	}

	release.SHA1, err = computeSha1Sum(releaseTarballPath)
	if err != nil {
		return common.Release{}, err
	}

	release.URL = fmt.Sprintf("%s/%s", compiledReleasesURLPrefix, releaseTarballName)

	return release, nil
}

func computeSha1Sum(filepath string) (string, error) {
	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(fileContents)), nil
}
