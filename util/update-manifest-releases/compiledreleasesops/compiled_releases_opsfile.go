package compiledreleasesops

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
)

const compiledReleasesURLPrefix = "https://storage.googleapis.com/cf-deployment-compiled-releases"

type Release struct {
	Name     string                    `yaml:"name"`
	SHA1     string                    `yaml:"sha1"`
	Stemcell common.StemcellForRelease `yaml:"stemcell,omitempty"`
	URL      string                    `yaml:"url"`
	Version  string                    `yaml:"version"`
}

func UpdateCompiledReleases(releaseNames []string, buildDir string, opsFile []byte, marshalFunc common.MarshalFunc, unmarshalFunc common.UnmarshalFunc) ([]byte, string, error) {
	if len(releaseNames) == 0 {
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
		var newRelease Release
		var err error

		foundRelease := false

		matchingReleasePath := fmt.Sprintf("/releases/name=%s", releaseName)

		for i, op := range deserializedOpsFile {
			if op.Path == matchingReleasePath {
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

func appendNewRelease(newRelease Release, opsFile []opsfile.Op) []opsfile.Op {
	newReleaseOps := opsfile.Op{
		TypeField: "replace",
		Path:      fmt.Sprintf("/releases/name=%s", newRelease.Name),
		Value:     newRelease,
	}

	return append(opsFile, newReleaseOps)
}

func getCompiledReleaseForBuild(buildDir, releaseName string) (Release, error) {
	releaseTarballGlob := filepath.Join(buildDir, fmt.Sprintf("%s-compiled-release-tarball", releaseName), "*.tgz")

	matches, err := filepath.Glob(releaseTarballGlob)
	if err != nil {
		return Release{}, err
	}
	if len(matches) != 1 {
		return Release{}, errors.New("expected to find exactly 1 compiled release tarball")
	}

	releaseTarballPath := matches[0]
	releaseTarballName := filepath.Base(releaseTarballPath)

	release := Release{Name: releaseName}
	release.Version, release.Stemcell.Version, release.Stemcell.OS, err = common.InfoFromTarballName(releaseTarballName, releaseName)
	if err != nil {
		return Release{}, err
	}

	release.SHA1, err = computeSha256Sum(releaseTarballPath)
	if err != nil {
		return Release{}, err
	}

	release.URL = fmt.Sprintf("%s/%s", compiledReleasesURLPrefix, releaseTarballName)

	return release, nil
}

func computeSha256Sum(filepath string) (string, error) {
	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return "sha256:" + fmt.Sprintf("%x", sha256.Sum256(fileContents)), nil
}
