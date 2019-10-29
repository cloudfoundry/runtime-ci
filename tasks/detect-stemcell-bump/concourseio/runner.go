package concourseio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
)

type Runner struct {
	stemcell         bosh.Stemcell
	manifestStemcell bosh.Stemcell
	bumpType string

	In  Inputs
	Out Outputs
}

type Inputs struct {
	cfDeploymentDir string
	stemcellDir     string
}

type Outputs struct {
	bumpTypeDir string
}

func NewRunner(buildDir string) (Runner, error) {
	inputs, err := setupInputs(buildDir)
	if err != nil {
		return Runner{}, err
	}

	outputs, err := setupOutputs(buildDir)
	if err != nil {
		return Runner{}, err
	}

	return Runner{In: inputs, Out: outputs}, nil
}

func (r *Runner) ReadStemcell() error {
	var err error

	r.stemcell.Version, err = readFile(filepath.Join(r.In.stemcellDir, "version"))
	if err != nil {
		return err
	}

	url, err := readFile(filepath.Join(r.In.stemcellDir, "url"))
	if err != nil {
		return err
	}

	r.stemcell.OS, err = parseOSfromURL(url)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) ReadCFDeploymentStemcell() error {
	contents, err := readFile(filepath.Join(r.In.cfDeploymentDir, "cf-deployment.yml"))
	if err != nil {
		return err
	}

	manifest, err := bosh.NewManifestFromFile([]byte(contents))
	if err != nil {
		return err
	}

	r.manifestStemcell = manifest.Stemcells[0]

	return nil
}

func (r *Runner) DetectStemcellBump() error {
	bumpType, err := r.stemcell.DetectBumpTypeFrom(r.manifestStemcell)
	if err != nil {
		return err
	}

	r.bumpType = bumpType
	return nil
}

func (r *Runner) WriteStemcellBumpTypeToFile() error {
	err := ioutil.WriteFile(filepath.Join(r.Out.bumpTypeDir, "result"), []byte(r.bumpType), 0644)
	if err != nil {
		return err
	}
	return nil
}

func readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		pathDir := filepath.Base(filepath.Dir(path))
		return "", fmt.Errorf("missing files: '%s'", filepath.Join(pathDir, filepath.Base(path)))
	}

	return string(content), err
}

func setupInputs(buildDir string) (Inputs, error) {
	cfDeploymentDir, err := buildSubDir(buildDir, "cf-deployment")
	if err != nil {
		return Inputs{}, err
	}

	stemcellDir, err := buildSubDir(buildDir, "stemcell")
	if err != nil {
		return Inputs{}, err
	}

	return Inputs{cfDeploymentDir, stemcellDir}, nil
}

func setupOutputs(buildDir string) (Outputs, error) {
	stemcellBumpTypeDir, err := buildSubDir(buildDir, "stemcell-bump-type")
	if err != nil {
		return Outputs{}, err
	}

	return Outputs{stemcellBumpTypeDir}, nil
}

func buildSubDir(buildDir, subDir string) (string, error) {
	dir := filepath.Join(buildDir, subDir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("missing sub directory '%s' in build directory '%s'", subDir, buildDir)
	}

	return dir, nil
}

func parseOSfromURL(url string) (string, error) {
	versionRegex := regexp.MustCompile(`(ubuntu-.*)-go_agent.tgz`)

	allMatches := versionRegex.FindAllStringSubmatch(url, 1)

	if len(allMatches) != 1 {
		return "", fmt.Errorf("stemcell URL does not contain a supported os (i.e. ubuntu): %s", strings.Trim(url, "\n"))
	}

	osMatch := allMatches[0][1]
	return osMatch, nil
}
