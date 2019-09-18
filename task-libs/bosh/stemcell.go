package bosh

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Stemcell struct {
	Alias   string
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
	content, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("missing files: %T", err)
	}

	return string(content), err
}
