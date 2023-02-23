package concourseio

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
)

type Runner struct {
	stemcell bosh.Stemcell

	In  Inputs
	Out Outputs
}

type Inputs struct {
	cfDeploymentDir string
	stemcellDir     string
}

type Outputs struct {
	UpdatedCFDeploymentDir string
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
	version, err := readFile(filepath.Join(r.In.stemcellDir, "version"))
	if err != nil {
		return err
	}
	r.stemcell.Version = version

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

type UpdateFunc func([]byte, bosh.Stemcell) ([]byte, error)

func (r *Runner) UpdateManifest(updateFunction UpdateFunc) error {
	originalFile := manifestPath(r.In.cfDeploymentDir)
	manifest, err := os.ReadFile(originalFile)
	if err != nil {
		return err
	}

	updatedContent, err := updateFunction(manifest, r.stemcell)
	if err != nil {
		return err
	}

	updatedFile := manifestPath(r.Out.UpdatedCFDeploymentDir)
	if err := os.WriteFile(updatedFile, updatedContent, 0755); err != nil {
		return err
	}

	return nil
}

func (r *Runner) WriteCommitMessage(commitMessagePath string) error {
	commitMessage := fmt.Sprintf("Update stemcell to %s \"%s\"", r.stemcell.OS, r.stemcell.Version)

	err := os.WriteFile(commitMessagePath, []byte(commitMessage), 0644)
	if err != nil {
		return err
	}

	return nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		pathDir := filepath.Base(filepath.Dir(path))
		return "", fmt.Errorf("missing files: '%s'", filepath.Join(pathDir, filepath.Base(path)))
	}

	return string(content), err
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
	updatedCFDeploymentDir, err := buildSubDir(buildDir, "updated-cf-deployment")
	if err != nil {
		return Outputs{}, err
	}

	return Outputs{updatedCFDeploymentDir}, nil
}

func buildSubDir(buildDir, subDir string) (string, error) {
	dir := filepath.Join(buildDir, subDir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("missing sub directory '%s' in build directory '%s'", subDir, buildDir)
	}

	return dir, nil
}

func manifestPath(cfDir string) string {
	return filepath.Join(cfDir, "cf-deployment.yml")
}
