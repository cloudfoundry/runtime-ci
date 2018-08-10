package common
import (
	"path/filepath"
	"fmt"
	"io/ioutil"
	"strings"
)

type MarshalFunc func(interface{}) ([]byte, error)
type UnmarshalFunc func([]byte, interface{}) error

type Release struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Version string `yaml:"version"`
	SHA1    string `yaml:"sha1"`
}

func GetReleaseFromFile(buildDir, releaseName string) (Release, error) {
	newRelease := Release{
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
			return Release{}, shaErr
		}
		return Release{}, urlErr
	}

	if verErr != nil {
		return Release{}, verErr
	}

	newRelease.URL = strings.TrimSpace(string(url))
	newRelease.SHA1 = strings.TrimSpace(string(sha1))
	newRelease.Version = strings.TrimSpace(string(version))

	return newRelease, nil
}