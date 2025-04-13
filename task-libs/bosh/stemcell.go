package bosh

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

type Stemcell struct {
	Alias   string `yaml:",omitempty"`
	OS      string
	Version string
}

// NewStemcellFromInput creates a Stemcell from a stemcell concourse resource
func NewStemcellFromInput(stemcellDir string) (Stemcell, error) {
	var stemcell Stemcell

	version, err := readFile(filepath.Join(stemcellDir, "version"))
	if err != nil {
		return stemcell, err
	}
	stemcell.Version = strings.Trim(version, "\n")

	url, err := readFile(filepath.Join(stemcellDir, "url"))
	if err != nil {
		return stemcell, err
	}

	stemcell.OS, err = parseOSfromURL(url)
	if err != nil {
		return stemcell, err
	}

	return stemcell, nil
}

func parseOSfromURL(url string) (string, error) {
	versionRegex := regexp.MustCompile(`(ubuntu-.*)-go_agent.tgz`)

	allMatches := versionRegex.FindAllStringSubmatch(url, 1)

	if len(allMatches) != 1 {
		return "", fmt.Errorf("stemcell URL does not contain an ubuntu stemcell: %s", strings.Trim(url, "\n"))
	}

	osMatch := allMatches[0][1]
	return osMatch, nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("missing files: %T", err)
	}

	return string(content), err
}

func (s Stemcell) CompareVersion(base Stemcell) (int, error) {
	if s.OS != base.OS {
		return 0, fmt.Errorf("stemcell OS mismatch: %q vs %q", s.OS, base.OS)
	}

	version, err := semver.Parse(fmt.Sprintf("%s.0", s.Version))
	if err != nil {
		return 0, fmt.Errorf("failed to parse stemcell version %q: %w", s.Version, err)
	}

	baseVersion, err := semver.Parse(fmt.Sprintf("%s.0", base.Version))
	if err != nil {
		return 0, fmt.Errorf("failed to parse stemcell version %q: %w", base.Version, err)
	}

	return version.Compare(baseVersion), nil
}

func (s Stemcell) DetectBumpTypeFrom(base Stemcell) (string, error) {
	if s.OS != base.OS {
		return "major", nil
	}

	version, err := semver.Parse(fmt.Sprintf("%s.0", s.Version))
	if err != nil {
		return "", fmt.Errorf("failed to parse stemcell version %q: %w", s.Version, err)
	}

	baseVersion, err := semver.Parse(fmt.Sprintf("%s.0", base.Version))
	if err != nil {
		return "", fmt.Errorf("failed to parse stemcell version %q: %w", base.Version, err)
	}

	if version.Major > baseVersion.Major {
		return "major", nil
	}

	if version.Major == baseVersion.Major && version.Minor > baseVersion.Minor {
		return "minor", nil
	}

	return "", fmt.Errorf("change from %s to %s is not a forward bump", base.Version, s.Version)
}
