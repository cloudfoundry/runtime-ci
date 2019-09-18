package compiledrelease

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/concourseio"
	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/manifest"
	"gopkg.in/yaml.v3"
)

type OpsfileUpdater struct {
	ops []Op

	compiledReleasesDir string
	opsFileOutPath      string

	releases []manifest.Release
}

type Op struct {
	Type  string
	Path  string
	Value manifest.Release
}

type NoReleasesErr struct{}

var _ error = new(NoReleasesErr)

func (*NoReleasesErr) Error() string {
	return "no releases found"
}

func NewOpsfileUpdater(compiledReleasesInDir string, opsFileOutPath string) *OpsfileUpdater {
	return &OpsfileUpdater{compiledReleasesDir: compiledReleasesInDir, opsFileOutPath: opsFileOutPath}
}

func (o *OpsfileUpdater) Load() error {
	err := filepath.Walk(o.compiledReleasesDir, o.extractReleases())
	if err != nil {
		return err
	}

	if len(o.releases) == 0 {
		return new(NoReleasesErr)
	}

	return nil
}

func (o *OpsfileUpdater) extractReleases() filepath.WalkFunc {
	const compiledReleaseGCSPrefix = "https://storage.googleapis.com/cf-deployment-compiled-releases"

	versionRegexString := `(.*)-([\d.]+)-(.*)-([\d.]+)-\d+-\d+-\d+.tgz`
	versionRegex := regexp.MustCompile(versionRegexString)

	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		tarballName := info.Name()

		allMatches := versionRegex.FindAllStringSubmatch(tarballName, 1)

		if len(allMatches) != 1 {
			return fmt.Errorf("invalid tarball name syntax: %s", tarballName)
		}

		sha1, err := computeSHA1Sum(path)
		if err != nil {
			return err
		}

		release := manifest.Release{
			Name: allMatches[0][1],
			SHA1: sha1,
			Stemcell: manifest.Stemcell{
				OS:      allMatches[0][3],
				Version: allMatches[0][4],
			},
			Version: allMatches[0][2],
			URL:     fmt.Sprintf("%s/%s", compiledReleaseGCSPrefix, tarballName),
		}

		o.releases = append(o.releases, release)

		return nil
	}
}

func (o *OpsfileUpdater) Update(stemcell manifest.Stemcell) error {
	if len(o.releases) == 0 {
		return new(NoReleasesErr)
	}

	for _, release := range o.releases {
		if release.Stemcell != stemcell {
			return errors.New("stemcell mismatch")
		}

		op := Op{
			Type:  "replace",
			Path:  fmt.Sprintf("/releases/name=%s", release.Name),
			Value: release,
		}

		o.ops = append(o.ops, op)
	}

	return nil
}

func (o OpsfileUpdater) Write() error {
	if len(o.ops) == 0 {
		return new(NoReleasesErr)
	}

	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "## GENERATED FILE. DO NOT EDIT")
	fmt.Fprintln(buf, "---")

	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(o.ops); err != nil {
		return err
	}

	return ioutil.WriteFile(o.opsFileOutPath, buf.Bytes(), 0755)
}

func computeSHA1Sum(filepath string) (string, error) {
	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(fileContents)), nil
}

var _ concourseio.StemcellUpdater = new(OpsfileUpdater)
